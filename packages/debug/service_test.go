package debug

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-ces-control/generated/maintenance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDebugModeService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		globalConfigMock := newMockDoguConfigurationContext(t)
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
		// The Run Method ensures the order
		doguInterActorMock.EXPECT().StopDoguWithWait(testCtx, "redmine", true).Return(nil).Run(func(_ context.Context, _ string, _ bool) {
			doguInterActorMock.EXPECT().StopDoguWithWait(testCtx, "postgresql", true).Return(nil)
		})
		doguInterActorMock.EXPECT().StartDoguWithWait(testCtx, "postgresql", true).Return(nil).Run(func(_ context.Context, _ string, _ bool) {
			doguInterActorMock.EXPECT().StartDoguWithWait(testCtx, "redmine", true).Return(nil)
		})
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetAll().Return([]*core.Dogu{{Name: "official/postgresql"}, {Name: "official/redmine", Dependencies: []core.Dependency{{Name: "postgresql"}}}}, nil)
		cesRegistryMock := newMockCesRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Deactivating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().RestoreDoguLogLevels(testCtx).Return(nil)
		debugModeRegistryMock.EXPECT().Disable(testCtx).Return(nil)
		sut := debugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, registry: cesRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Disable(testCtx, nil)

		// then
		require.NoError(t, err)
	})
}

func Test_debugModeService_Enable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		// The Run Method ensures the order
		doguInterActorMock.EXPECT().StopDoguWithWait(testCtx, "redmine", true).Return(nil).Run(func(_ context.Context, _ string, _ bool) {
			doguInterActorMock.EXPECT().StopDoguWithWait(testCtx, "postgresql", true).Return(nil)
		})
		doguInterActorMock.EXPECT().StartDoguWithWait(testCtx, "postgresql", true).Return(nil).Run(func(_ context.Context, _ string, _ bool) {
			doguInterActorMock.EXPECT().StartDoguWithWait(testCtx, "redmine", true).Return(nil)
		})
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetAll().Return([]*core.Dogu{{Name: "official/postgresql"}, {Name: "official/redmine", Dependencies: []core.Dependency{{Name: "postgresql"}}}}, nil)
		doguPsqlConfigMock := newMockDoguConfigurationContext(t)
		doguPsqlConfigMock.EXPECT().Set("logging/root", "DEBUG").Return(nil)
		doguRedmineConfigMock := newMockDoguConfigurationContext(t)
		doguRedmineConfigMock.EXPECT().Set("logging/root", "DEBUG").Return(nil)
		cesRegistryMock := newMockCesRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)
		cesRegistryMock.EXPECT().DoguConfig("postgresql").Return(doguPsqlConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("redmine").Return(doguRedmineConfigMock)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		maintenanceModeSwitchMock.EXPECT().ActivateMaintenanceMode("Service unavailable", "Activating debug mode").Return(nil)
		maintenanceModeSwitchMock.EXPECT().DeactivateMaintenanceMode().Return(nil)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		debugModeRegistryMock.EXPECT().Enable(testCtx, int32(15)).Return(nil)
		debugModeRegistryMock.EXPECT().BackupDoguLogLevels(testCtx).Return(nil)
		sut := debugModeService{maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, registry: cesRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.NoError(t, err)
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
}
