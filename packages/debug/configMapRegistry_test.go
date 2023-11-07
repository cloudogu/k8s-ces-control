package debug

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strconv"
	"testing"
)

const testNamespace = "ecosystem"

var testCtx = context.TODO()

func TestNewConfigMapDebugModeRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		clientSetMock := newMockClusterClientSet(t)

		// when
		cmRegistry := NewConfigMapDebugModeRegistry(cesRegistryMock, clientSetMock, testNamespace)

		// then
		require.NotNil(t, cmRegistry)
	})
}

func Test_configMapDebugModeRegistry_BackupDoguLogLevels(t *testing.T) {

}

func Test_configMapDebugModeRegistry_Disable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock}

		// when
		result := sut.Disable(testCtx)

		// then
		require.NoError(t, result)
	})

	t.Run("do not retry if config map is not found", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(errors.NewNotFound(schema.GroupResource{}, "debug-mode-registry"))

		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock}

		// when
		result := sut.Disable(testCtx)

		// then
		require.NoError(t, result)
	})

	t.Run("should delete configmap after error with retry", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(assert.AnError).Times(1)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(nil).Times(1)

		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock}

		// when
		result := sut.Disable(testCtx)

		// then
		require.NoError(t, result)
	})
}

func Test_configMapDebugModeRegistry_Enable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       nil,
		}
		expectedUpdatedConfigMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "disable-at-timestamp": strconv.FormatInt(123456789, 10)},
		}

		cesRegistryMock := newMockCesRegistry(t)
		doguLogeLevelRegistry := newMockDoguLogLevelRegistry(t)
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil).Times(2)
		configMapInterfaceMock.EXPECT().Update(testCtx, expectedUpdatedConfigMapRegistry, metav1.UpdateOptions{}).Return(nil, nil)

		sut := &configMapDebugModeRegistry{cesRegistry: cesRegistryMock, configMapInterface: configMapInterfaceMock, doguLogLevelRegistry: doguLogeLevelRegistry, namespace: testNamespace}

		// when
		result := sut.Enable(testCtx, 123456789)

		// then
		require.NoError(t, result)
	})
}

func Test_configMapDebugModeRegistry_RestoreDoguLogLevels(t *testing.T) {

}

func Test_configMapDebugModeRegistry_Status(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		registryCm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "disable-at-timestamp": strconv.FormatInt(123456789, 10)},
		}
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(registryCm, nil)

		sut := &configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		enabled, timestamp, err := sut.Status(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, true, enabled)
		assert.Equal(t, int64(123456789), timestamp)
	})
}
