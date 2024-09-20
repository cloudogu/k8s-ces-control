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
		doguConfigA := config.CreateDoguConfig("dogua", config.Entries{})
		doguConfigB := config.CreateDoguConfig("dogub", config.Entries{})
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
		doguConfigA := config.CreateDoguConfig("dogua", config.Entries{})
		doguConfigB := config.CreateDoguConfig("dogub", config.Entries{})
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

	t.Run("should return multi error on error update configmaps", func(t *testing.T) {
		// given
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)
		doguConfigA := config.CreateDoguConfig("dogua", config.Entries{})
		doguConfigB := config.CreateDoguConfig("dogub", config.Entries{})
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguConfigRepositoryMock.EXPECT().Update(context.TODO(), mock.Anything).Return(config.DoguConfig{}, assert.AnError).Times(2)

		sut := &doguLogLevelYamlRegistryMap{
			doguConfigRepository: doguConfigRepositoryMock,
		}

		// when
		err := sut.UnMarshalFromStringToCesRegistry(context.TODO(), "dogua: ERROR\ndogub: \"\"\n")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to update dogu config for dogu dogua: assert.AnError general error for testing\nfailed to update dogu config for dogu dogub: assert.AnError general error for testing")
	})
}

func Test_doguLogLevelRegistry_MarshalFromCesRegistryToString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedRegistryStr := "dogua: ERROR\ndogub: INFO\n"
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)

		doguConfigA := config.CreateDoguConfig("dogua", config.Entries{"logging/root": "ERROR"})
		doguConfigB := config.CreateDoguConfig("dogub", config.Entries{"logging/root": "INFO"})

		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		doguDescriptionGetter := newMockDoguDescriptorGetter(t)
		doguDescriptionGetter.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:              doguDescriptionGetter,
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
		descriptorGetter := newMockDoguDescriptorGetter(t)
		descriptorGetter.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:             descriptorGetter,
			logLevelRegistryMap: map[string]string{},
		}

		// when
		_, err := sut.MarshalFromCesRegistryToString(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get all dogus")
	})

	t.Run("should write entry with empty string for dogus not containing a log level entry", func(t *testing.T) {
		// given
		expectedRegistryStr := "dogua: \"\"\ndogub: INFO\n"
		doguA := &core.Dogu{Name: "official/dogua"}
		doguB := &core.Dogu{Name: "official/dogub"}
		doguConfigRepositoryMock := newMockDoguConfigRepository(t)

		doguConfigA := config.CreateDoguConfig("dogua", config.Entries{})
		doguConfigB := config.CreateDoguConfig("dogub", config.Entries{"logging/root": "INFO"})

		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogua")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("dogub")).Return(doguConfigB, nil)
		descriptorGetter := newMockDoguDescriptorGetter(t)
		descriptorGetter.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{doguA, doguB}, nil)

		sut := &doguLogLevelYamlRegistryMap{
			doguReg:              descriptorGetter,
			doguConfigRepository: doguConfigRepositoryMock,
			logLevelRegistryMap:  map[string]string{},
		}

		// when
		result, err := sut.MarshalFromCesRegistryToString(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRegistryStr, result)
	})
}
