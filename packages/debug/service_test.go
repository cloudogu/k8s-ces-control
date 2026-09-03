package debug

import (
	"testing"
	"time"

	"github.com/cloudogu/ces-control-api/generated/maintenance"
	debugModeV1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	"github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
		service := NewDebugModeService(debugModeClientMock, doguInterActorMock, repository.DoguConfigRepository{}, doguDescriptionGetterMock, clientSetMock, testNamespace)

		// then
		require.NotNil(t, service)
	})
}

func Test_defaultDebugModeService_Disable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		doguInterActorMock := newMockDoguInterActor(t)
		debugModeRegistryMock := newMockDebugModeRegistry(t)
		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

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
		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, assert.AnError)

		sut := defaultDebugModeService{debugModeClient: debugModeClientMock}

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

		debugModeRegistryMock := newMockDebugModeRegistry(t)

		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, nil)
		debugModeClientMock.EXPECT().Update(testCtx, debugMode, metav1.UpdateOptions{}).Return(debugMode, nil)

		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.NoError(t, err)
	})

	t.Run("should return and do not call update if the debug mode is not existent before", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		doguInterActorMock := newMockDoguInterActor(t)

		debugModeRegistryMock := newMockDebugModeRegistry(t)

		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, "debug-mode"))
		debugModeClientMock.EXPECT().Create(testCtx, mock.Anything, metav1.CreateOptions{}).Return(nil, nil)

		sut := defaultDebugModeService{debugModeClient: debugModeClientMock, debugModeRegistry: debugModeRegistryMock, doguInterActor: doguInterActorMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{WithMaintenanceMode: true, Timer: 15})

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error enable maintenance mode", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		debugMode := &debugModeV1.DebugMode{}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, assert.AnError)

		sut := defaultDebugModeService{debugModeClient: debugModeClientMock}

		// when
		_, err := sut.Enable(testCtx, &maintenance.ToggleDebugModeRequest{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "ERROR: failed to get debug-mode")
	})
}

func Test_defaultDebugModeService_Status(t *testing.T) {
	deactivateTimestamp := metav1.NewTime(time.Date(2026, time.September, 3, 12, 34, 56, 0, time.UTC))

	t.Run("should return active debug mode with failed condition message", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		debugMode := &debugModeV1.DebugMode{
			Spec: debugModeV1.DebugModeSpec{DeactivateTimestamp: deactivateTimestamp},
			Status: debugModeV1.DebugModeStatus{
				Phase: debugModeV1.DebugModeStatusSet,
				Conditions: []metav1.Condition{{
					Type:    debugModeV1.ConditionFailed,
					Status:  metav1.ConditionTrue,
					Message: "failed to set log level",
				}},
			},
		}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, nil)
		sut := defaultDebugModeService{debugModeClient: debugModeClientMock}

		// when
		response, err := sut.Status(testCtx, nil)

		// then
		require.NoError(t, err)
		assert.True(t, response.IsEnabled)
		assert.Equal(t, deactivateTimestamp.UnixMilli(), response.DisableAtTimestamp)
		assert.Equal(t, "failed to set log level", response.Error)
	})

	t.Run("should return completed debug mode without error for false failed condition", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		debugMode := &debugModeV1.DebugMode{
			Spec: debugModeV1.DebugModeSpec{DeactivateTimestamp: deactivateTimestamp},
			Status: debugModeV1.DebugModeStatus{
				Phase: debugModeV1.DebugModeStatusCompleted,
				Conditions: []metav1.Condition{{
					Type:    debugModeV1.ConditionFailed,
					Status:  metav1.ConditionFalse,
					Message: "obsolete error",
				}},
			},
		}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, nil)
		sut := defaultDebugModeService{debugModeClient: debugModeClientMock}

		// when
		response, err := sut.Status(testCtx, nil)

		// then
		require.NoError(t, err)
		assert.False(t, response.IsEnabled)
		assert.Equal(t, deactivateTimestamp.UnixMilli(), response.DisableAtTimestamp)
		assert.Empty(t, response.Error)
	})

	t.Run("should return no error if failed condition is missing", func(t *testing.T) {
		// given
		debugModeClientMock := newMockDebugModeInterface(t)
		debugMode := &debugModeV1.DebugMode{
			Spec:   debugModeV1.DebugModeSpec{DeactivateTimestamp: deactivateTimestamp},
			Status: debugModeV1.DebugModeStatus{Phase: debugModeV1.DebugModeStatusRollback},
		}
		debugModeClientMock.EXPECT().Get(testCtx, "debug-mode", metav1.GetOptions{}).Return(debugMode, nil)
		sut := defaultDebugModeService{debugModeClient: debugModeClientMock}

		// when
		response, err := sut.Status(testCtx, nil)

		// then
		require.NoError(t, err)
		assert.True(t, response.IsEnabled)
		assert.Equal(t, deactivateTimestamp.UnixMilli(), response.DisableAtTimestamp)
		assert.Empty(t, response.Error)
	})
}
