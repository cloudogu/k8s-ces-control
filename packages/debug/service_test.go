package debug

import (
	"context"
	"github.com/cloudogu/k8s-ces-control/generated/maintenance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDebugModeService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		globalConfigMock := newMockConfigurationContext(t)
		cesRegistryMock.EXPECT().GlobalConfig().Return(globalConfigMock)
		doguRegistryMock := newMockDoguRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)

		clientSetMock := newMockClusterClientSet(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		configMapClientMock := newMockConfigMapInterface(t)
		coreV1Mock.EXPECT().ConfigMaps(testNamespace).Return(configMapClientMock)

		// when
		service := NewDebugModeService(cesRegistryMock, clientSetMock, testNamespace)

		// then
		require.NotNil(t, service)
	})
}

func Test_debugModeService_Disable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(testCtx).Return(nil).Run(func(ctx context.Context) {
			doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(nil)
		})
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		sut := debugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, nil)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error enable maintenance mode", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Deactivating debug mode").Return(assert.AnError)

		sut := debugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to activate maintenance mode")
	})

	t.Run("should return error on error restore debug mode registry to ces registry", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(assert.AnError)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(assert.AnError)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to restore log levels to ces registry")
	})

	t.Run("should return error on error stop all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(testCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to stop all dogus")
	})

	t.Run("should return error on start all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(testCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to start all dogus")
	})

	t.Run("should return error on error disable debug mode registry", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(testCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(nil)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(assert.AnError)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to disable the debug mode registry")
	})
}

func Test_debugModeService_Enable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(testCtx).Return(nil).Run(func(ctx context.Context) {
			doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(nil)
		})
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus("DEBUG").Return(nil)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		sut := debugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error enable maintenance mode", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(assert.AnError)

		sut := debugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to activate maintenance mode")
	})

	t.Run("should return error on error enable debug mode registry", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(assert.AnError)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(assert.AnError)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to enable debug mode registry")
	})

	t.Run("should return error on error backup dogu log levels", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(assert.AnError)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to backup dogu log levels")
	})

	t.Run("should return error on error set debug log level", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus("DEBUG").Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to set dogu log levels to debug")
	})

	t.Run("should return error on error stop all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus("DEBUG").Return(nil)
		doguInterActorMock.EXPECT().StopAllDogus(testCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to stop all dogus")
	})

	t.Run("should return error on error start all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus("DEBUG").Return(nil)
		doguInterActorMock.EXPECT().StopAllDogus(testCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to start all dogus")
	})
}

func Test_debugModeService_Status(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Status(testCtx).Return(true, 15, nil)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock}

		// when
		response, err := sut.Status(context.TODO(), nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, true, response.IsEnabled)
		assert.Equal(t, int64(15), response.DisableAtTimestamp)
	})

	t.Run("should return error on status error", func(t *testing.T) {
		// given
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Status(testCtx).Return(false, 0, assert.AnError)
		sut := debugModeService{debugModeRegistry: debugModeRegistryMock}

		// when
		_, err := sut.Status(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to get status of debug mode registry")
	})
}
