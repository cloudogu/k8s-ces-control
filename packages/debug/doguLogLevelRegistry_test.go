package debug

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDoguLogLevelRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// when
		debugReg := NewDoguLogLevelRegistryMap()

		// then
		require.NotNil(t, debugReg)
		assert.Equal(t, map[string]string{}, debugReg.registry)
	})
}

func Test_doguLogLevelRegistry_UnMarshalFromStringToCesRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
		doguAConfigMock.EXPECT().Set("logging/root", "ERROR").Return(nil)
		doguBConfigMock.EXPECT().Set("logging/root", "INFO").Return(nil)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)

		sut := &doguLogLevelYamlRegistryMap{}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(cesRegistryMock, "dogua: ERROR\ndogub: INFO\n")

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on invalid registry string", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		notAMapString := "notayamlMap::;;!"
		sut := &doguLogLevelYamlRegistryMap{}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(cesRegistryMock, notAMapString)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to unmarshal dogu log level from string")
	})

	t.Run("should delete log level if the key is an empty string (restore to default level)", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
		doguAConfigMock.EXPECT().Set("logging/root", "ERROR").Return(nil)
		doguBConfigMock.EXPECT().Delete("logging/root").Return(nil)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)

		sut := &doguLogLevelYamlRegistryMap{}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(cesRegistryMock, "dogua: ERROR\ndogub: \"\"\n")

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
		doguAConfigMock.EXPECT().Set("logging/root", "ERROR").Return(nil)
		doguBConfigMock.EXPECT().Delete("logging/root").Return(assert.AnError)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)

		sut := &doguLogLevelYamlRegistryMap{}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(cesRegistryMock, "dogua: ERROR\ndogub: \"\"\n")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "assert.AnError general error for testing")
	})

	t.Run("should return error on set log level error", func(t *testing.T) {
		// given
		cesRegistryMock := newMockCesRegistry(t)
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
		doguAConfigMock.EXPECT().Set("logging/root", "ERROR").Return(assert.AnError)
		doguBConfigMock.EXPECT().Delete("logging/root").Return(nil)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)

		sut := &doguLogLevelYamlRegistryMap{}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(cesRegistryMock, "dogua: ERROR\ndogub: \"\"\n")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "assert.AnError general error for testing")
	})
}

func Test_doguLogLevelRegistry_MarshalFromCesRegistryToString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedRegistryStr := "dogua: ERROR\ndogub: INFO\n"
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
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

		sut := &doguLogLevelYamlRegistryMap{registry: map[string]string{}}

		// when
		result, err := sut.MarshalFromCesRegistryToString(cesRegistryMock)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryStr, result)
	})

	t.Run("should return error on error getting all dogus", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguRegistry(t)
		cesRegistryMock := newMockCesRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)
		doguRegistryMock.EXPECT().GetAll().Return(nil, assert.AnError)

		sut := &doguLogLevelYamlRegistryMap{registry: map[string]string{}}

		// when
		_, err := sut.MarshalFromCesRegistryToString(cesRegistryMock)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get all dogus")
	})

	t.Run("should return error on error getting logging key", func(t *testing.T) {
		// given
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
		doguRegistryMock := newMockDoguRegistry(t)
		cesRegistryMock := newMockCesRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)
		doguRegistryMock.EXPECT().GetAll().Return([]*core.Dogu{doguA, doguB}, nil)
		doguAConfigMock.EXPECT().Exists("logging/root").Return(false, assert.AnError)
		doguBConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		doguBConfigMock.EXPECT().Get("logging/root").Return("INFO", nil)

		sut := &doguLogLevelYamlRegistryMap{registry: map[string]string{}}

		// when
		_, err := sut.MarshalFromCesRegistryToString(cesRegistryMock)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "assert.AnError general error for testing")
	})

	t.Run("should write entry with empty string for dogus not containing a log level entry", func(t *testing.T) {
		// given
		expectedRegistryStr := "dogua: \"\"\ndogub: INFO\n"
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
		doguRegistryMock := newMockDoguRegistry(t)
		cesRegistryMock := newMockCesRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)
		doguRegistryMock.EXPECT().GetAll().Return([]*core.Dogu{doguA, doguB}, nil)
		doguAConfigMock.EXPECT().Exists("logging/root").Return(false, nil)
		doguBConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		doguBConfigMock.EXPECT().Get("logging/root").Return("INFO", nil)

		sut := &doguLogLevelYamlRegistryMap{registry: map[string]string{}}

		// when
		result, err := sut.MarshalFromCesRegistryToString(cesRegistryMock)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryStr, result)
	})

	t.Run("should return error on error getting dogu log level", func(t *testing.T) {
		// given
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguAConfigMock := newMockConfigurationContext(t)
		doguBConfigMock := newMockConfigurationContext(t)
		doguRegistryMock := newMockDoguRegistry(t)
		cesRegistryMock := newMockCesRegistry(t)
		cesRegistryMock.EXPECT().DoguRegistry().Return(doguRegistryMock)
		cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)
		doguRegistryMock.EXPECT().GetAll().Return([]*core.Dogu{doguA, doguB}, nil)
		doguAConfigMock.EXPECT().Exists("logging/root").Return(false, nil)
		doguBConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		doguBConfigMock.EXPECT().Get("logging/root").Return("", assert.AnError)

		sut := &doguLogLevelYamlRegistryMap{registry: map[string]string{}}

		// when
		_, err := sut.MarshalFromCesRegistryToString(cesRegistryMock)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "assert.AnError general error for testing")
	})
}
