package debug

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-ces-control/packages/doguinteraction"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbMaintenance "github.com/cloudogu/k8s-ces-control/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/generated/types"
)

const (
	maintenanceTitle          = "Service unavailable"
	activateMaintenanceText   = "Activating debug mode"
	deactivateMaintenanceText = "Deactivating debug mode"
	logLevelDebug             = "DEBUG"
	interErrMsg               = "internal error"
)

type debugModeService struct {
	pbMaintenance.UnimplementedDebugModeServer
	globalConfig          configurationContext
	doguConfig            doguRegistry
	clientSet             clusterClientSet
	debugModeRegistry     debugModeRegistry
	maintenanceModeSwitch maintenanceModeSwitch
	namespace             string
	doguInterActor        doguInterActor
}

// NewDebugModeService returns an instance of debugModeService.
func NewDebugModeService(registry cesRegistry, clusterClient clusterClientSet, namespace string) *debugModeService {
	cmDebugModeRegistry := NewConfigMapDebugModeRegistry(registry, clusterClient, namespace)
	globalConfig := registry.GlobalConfig()
	watcher := NewDefaultConfigMapRegistryWatcher(clusterClient.CoreV1().ConfigMaps(namespace), cmDebugModeRegistry)
	watcher.StartWatch(context.Background())
	return &debugModeService{
		globalConfig:          globalConfig,
		doguConfig:            registry.DoguRegistry(),
		clientSet:             clusterClient,
		debugModeRegistry:     cmDebugModeRegistry,
		maintenanceModeSwitch: NewDefaultMaintenanceModeSwitch(globalConfig),
		namespace:             namespace,
		doguInterActor:        doguinteraction.NewDefaultDoguInterActor(clusterClient, namespace),
	}
}

// Enable enables the debug mode, sets dogu log level to debug and restarts all dogus.
func (s *debugModeService) Enable(ctx context.Context, req *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	err := s.maintenanceModeSwitch.ActivateMaintenanceMode(maintenanceTitle, activateMaintenanceText)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to activate maintenance mode: %w", err))
	}

	defer func() {
		err = s.maintenanceModeSwitch.DeactivateMaintenanceMode()
		if err != nil {
			log.FromContext(ctx).Error(fmt.Errorf("failed to deactivate maintenance mode: %w", err), interErrMsg)
		}
	}()

	err = s.debugModeRegistry.Enable(ctx, req.Timer)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to enable debug mode registry: %w", err))
	}

	err = s.debugModeRegistry.BackupDoguLogLevels(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to backup dogu log levels: %w", err))
	}

	err = s.doguInterActor.SetLogLevelInAllDogus(logLevelDebug)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to set dogu log levels to debug: %w", err))
	}

	// Create new context because the admin dogu itself will be canceled
	ctx = context.Background()

	err = s.doguInterActor.StopAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to stop all dogus: %w", err))
	}

	err = s.doguInterActor.StartAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to start all dogus %w", err))
	}

	return &types.BasicResponse{}, nil
}

// Disable returns an error because the method is unimplemented.
func (s *debugModeService) Disable(ctx context.Context, _ *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	err := s.maintenanceModeSwitch.ActivateMaintenanceMode(maintenanceTitle, deactivateMaintenanceText)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to activate maintenance mode: %w", err))
	}

	defer func() {
		err = s.maintenanceModeSwitch.DeactivateMaintenanceMode()
		if err != nil {
			log.FromContext(ctx).Error(fmt.Errorf("failed to deactivate maintenance mode: %w", err), interErrMsg)
		}
	}()

	err = s.debugModeRegistry.RestoreDoguLogLevels(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to restore log levels to ces registry: %w", err))
	}

	// Create new context because the admin dogu itself will be canceled
	ctx = context.Background()

	err = s.doguInterActor.StopAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to stop all dogus: %w", err))
	}

	err = s.doguInterActor.StartAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to start all dogus: %w", err))
	}

	err = s.debugModeRegistry.Disable(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to disable the debug mode registry: %w", err))
	}

	return &types.BasicResponse{}, nil
}

// Status return an error because the method is unimplemented.
func (s *debugModeService) Status(ctx context.Context, _ *types.BasicRequest) (*pbMaintenance.DebugModeStatusResponse, error) {
	enabled, timestamp, err := s.debugModeRegistry.Status(ctx)
	if err != nil {
		return nil, createInternalError(ctx, fmt.Errorf("failed to get status of debug mode registry: %w", err))
	}

	return &pbMaintenance.DebugModeStatusResponse{IsEnabled: enabled, DisableAtTimestamp: timestamp}, nil
}

func createInternalError(ctx context.Context, err error) error {
	logger := log.FromContext(ctx)
	logger.Error(err, interErrMsg)
	return status.Errorf(codes.Internal, err.Error())
}
