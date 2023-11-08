package debug

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
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
	doguConfigKeyLogLevel     = "logging/root"
	logLevelDebug             = "DEBUG"
	interErrMsg               = "internal error"
)

type debugModeService struct {
	pbMaintenance.UnimplementedDebugModeServer
	registry              cesRegistry
	globalConfig          doguConfigurationContext
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
	return &debugModeService{
		registry:              registry,
		globalConfig:          globalConfig,
		doguConfig:            registry.DoguRegistry(),
		clientSet:             clusterClient,
		debugModeRegistry:     cmDebugModeRegistry,
		maintenanceModeSwitch: NewDefaultMaintenanceModeSwitch(globalConfig),
		namespace:             namespace,
		doguInterActor:        doguinteraction.NewDefaultDoguInterActor(clusterClient, namespace),
	}
}

// Enable returns an error because the method is unimplemented.
// TODO Use timer to disable debug mode.
// TODO rotate logs on enable and disable
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

	// TODO Check if this is a timestamp or just 15Minutes in seconds
	err = s.debugModeRegistry.Enable(ctx, req.Timer)
	if err != nil {
		return nil, createInternalError(ctx, err)
	}

	err = s.debugModeRegistry.BackupDoguLogLevels(ctx)
	if err != nil {
		return nil, createInternalError(ctx, err)
	}

	err = s.setLogLevelInAllDogus(logLevelDebug)
	if err != nil {
		// TODO only log errors here?
		return nil, createInternalError(ctx, err)
	}

	// Create new context because the admin dogu itself will be canceled
	ctx = context.Background()

	err = s.stopAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, err)
	}

	err = s.startAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, err)
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
		return nil, createInternalError(ctx, err)
	}

	// Create new context because the admin dogu itself will be canceled
	ctx = context.Background()

	err = s.stopAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, err)
	}

	err = s.startAllDogus(ctx)
	if err != nil {
		return nil, createInternalError(ctx, err)
	}

	err = s.debugModeRegistry.Disable(ctx)
	if err != nil {
		return nil, createInternalError(ctx, err)
	}

	return &types.BasicResponse{}, nil
}

// Status return an error because the method is unimplemented.
func (s *debugModeService) Status(ctx context.Context, _ *types.BasicRequest) (*pbMaintenance.DebugModeStatusResponse, error) {
	enabled, timestamp, err := s.debugModeRegistry.Status(ctx)
	if err != nil {
		return nil, createInternalError(ctx, err)
	}

	return &pbMaintenance.DebugModeStatusResponse{IsEnabled: enabled, DisableAtTimestamp: timestamp}, nil
}

// TODO move this to doguinteraction? Or to debugmoderegistry?
func (s *debugModeService) setLogLevelInAllDogus(logLevel string) error {
	allDogus, err := s.registry.DoguRegistry().GetAll()
	if err != nil {
		return getAllDogusError(err)
	}

	var multiError error
	for _, dogu := range allDogus {
		doguConfig := s.registry.DoguConfig(dogu.GetSimpleName())
		setErr := doguConfig.Set(doguConfigKeyLogLevel, logLevel)
		if setErr != nil {
			multiError = errors.Join(multiError, setErr)
		}
	}

	return multiError
}

func (s *debugModeService) stopAllDogus(ctx context.Context) error {
	allDogus, err := s.registry.DoguRegistry().GetAll()
	if err != nil {
		return getAllDogusError(err)
	}

	sortedDogus := core.SortDogusByInvertedDependency(allDogus)
	var multiError error
	for _, dogu := range sortedDogus {
		stopErr := s.doguInterActor.StopDoguWithWait(ctx, dogu.GetSimpleName(), true)
		if stopErr != nil {
			multiError = errors.Join(multiError, stopErr)
		}
	}

	return multiError
}

func (s *debugModeService) startAllDogus(ctx context.Context) error {
	allDogus, err := s.registry.DoguRegistry().GetAll()
	if err != nil {
		return getAllDogusError(err)
	}

	sortedDogus := core.SortDogusByDependency(allDogus)
	var multiError error
	for _, dogu := range sortedDogus {
		startErr := s.doguInterActor.StartDoguWithWait(ctx, dogu.GetSimpleName(), true)
		if startErr != nil {
			multiError = errors.Join(multiError, startErr)
		}
	}

	return multiError
}

func getAllDogusError(err error) error {
	return fmt.Errorf("failed to get all dogus: %w", err)
}

func createInternalError(ctx context.Context, err error) error {
	logger := log.FromContext(ctx)
	logger.Error(err, interErrMsg)
	return status.Errorf(codes.Internal, err.Error())
}
