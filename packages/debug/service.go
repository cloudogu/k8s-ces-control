package debug

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	"github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/ces-control-api/generated/types"
)

const (
	maintenanceTitle          = "Service unavailable"
	activateMaintenanceText   = "Activating debug mode"
	deactivateMaintenanceText = "Deactivating debug mode"
	logLevelDebug             = "DEBUG"
	interErrMsg               = "internal error"
)

type defaultDebugModeService struct {
	pbMaintenance.UnimplementedDebugModeServer
	debugModeClient       debugModeInterface
	debugModeRegistry     debugModeRegistry
	maintenanceModeSwitch maintenanceModeSwitch
	doguInterActor        doguInterActor
}

// NewDebugModeService returns an instance of debugModeService.
func NewDebugModeService(debugMode debugModeInterface, doguInterActor doguInterActor, doguConfigRepository doguConfigRepository, globalConfigRepository globalConfigRepository, doguDescriptorGetter doguDescriptorGetter, clusterClient clusterClientSet, namespace string) *defaultDebugModeService {
	cmDebugModeRegistry := NewConfigMapDebugModeRegistry(doguConfigRepository, doguDescriptorGetter, clusterClient, namespace)
	return &defaultDebugModeService{
		debugModeClient:       debugMode,
		debugModeRegistry:     cmDebugModeRegistry,
		maintenanceModeSwitch: NewDefaultMaintenanceModeSwitch(globalConfigRepository),
		doguInterActor:        doguInterActor,
	}
}

// Enable enables the debug mode, sets dogu log level to debug and restarts all dogus.
func (s *defaultDebugModeService) Enable(ctx context.Context, req *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("Starting to enable debug-mode...")

	timestamp := time.Now().Add(time.Duration(req.Timer) * time.Minute)

	debugMode, err := s.debugModeClient.Get(ctx, "debug-mode", metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		debugMode = &v1.DebugMode{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: "debug-mode",
			},
			Spec: v1.DebugModeSpec{
				DeactivateTimestamp: metav1.NewTime(timestamp),
				TargetLogLevel:      "DEBUG",
			},
		}
		_, err = s.debugModeClient.Create(ctx, debugMode, metav1.CreateOptions{})
		if err != nil {
			logrus.Errorf("ERROR: failed to create debug-mode: %v", err)
			return nil, fmt.Errorf("ERROR: failed to create debug-mode: %q", err)
		}
	} else if err != nil {
		logrus.Errorf("ERROR: failed to get debug-mode: %v", err)
		return nil, fmt.Errorf("ERROR: failed to get debug-mode: %q", err)
	}

	debugMode.Spec.DeactivateTimestamp = metav1.NewTime(timestamp)

	_, err = s.debugModeClient.Update(ctx, debugMode, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("ERROR: failed to update debug-mode: %v", err)
		return nil, fmt.Errorf("ERROR: failed to update debug-mode: %q", err)
	}

	return &types.BasicResponse{}, nil
}

func (s *defaultDebugModeService) rollbackRestoreStartDisable(ctx context.Context, err error) error {
	startErr := s.doguInterActor.StartAllDogus(ctx)
	if startErr != nil {
		err = errors.Join(wrapRollBackErr(startErr), err)
	}

	return s.rollbackRestoreDisable(ctx, err)
}

func (s *defaultDebugModeService) rollbackRestoreDisable(ctx context.Context, err error) error {
	restoreErr := s.debugModeRegistry.RestoreDoguLogLevels(ctx)
	if restoreErr != nil {
		err = errors.Join(wrapRollBackErr(restoreErr), err)
	}

	return s.rollbackDisable(ctx, err)
}

func (s *defaultDebugModeService) rollbackDisable(ctx context.Context, err error) error {
	rollbackErr := s.debugModeRegistry.Disable(ctx)
	if rollbackErr != nil {
		err = errors.Join(wrapRollBackErr(rollbackErr), err)
	}
	return err
}

func wrapRollBackErr(err error) error {
	return fmt.Errorf("rollback error: %w", err)
}

// Disable returns an error because the method is unimplemented.
func (s *defaultDebugModeService) Disable(ctx context.Context, _ *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("Starting to disable debug-mode...")

	debugMode, err := s.debugModeClient.Get(ctx, "debug-mode", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("ERROR: failed to get debug-mode: %q", err)
	}

	debugMode.Spec.DeactivateTimestamp = metav1.NewTime(time.Now())

	_, err = s.debugModeClient.Update(ctx, debugMode, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("ERROR: failed to update debug-mode: %q", err)
	}

	return &types.BasicResponse{}, nil
}

// Status return an error because the method is unimplemented.
func (s *defaultDebugModeService) Status(ctx context.Context, _ *types.BasicRequest) (result *pbMaintenance.DebugModeStatusResponse, e error) {
	debugMode, err := s.debugModeClient.Get(ctx, "debug-mode", metav1.GetOptions{})

	if err != nil {
		return nil, fmt.Errorf("ERROR: failed to get debug-mode: %q", err)
	}

	return &pbMaintenance.DebugModeStatusResponse{IsEnabled: debugMode.Status.Phase != v1.DebugModeStatusCompleted, DisableAtTimestamp: debugMode.Spec.DeactivateTimestamp.Unix()}, nil
}

func noInheritCancel(_ context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
