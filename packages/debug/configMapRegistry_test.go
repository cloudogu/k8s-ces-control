package debug

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"testing"
	"time"
)

const testNamespace = "ecosystem"

var (
	testCtx               = context.TODO()
	noInheritedTestCtx, _ = noInheritCancel(testCtx)
)

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
		cmRegistry := NewConfigMapDebugModeRegistry(cesRegistryMock, nil, clientSetMock, testNamespace)

		// then
		require.NotNil(t, cmRegistry)
	})
}

func Test_configMapDebugModeRegistry_BackupDoguLogLevels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedDoguLogLevelRegistryStr := "dogua: DEBUG\ndogub: ERROR\n"
		cesRegistryMock := newMockCesRegistry(t)
		doguLogLevelRegistryMock := newMockDoguLogLevelRegistry(t)
		doguLogLevelRegistryMock.EXPECT().MarshalFromCesRegistryToString(testCtx).Return(expectedDoguLogLevelRegistryStr, nil)

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

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, doguLogLevelRegistry: doguLogLevelRegistryMock, cesRegistry: cesRegistryMock}

		// when
		err := sut.BackupDoguLogLevels(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error getting registry config map", func(t *testing.T) {
		// given
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, assert.AnError)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, namespace: testNamespace}

		// when
		err := sut.BackupDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to get config map %s/%s", testNamespace, "debug-mode-registry"))
	})

	t.Run("should return error on error checking if registry is enabled with empty data", func(t *testing.T) {
		// given
		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
		}
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, namespace: testNamespace}

		// when
		err := sut.BackupDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "registry debug-mode-registry is not initialized")
	})

	t.Run("should return error if registry is not enabled", func(t *testing.T) {
		// given
		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "false"},
		}
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, namespace: testNamespace}

		// when
		err := sut.BackupDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "registry is not enabled")
	})

	t.Run("should return error on error marshal from string ", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		doguLogLevelRegistryMock := newMockDoguLogLevelRegistry(t)
		doguLogLevelRegistryMock.EXPECT().MarshalFromCesRegistryToString(testCtx).Return("", assert.AnError)

		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true"},
		}

		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, doguLogLevelRegistry: doguLogLevelRegistryMock, cesRegistry: cesRegistryMock}

		// when
		err := sut.BackupDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to renew dogu log level registry")
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

	t.Run("should return error on error deleting configmap", func(t *testing.T) {
		// given
		oldBackOff := maxThirtySecondsBackoff
		maxThirtySecondsBackoff = wait.Backoff{
			Steps:    1,
			Duration: 1 * time.Second,
			Factor:   1.0,
		}
		defer func() { maxThirtySecondsBackoff = oldBackOff }()

		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Delete(testCtx, "debug-mode-registry", metav1.DeleteOptions{}).Return(assert.AnError)

		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock}

		// when
		err := sut.Disable(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to delete debug mode registry configmap debug-mode-registry")
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

	t.Run("should return error on get registry error", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: registryName, Namespace: testNamespace},
		}
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))
		configMapInterfaceMock.EXPECT().Create(testCtx, cm, metav1.CreateOptions{}).Return(nil, assert.AnError)
		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		err := sut.Enable(testCtx, 15)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to create config map %s/%s", testNamespace, "debug-mode-registry"))
	})
}

func Test_configMapDebugModeRegistry_RestoreDoguLogLevels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		doguLogLevelRegistryStr := "dogua: DEBUG\ndogub: ERROR\n"
		doguLogLevelRegistryMock := newMockDoguLogLevelRegistry(t)
		doguLogLevelRegistryMock.EXPECT().UnMarshalFromStringToCesRegistry(doguLogLevelRegistryStr).Return(nil)

		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "dogus": doguLogLevelRegistryStr},
		}
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := &configMapDebugModeRegistry{configMapInterface: configMapClientMock, doguLogLevelRegistry: doguLogLevelRegistryMock, cesRegistry: cesRegistryMock}

		// when
		err := sut.RestoreDoguLogLevels(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error getting registry config map", func(t *testing.T) {
		// given
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, assert.AnError)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, namespace: testNamespace}

		// when
		err := sut.RestoreDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to get config map %s/%s", testNamespace, "debug-mode-registry"))
	})

	t.Run("should return error on error checking if registry is enabled with empty data", func(t *testing.T) {
		// given
		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
		}
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, namespace: testNamespace}

		// when
		err := sut.RestoreDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "registry debug-mode-registry is not initialized")
	})

	t.Run("should return error if registry is not enabled", func(t *testing.T) {
		// given
		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "false"},
		}
		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, namespace: testNamespace}

		// when
		err := sut.RestoreDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "registry is not enabled")
	})

	t.Run("should return error on missing dogu log level key", func(t *testing.T) {
		// given
		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true"},
		}

		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock}

		// when
		err := sut.RestoreDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "missing registry key dogus")
	})

	t.Run("should return error on error unmarshal from registry", func(t *testing.T) {
		// given
		doguLogLevelRegistryMock := newMockDoguLogLevelRegistry(t)
		cesRegistryMock := newMockCesRegistry(t)
		doguLogLevelRegistryMock.EXPECT().UnMarshalFromStringToCesRegistry("dogua: test\ndogub test").Return(assert.AnError)

		configMapRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "dogus": "dogua: test\ndogub test"},
		}

		configMapClientMock := newMockConfigMapInterface(t)
		configMapClientMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(configMapRegistry, nil)

		sut := configMapDebugModeRegistry{configMapInterface: configMapClientMock, doguLogLevelRegistry: doguLogLevelRegistryMock, cesRegistry: cesRegistryMock}

		// when
		err := sut.RestoreDoguLogLevels(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to restore dogu log level")
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

	t.Run("should return false on missing registry config map", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))
		sut := &configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		enabled, timestamp, err := sut.Status(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, false, enabled)
		assert.Equal(t, int64(0), timestamp)
	})

	t.Run("should return error on error getting registry config map", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, assert.AnError)
		sut := &configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		_, _, err := sut.Status(testCtx)

		// then
		require.Error(t, err)
	})

	t.Run("should return error on error checking if registry is enabled", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		emptyDataRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
		}
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(emptyDataRegistry, nil)
		sut := &configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		_, _, err := sut.Status(testCtx)

		// then
		require.Error(t, err)
	})

	t.Run("should return error on error getting disabledAtTimestamp with wrong format", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		emptyDataRegistry := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "debug-mode-registry", Namespace: testNamespace},
			Data:       map[string]string{"enabled": "true", "disable-at-timestamp": "something else"},
		}
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(emptyDataRegistry, nil)
		sut := &configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		_, _, err := sut.Status(testCtx)

		// then
		require.Error(t, err)
	})
}

