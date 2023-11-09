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
		cesRegistryMock := newMockCesRegistry(t)
		doguAConfigMock := newMockDoguConfigurationContext(t)
		doguBConfigMock := newMockDoguConfigurationContext(t)
		doguAConfigMock.EXPECT().Set("logging/root", "ERROR").Return(nil)
		doguBConfigMock.EXPECT().Set("logging/root", "INFO").Return(nil)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)

		sut := &doguLogLevelYamlRegistryMap{cesRegistry: cesRegistryMock}

		// when
		err := sut.UnMarshalFromStringToCesRegistry("dogua: ERROR\ndogub: INFO\n")

		// then
		require.NoError(t, err)
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
		doguAConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		doguAConfigMock.EXPECT().Get("logging/root").Return("ERROR", nil)
		doguBConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		doguBConfigMock.EXPECT().Get("logging/root").Return("INFO", nil)

		sut := &doguLogLevelYamlRegistryMap{cesRegistry: cesRegistryMock, registry: map[string]string{}}

		// when
		result, err := sut.MarshalFromCesRegistryToString()

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryStr, result)
	})
}
