package debug

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	debugClient "github.com/cloudogu/k8s-debug-mode-cr-lib/pkg/client/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
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
	debugModeRegistry     debugModeRegistry
	maintenanceModeSwitch maintenanceModeSwitch
	doguInterActor        doguInterActor
}

// NewDebugModeService returns an instance of debugModeService.
func NewDebugModeService(doguInterActor doguInterActor, doguConfigRepository doguConfigRepository, globalConfigRepository globalConfigRepository, doguDescriptorGetter doguDescriptorGetter, clusterClient clusterClientSet, namespace string) *defaultDebugModeService {
	cmDebugModeRegistry := NewConfigMapDebugModeRegistry(doguConfigRepository, doguDescriptorGetter, clusterClient, namespace)
	return &defaultDebugModeService{
		debugModeRegistry:     cmDebugModeRegistry,
		maintenanceModeSwitch: NewDefaultMaintenanceModeSwitch(globalConfigRepository),
		doguInterActor:        doguInterActor,
	}
}

func newDebugModeClient(restConfig rest.Config) (debugClient.DebugModeInterface, error) {
	config, err := debugClient.NewForConfig(&restConfig)
	if err != nil {
		return nil, err
	}

	return config.DebugMode("ecosystem"), nil
}

// Enable enables the debug mode, sets dogu log level to debug and restarts all dogus.
func (s *defaultDebugModeService) Enable(ctx context.Context, req *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("Starting to enable debug-mode...")

	client, err := newDebugModeClient(rest.Config{})
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Add(time.Duration(req.Timer) * time.Second)

	debugMode := v1.DebugMode{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1.DebugModeSpec{
			DeactivateTimestamp: metav1.NewTime(timestamp),
			TargetLogLevel:      "DEBUG",
		},
		Status: v1.DebugModeStatus{
			Phase:  v1.DebugModeStatusSet,
			Errors: "",
			Conditions: []metav1.Condition{{
				Type:    v1.ConditionLogLevelSet,
				Status:  metav1.ConditionFalse,
				Reason:  "DebugModeStatusSet",
				Message: "Set condition to false",
			},
			},
		},
	}

	_, err = client.Create(ctx, &debugMode, metav1.CreateOptions{})
	if err != nil {
		return nil, err
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

	client, err := newDebugModeClient(rest.Config{})
	if err != nil {
		return nil, err
	}

	debugMode, err := client.Get(ctx, "debugmode", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	debugMode.Spec.DeactivateTimestamp = metav1.NewTime(time.Now())

	_, err = client.Update(ctx, debugMode, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return &types.BasicResponse{}, nil
}

// Status return an error because the method is unimplemented.
func (s *defaultDebugModeService) Status(ctx context.Context, _ *types.BasicRequest) (*pbMaintenance.DebugModeStatusResponse, error) {
	client, err := newDebugModeClient(rest.Config{})
	if err != nil {
		return nil, err
	}

	debugMode, err := client.Get(ctx, "debugmode", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &pbMaintenance.DebugModeStatusResponse{IsEnabled: debugMode.Status.Phase == v1.DebugModeStatusCompleted, DisableAtTimestamp: debugMode.Spec.DeactivateTimestamp.Unix()}, nil
}

func createInternalError(err error) error {
	logrus.Error(err, interErrMsg)
	return status.Errorf(codes.Internal, "%v", err.Error())
}

func noInheritCancel(_ context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
