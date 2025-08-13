package debug

import (
	"context"
	"github.com/cloudogu/ces-control-api/generated/maintenance"
	debugModeV1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	"github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewdefaultDebugModeService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)

		doguDescriptionGetterMock := newMockDoguDescriptorGetter(t)

		doguInterActorMock := newMockDoguInterActor(t)
		clientSetMock := newMockClusterClientSet(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		configMapClientMock := newMockConfigMapInterface(t)
		coreV1Mock.EXPECT().ConfigMaps(testNamespace).Return(configMapClientMock)

		// when
		service := NewDebugModeService(debugModeClientMock, doguInterActorMock, repository.DoguConfigRepository{}, repository.GlobalConfigRepository{}, doguDescriptionGetterMock, clientSetMock, testNamespace)

		// then
		require.NotNil(t, service)
	})
}

func Test_defaultDebugModeService_Disable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(nil).Run(func(ctx context.Context) {
			doguInterActorMock.EXPECT().StartAllDogus(noInheritedTestCtx).Return(nil)
		})
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(noInheritedTestCtx).Return(nil)
		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, nil)

		// when
		_, err := sut.Disable(testCtx, nil)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error enable maintenance mode", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Deactivating debug mode").Return(assert.AnError)

		sut := defaultDebugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to activate maintenance mode")
	})

	t.Run("should return error on error restore debug mode registry to ces registry", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(assert.AnError)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(assert.AnError)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to restore log levels to ces registry")
	})

	t.Run("should return error on error stop all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to stop all dogus")
	})

	t.Run("should return error on start all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(noInheritedTestCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to start all dogus")
	})

	t.Run("should return error on error disable debug mode registry", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(noInheritedTestCtx).Return(nil)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(noInheritedTestCtx).Return(assert.AnError)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to disable the debug mode registry")
	})
}

func Test_defaultDebugModeService_Enable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(nil).Run(func(ctx context.Context) {
			doguInterActorMock.EXPECT().StartAllDogus(noInheritedTestCtx).Return(nil)
		})
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus(testCtx, "DEBUG").Return(nil)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		sut := defaultDebugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error enable maintenance mode", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(assert.AnError)

		sut := defaultDebugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to activate maintenance mode")
	})

	t.Run("should return error on error enable debug mode registry", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(assert.AnError)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(assert.AnError)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to enable debug mode registry")
	})

	t.Run("should return wrapped error on error rollback enable debug mode registry", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(assert.AnError)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(assert.AnError)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(assert.AnError)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rollback error: assert.AnError general error for testing")
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to enable debug mode registry: assert.AnError general error for testing")
	})

	t.Run("should return error on error backup dogu log levels", func(t *testing.T) {
		// given
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(assert.AnError)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to backup dogu log levels")
	})

	t.Run("should return error on error set debug log level", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus(testCtx, "DEBUG").Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to set dogu log levels to debug")
	})

	t.Run("should return wrapped error on rollback error set debug log level", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus(testCtx, "DEBUG").Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(assert.AnError)
		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to set dogu log levels to debug")
		assert.ErrorContains(t, err, "rollback error: assert.AnError general error for testing")
	})

	t.Run("should return error on error stop all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus(testCtx, "DEBUG").Return(nil)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(nil)

		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to stop all dogus")
	})

	t.Run("should return wrapped error on error rollback stop all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus(testCtx, "DEBUG").Return(nil)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(assert.AnError)

		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to stop all dogus")
		assert.ErrorContains(t, err, "rollback error: assert.AnError general error for testing")
	})

	t.Run("should return error on error start all dogus", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().SetLogLevelInAllDogus(testCtx, "DEBUG").Return(nil)
		doguInterActorMock.EXPECT().StopAllDogus(noInheritedTestCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(noInheritedTestCtx).Return(assert.AnError)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode(context.TODO(), "Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode(mock.AnythingOfType("*context.cancelCtx")).Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		doguInterActorMock.EXPECT().StartAllDogus(testCtx).Return(nil)

		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock, maintenanceModeSwitch: maintenanceModeSwitchMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to start all dogus")
	})
}

// func Test_defaultDebugModeService_Status(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		// given
//		debugModeRegistryMock := newMockDebugModeRegistry(t)
//		debugModeRegistryMock.EXPECT().Status(testCtx).Return(true, 15, nil)
//		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock}
//
//		// when
//		response, err := sut.Status(context.TODO(), nil)
//
//		// then
//		require.NoError(t, err)
//		assert.Equal(t, true, response.IsEnabled)
//		assert.Equal(t, int64(15), response.DisableAtTimestamp)
//	})
//
//	t.Run("should return error on status error", func(t *testing.T) {
//		// given
//		debugModeRegistryMock := newMockDebugModeRegistry(t)
//		debugModeRegistryMock.EXPECT().Status(testCtx).Return(false, 0, assert.AnError)
//		sut := defaultDebugModeService{debugModeRegistry: debugModeRegistryMock}
//
//		// when
//		_, err := sut.Status(context.TODO(), nil)
//
//		// then
//		require.Error(t, err)
//		assert.ErrorContains(t, err, "rpc error: code = Internal desc = failed to get status of debug mode registry")
//	})
// }
