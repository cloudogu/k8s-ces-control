package debug

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-ces-control/packages/doguinteraction"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	globalConfig          configurationContext
	clientSet             clusterClientSet
	debugModeRegistry     debugModeRegistry
	maintenanceModeSwitch maintenanceModeSwitch
	namespace             string
	doguInterActor        doguInterActor
}

// NewDebugModeService returns an instance of debugModeService.
func NewDebugModeService(registry cesRegistry, doguReg doguRegistry, clusterClient clusterClientSet, namespace string) *defaultDebugModeService {
	cmDebugModeRegistry := NewConfigMapDebugModeRegistry(registry, doguReg, clusterClient, namespace)
	return &defaultDebugModeService{
		globalConfig:          registry.GlobalConfig(),
		clientSet:             clusterClient,
		debugModeRegistry:     cmDebugModeRegistry,
		maintenanceModeSwitch: NewDefaultMaintenanceModeSwitch(registry.GlobalConfig()),
		namespace:             namespace,
		doguInterActor:        doguinteraction.NewDefaultDoguInterActor(clusterClient, namespace, registry, doguReg),
	}
}

// Enable enables the debug mode, sets dogu log level to debug and restarts all dogus.
func (s *defaultDebugModeService) Enable(ctx context.Context, req *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("Starting to enable debug-mode...")

	err := s.maintenanceModeSwitch.ActivateMaintenanceMode(maintenanceTitle, activateMaintenanceText)
	if err != nil {
		return nil, createInternalError(fmt.Errorf("failed to activate maintenance mode: %w", err))
	}

	defer func() {
		err = s.maintenanceModeSwitch.DeactivateMaintenanceMode()
		if err != nil {
			logrus.Error(fmt.Errorf("failed to deactivate maintenance mode: %w", err), interErrMsg)
		}
		logrus.Info("...Finished enabling debug-mode.")
	}()

	err = s.debugModeRegistry.Enable(ctx, req.Timer)
	if err != nil {
		return nil, s.rollbackDisable(ctx, createInternalError(fmt.Errorf("failed to enable debug mode registry: %w", err)))
	}

	err = s.debugModeRegistry.BackupDoguLogLevels(ctx)
	if err != nil {
		return nil, s.rollbackDisable(ctx, createInternalError(fmt.Errorf("failed to backup dogu log levels: %w", err)))
	}

	err = s.doguInterActor.SetLogLevelInAllDogus(ctx, logLevelDebug)
	if err != nil {
		return nil, s.rollbackRestoreDisable(ctx, createInternalError(fmt.Errorf("failed to set dogu log levels to debug: %w", err)))
	}

	// Create new context because the admin dogu itself will be canceled
	noInheritedCtx, cancel := noInheritCancel(ctx)
	defer cancel()

	err = s.doguInterActor.StopAllDogus(noInheritedCtx)
	if err != nil {
		return nil, s.rollbackRestoreStartDisable(ctx, createInternalError(fmt.Errorf("failed to stop all dogus: %w", err)))
	}

	err = s.doguInterActor.StartAllDogus(noInheritedCtx)
	if err != nil {
		return nil, s.rollbackRestoreStartDisable(ctx, createInternalError(fmt.Errorf("failed to start all dogus %w", err)))
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
	err := s.maintenanceModeSwitch.ActivateMaintenanceMode(maintenanceTitle, deactivateMaintenanceText)
	if err != nil {
		return nil, createInternalError(fmt.Errorf("failed to activate maintenance mode: %w", err))
	}

	defer func() {
		err = s.maintenanceModeSwitch.DeactivateMaintenanceMode()
		if err != nil {
			logrus.Error(fmt.Errorf("failed to deactivate maintenance mode: %w", err), interErrMsg)
		}
		logrus.Info("...Finished disabling debug-mode.")
	}()

	err = s.debugModeRegistry.RestoreDoguLogLevels(ctx)
	if err != nil {
		return nil, createInternalError(fmt.Errorf("failed to restore log levels to ces registry: %w", err))
	}

	// Create new context because the admin dogu itself will be canceled
	noInheritedCtx, cancel := noInheritCancel(ctx)
	defer cancel()

	err = s.doguInterActor.StopAllDogus(noInheritedCtx)
	if err != nil {
		return nil, createInternalError(fmt.Errorf("failed to stop all dogus: %w", err))
	}

	err = s.doguInterActor.StartAllDogus(noInheritedCtx)
	if err != nil {
		return nil, createInternalError(fmt.Errorf("failed to start all dogus: %w", err))
	}

	err = s.debugModeRegistry.Disable(noInheritedCtx)
	if err != nil {
		return nil, createInternalError(fmt.Errorf("failed to disable the debug mode registry: %w", err))
	}

	return &types.BasicResponse{}, nil
}

// Status return an error because the method is unimplemented.
func (s *defaultDebugModeService) Status(ctx context.Context, _ *types.BasicRequest) (*pbMaintenance.DebugModeStatusResponse, error) {
	enabled, timestamp, err := s.debugModeRegistry.Status(ctx)
	if err != nil {
		return nil, createInternalError(fmt.Errorf("failed to get status of debug mode registry: %w", err))
	}

	return &pbMaintenance.DebugModeStatusResponse{IsEnabled: enabled, DisableAtTimestamp: timestamp}, nil
}

func createInternalError(err error) error {
	logrus.Error(err, interErrMsg)
	return status.Errorf(codes.Internal, err.Error())
}

func noInheritCancel(_ context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
