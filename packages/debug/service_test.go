package debug

import (
	"github.com/cloudogu/ces-control-api/generated/maintenance"
	debugModeV1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	"github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
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
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, nil)
		debugModeClientMock.EXPECT().Update(testCtx, debugMode, metav1.UpdateOptions{}).Return(debugMode, nil)

		// when
		_, err := sut.Disable(testCtx, nil)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error enable maintenance mode", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, assert.AnError)

		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Disable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "ERROR: failed to get debug-mode")
	})

}

func Test_defaultDebugModeService_Enable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		doguInterActorMock := newMockDoguInterActor(t)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)

		debugModeRegistryMock := newMockDebugModeRegistry(t)

		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, nil)
		debugModeClientMock.EXPECT().Update(testCtx, debugMode, metav1.UpdateOptions{}).Return(debugMode, nil)

		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, maintenanceModeSwitch: maintenanceModeSwitchMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.NoError(t, err)
	})
	t.Run("should return error on error enable maintenance mode", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		maintenanceModeSwitchMock := newMockMaintenanceModeSwitch(t)
		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, assert.AnError)

		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, maintenanceModeSwitch: maintenanceModeSwitchMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "ERROR: failed to get debug-mode")
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
