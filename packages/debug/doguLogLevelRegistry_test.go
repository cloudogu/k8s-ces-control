package debug

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDoguLogLevelRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		registryMock := newMockCesRegistry(t)

		// when
		debugReg := NewDoguLogLevelRegistryMap(registryMock)

		// then
		require.NotNil(t, debugReg)
		assert.Equal(t, map[string]string{}, debugReg.registry)
		assert.Equal(t, registryMock, debugReg.cesRegistry)
	})
}

func TestUnMarshalFromString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedRegistryMap := map[string]string{"dogua": "ERROR", "dogub": "INFO"}
		cesRegistryMock := newMockCesRegistry(t)

		sut := &doguLogLevelYamlRegistryMap{cesRegistry: cesRegistryMock}

		// when
		result, err := sut.UnMarshalFromString("dogua: ERROR\ndogub: INFO\n")

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryMap, result.registry)
	})
}

func Test_doguLogLevelRegistry_MarshalToString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedRegistryStr := "dogua: ERROR\ndogub: INFO\n"
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguAConfigMock := newMockDoguConfigurationContext(t)
		doguBConfigMock := newMockDoguConfigurationContext(t)
		doguRegistryMock := newMockDoguRegistry(t)
		cesRegistryMock := newMockCesRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)
		doguRegistryMock.EXPECT().GetAll().Return([]*core.Dogu{doguA, doguB}, nil)
		doguAConfigMock.EXPECT().Get("logging/root").Return("ERROR", nil)
		doguBConfigMock.EXPECT().Get("logging/root").Return("INFO", nil)

		sut := &doguLogLevelYamlRegistryMap{cesRegistry: cesRegistryMock, registry: map[string]string{}}

		// when
		result, err := sut.MarshalToString()

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryStr, result)
	})
}

func Test_doguLogLevelRegistry_RestoreToCesRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		doguAConfigMock := newMockDoguConfigurationContext(t)
		doguBConfigMock := newMockDoguConfigurationContext(t)
		doguAConfigMock.EXPECT().Set("logging/root", "ERROR").Return(nil)
		doguBConfigMock.EXPECT().Set("logging/root", "INFO").Return(nil)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)

		sut := &doguLogLevelYamlRegistryMap{cesRegistry: cesRegistryMock, registry: map[string]string{"dogua": "ERROR", "dogub": "INFO"}}

		// when
		result := sut.RestoreToCesRegistry()

		// then
		require.NoError(t, result)
	})
}
