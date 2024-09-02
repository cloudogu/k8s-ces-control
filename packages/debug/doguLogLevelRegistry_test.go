package debug

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDoguLogLevelRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// when
		debugReg := NewDoguLogLevelRegistryMap(repository.DoguConfigRepository{}, nil)

		// then
		require.NotNil(t, debugReg)
		assert.Equal(t, map[string]string{}, debugReg.logLevelRegistryMap)
	})
}

func Test_doguLogLevelRegistry_UnMarshalFromStringToCesRegistry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)
		doguConfigA := config.DoguConfig{DoguName: config.SimpleDoguName("dogua")}
		doguConfigA.Config = config.CreateConfig(make(config.Entries))
		doguConfigB := config.DoguConfig{DoguName: config.SimpleDoguName("dogua")}
		doguConfigB.Config = config.CreateConfig(make(config.Entries))
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguConfigRepositoryMock.EXPECT().Update(context.TODO(), mock.Anything).Return(config.DoguConfig{}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguConfigRepository: doguConfigRepositoryMock,
		}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(context.TODO(), "dogua: ERROR\ndogub: INFO\n")

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on invalid registry string", func(t *testing.T) {
		// given
		notAMapString := "notayamlMap::;;!"
		sut := &doguLogLevelYamlRegistryMap{}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(context.TODO(), notAMapString)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to unmarshal dogu log level from string")
	})

	t.Run("should delete log level if the key is an empty string (restore to default level)", func(t *testing.T) {
		// given
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)
		doguConfigA := config.DoguConfig{DoguName: config.SimpleDoguName("dogua")}
		doguConfigA.Config = config.CreateConfig(make(config.Entries))
		doguConfigB := config.DoguConfig{DoguName: config.SimpleDoguName("dogua")}
		doguConfigB.Config = config.CreateConfig(make(config.Entries))
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguConfigRepositoryMock.EXPECT().Update(context.TODO(), mock.Anything).Return(config.DoguConfig{}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguConfigRepository: doguConfigRepositoryMock,
		}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(context.TODO(), "dogua: ERROR\ndogub: \"\"\n")

		// then
		require.NoError(t, err)
	})
}

func Test_doguLogLevelRegistry_MarshalFromCesRegistryToString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedRegistryStr := "dogua: ERROR\ndogub: INFO\n"
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)

		configA := config.CreateConfig(make(config.Entries))
		configA.Set("logging/root", "ERROR")
		doguConfigA := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configA}
		configB := config.CreateConfig(make(config.Entries))
		configB.Set("logging/root", "INFO")
		doguConfigB := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configB}

		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:              doguRegistryMock,
			doguConfigRepository: doguConfigRepositoryMock,
			logLevelRegistryMap:  map[string]string{},
		}

		// when
		result, err := sut.MarshalFromCesRegistryToString(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryStr, result)
	})

	t.Run("should return error on error getting all dogus", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:             doguRegistryMock,
			logLevelRegistryMap: map[string]string{},
		}

		// when
		_, err := sut.MarshalFromCesRegistryToString(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get all dogus")
	})

	t.Run("should return error on error getting logging key", func(t *testing.T) {
		//// given
		//doguA := &core.Dogu{Name: "official/dogua"}
		//doguB := &core.Dogu{Name: "official/dogub"}
		//doguAConfigMock := newMockConfigurationContext(t)
		//doguBConfigMock := newMockConfigurationContext(t)
		//doguRegistryMock := newMockDoguRegistry(t)
		//cesRegistryMock := newMockCesRegistry(t)
		//cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		//cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)
		//doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)
		//doguAConfigMock.EXPECT().Exists("logging/root").Return(false, assert.AnError)
		//doguBConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		//doguBConfigMock.EXPECT().Get("logging/root").Return("INFO", nil)

		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)

		configA := config.CreateConfig(make(config.Entries))
		configA.Set("logging/root", "ERROR")
		doguConfigA := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configA}
		configB := config.CreateConfig(make(config.Entries))
		configB.Set("logging/root", "INFO")
		doguConfigB := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configB}

		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:              doguRegistryMock,
			doguConfigRepository: doguConfigRepositoryMock,
			logLevelRegistryMap:  map[string]string{},
		}

		// when
		_, err := sut.MarshalFromCesRegistryToString(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "assert.AnError general error for testing")
	})

	t.Run("should write entry with empty string for dogus not containing a log level entry", func(t *testing.T) {
		// given
		//expectedRegistryStr := "dogua: \"\"\ndogub: INFO\n"
		//doguA := &core.Dogu{Name: "official/dogua"}
		//doguB := &core.Dogu{Name: "official/dogub"}
		//doguAConfigMock := newMockConfigurationContext(t)
		//doguBConfigMock := newMockConfigurationContext(t)
		//doguRegistryMock := newMockDoguRegistry(t)
		//cesRegistryMock := newMockCesRegistry(t)
		//cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		//cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)
		//doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)
		//doguAConfigMock.EXPECT().Exists("logging/root").Return(false, nil)
		//doguBConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		//doguBConfigMock.EXPECT().Get("logging/root").Return("INFO", nil)

		expectedRegistryStr := "dogua: \"\"\ndogub: INFO\n"
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)

		configA := config.CreateConfig(make(config.Entries))
		configA.Set("logging/root", "ERROR")
		doguConfigA := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configA}
		configB := config.CreateConfig(make(config.Entries))
		configB.Set("logging/root", "INFO")
		doguConfigB := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configB}

		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:              doguRegistryMock,
			doguConfigRepository: doguConfigRepositoryMock,
			logLevelRegistryMap:  map[string]string{},
		}

		// when
		result, err := sut.MarshalFromCesRegistryToString(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryStr, result)
	})

	t.Run("should return error on error getting dogu log level", func(t *testing.T) {
		// given
		//doguA := &core.Dogu{Name: "official/dogua"}
		//doguB := &core.Dogu{Name: "official/dogub"}
		//doguAConfigMock := newMockConfigurationContext(t)
		//doguBConfigMock := newMockConfigurationContext(t)
		//doguRegistryMock := newMockDoguRegistry(t)
		//cesRegistryMock := newMockCesRegistry(t)
		//cesRegistryMock.EXPECT().DoguConfig("dogua").Return(doguAConfigMock)
		//cesRegistryMock.EXPECT().DoguConfig("dogub").Return(doguBConfigMock)
		//doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)
		//doguAConfigMock.EXPECT().Exists("logging/root").Return(false, nil)
		//doguBConfigMock.EXPECT().Exists("logging/root").Return(true, nil)
		//doguBConfigMock.EXPECT().Get("logging/root").Return("", assert.AnError)

		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)

		configA := config.CreateConfig(make(config.Entries))
		configA.Set("logging/root", "ERROR")
		doguConfigA := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configA}
		configB := config.CreateConfig(make(config.Entries))
		configB.Set("logging/root", "INFO")
		doguConfigB := config.DoguConfig{DoguName: config.SimpleDoguName("dogua"), Config: configB}

		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:              doguRegistryMock,
			doguConfigRepository: doguConfigRepositoryMock,
			logLevelRegistryMap:  map[string]string{},
		}

		// when
		_, err := sut.MarshalFromCesRegistryToString(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "assert.AnError general error for testing")
	})
}
