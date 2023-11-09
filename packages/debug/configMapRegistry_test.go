package debug

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
	"time"
)

const testNamespace = "ecosystem"

var testCtx = context.TODO()

func TestNewConfigMapDebugModeRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		clientSetMock := newMockClusterClientSet(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		configMapInterfaceMock := newMockConfigMapInterface(t)
		coreV1Mock.EXPECT().ConfigMaps(testNamespace).Return(configMapInterfaceMock)

		// when
		cmRegistry := NewConfigMapDebugModeRegistry(cesRegistryMock, clientSetMock, testNamespace)

		// then
		require.NotNil(t, cmRegistry)
	})
}

func Test_configMapDebugModeRegistry_BackupDoguLogLevels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedDoguLogLevelRegistryStr := "dogua: DEBUG\ndogub: ERROR\n"
		doguLogLevelRegistryMock := newMockDoguLogLevelRegistry(t)
		doguLogLevelRegistryMock.EXPECT().MarshalFromCesRegistryToString().Return(expectedDoguLogLevelRegistryStr, nil)

		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true"},
		}

		expectedUpdatedConfigMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "dogus": expectedDoguLogLevelRegistryStr},
			BinaryData: nil,
		}

		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)
		configMapClientMock.EXPECT().Update(testCtx, expectedUpdatedConfigMapRegistry, metav1.UpdateOptions{}).Return(nil, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, doguLogLevelRegistry: doguLogLevelRegistryMock}

		// when
		err := sut.BackupDoguLogLevels(testCtx)

		// then
		require.NoError(t, err)
	})
}

func Test_configMapDebugModeRegistry_Disable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock}

		// when
		err := sut.Disable(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("do not retry if config map is not found", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(errors.NewNotFound(schema.GroupResource{}, "debug-mode-registry"))

		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock}

		// when
		err := sut.Disable(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete configmap after error with retry", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(assert.AnError).Times(1)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(nil).Times(1)

		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock}

		// when
		err := sut.Disable(testCtx)

		// then
		require.NoError(t, err)
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
			Data:       map[string]string{"enabled": "true", "disable-at-timestamp": time.Now().Add(time.Minute * 15).Format(time.RFC822)},
		}

		cesRegistryMock := newMockCesRegistry(t)
		doguLogeLevelRegistry := newMockDoguLogLevelRegistry(t)
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil).Times(2)
		configMapInterfaceMock.EXPECT().Update(testCtx, expectedUpdatedConfigMapRegistry, metav1.UpdateOptions{}).Return(nil, nil)

		sut := &configMapDebugModeRegistry{cesRegistry: cesRegistryMock, configMapInterface: configMapInterfaceMock, doguLogLevelRegistry: doguLogeLevelRegistry, namespace: testNamespace}

		// when
		err := sut.Enable(testCtx, 15)

		// then
		require.NoError(t, err)
	})
}

func Test_configMapDebugModeRegistry_RestoreDoguLogLevels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguLogLevelRegistryStr := "dogua: DEBUG\ndogub: ERROR\n"
		doguLogLevelRegistryMock := newMockDoguLogLevelRegistry(t)
		doguLogLevelRegistryMock.EXPECT().UnMarshalFromStringToCesRegistry(doguLogLevelRegistryStr).Return(nil)

		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "dogus": doguLogLevelRegistryStr},
		}
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := &configMapDebugModeRegistry{configMapInterface: configMapClientMock, doguLogLevelRegistry: doguLogLevelRegistryMock}

		// when
		err := sut.RestoreDoguLogLevels(testCtx)

		// then
		require.NoError(t, err)
	})
}

func Test_configMapDebugModeRegistry_Status(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		now := time.Now()
		expectedFormat := now.Format(time.RFC822)
		registryCm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "disable-at-timestamp": expectedFormat},
		}
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(registryCm, nil)

		sut := &configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		enabled, timestamp, err := sut.Status(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, true, enabled)
		assert.Equal(t, expectedFormat, time.UnixMilli(timestamp).Format(time.RFC822))
	})
}