func Test_configMapDebugModeRegistry_createRegistryIfNotFound(t *testing.T) {
	t.Run("should create a new empty registry if no found", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: registryName, Namespace: testNamespace},
		}
		configMapInterfaceMock.EXPECT().Create(testCtx, cm, metav1.CreateOptions{}).Return(cm, nil)
		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		registry, err := sut.createRegistryIfNotFound(testCtx, false)

		// then
		require.NoError(t, err)
		assert.Equal(t, cm, registry)
	})

	t.Run("should return error on registry create error", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: registryName, Namespace: testNamespace},
		}
		configMapInterfaceMock.EXPECT().Create(testCtx, cm, metav1.CreateOptions{}).Return(nil, assert.AnError)
		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		_, err := sut.createRegistryIfNotFound(testCtx, false)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to create config map %s/%s", testNamespace, "debug-mode-registry"))
	})
}

func Test_configMapDebugModeRegistry_updateConfigMap(t *testing.T) {
	t.Run("should return error on registry create error", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: registryName, Namespace: testNamespace},
		}
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(nil, assert.AnError)
		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		err := sut.updateConfigMap(testCtx, cm)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to get config map %s/%s", testNamespace, "debug-mode-registry"))
	})

	t.Run("should return error on update error", func(t *testing.T) {
		// given
		configMapInterfaceMock := newMockConfigMapInterface(t)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: registryName, Namespace: testNamespace},
		}
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(cm, nil)
		configMapInterfaceMock.EXPECT().Update(testCtx, cm, metav1.UpdateOptions{}).Return(nil, assert.AnError)
		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		err := sut.updateConfigMap(testCtx, cm)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to update config map %s/%s", testNamespace, "debug-mode-registry"))
	})

	t.Run("should return error if conflict timout is reached", func(t *testing.T) {
		// given
		oldBackOff := maxThirtySecondsBackoff
		maxThirtySecondsBackoff = wait.Backoff{
			Steps:    1,
			Duration: 1 * time.Second,
			Factor:   1.0,
		}
		defer func() { maxThirtySecondsBackoff = oldBackOff }()

		configMapInterfaceMock := newMockConfigMapInterface(t)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: registryName, Namespace: testNamespace},
		}
		configMapInterfaceMock.EXPECT().Get(testCtx, "debug-mode-registry", metav1.GetOptions{}).Return(cm, nil)
		configMapInterfaceMock.EXPECT().Update(testCtx, cm, metav1.UpdateOptions{}).Return(nil, errors.NewConflict(schema.GroupResource{}, "", assert.AnError))
		sut := configMapDebugModeRegistry{configMapInterface: configMapInterfaceMock, namespace: testNamespace}

		// when
		err := sut.updateConfigMap(testCtx, cm)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to update config map %s/%s", testNamespace, "debug-mode-registry"))
	})
}

func Test_getDisableAtTimeStamp(t *testing.T) {
	t.Run("should return error on empty configmap", func(t *testing.T) {
		// when
		_, err := getDisableAtTimeStamp(&v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "test"}})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "registry test is not initialized")
	})

	t.Run("should return zero timestamp on missing key", func(t *testing.T) {
		// given
		configMap := v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "test"}, Data: map[string]string{}}

		// when
		timestamp, err := getDisableAtTimeStamp(&configMap)

		// then
		require.NoError(t, err)
		assert.Equal(t, int64(0), timestamp)
	})
}

func Test_isRegistryEnabled(t *testing.T) {
	t.Run("should return false if key is missing", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "test"}, Data: map[string]string{}}

		// when
		enabled, err := isRegistryEnabled(configMap)

		// then
		require.NoError(t, err)
		assert.False(t, enabled)
	})

	t.Run("should return error on wrong string format", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "test"}, Data: map[string]string{"enabled": "invalidBool"}}

		// when
		_, err := isRegistryEnabled(configMap)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse bool")
	})
}
