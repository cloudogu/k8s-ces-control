package debug

import (
	"context"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_defaultMaintenanceModeSwitch_ActivateMaintenanceMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(config.CreateGlobalConfig(config.Entries{}), nil)
		globalConfigRepoMock.EXPECT().Update(testCtx, mock.Anything).RunAndReturn(func(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error) {
			get, b := globalConfig.Get("maintenance")
			require.True(t, b)
			assert.Equal(t, "{\"title\":\"title\",\"text\":\"text\"}", get.String())

			return globalConfig, nil
		})
		sut := defaultMaintenanceModeSwitch{globalConfigRepo: globalConfigRepoMock}

		// when
		err := sut.ActivateMaintenanceMode(testCtx, "title", "text")

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on getting global config", func(t *testing.T) {
		// given
		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(config.CreateGlobalConfig(config.Entries{}), assert.AnError)
		sut := defaultMaintenanceModeSwitch{globalConfigRepo: globalConfigRepoMock}

		// when
		err := sut.ActivateMaintenanceMode(testCtx, "title", "text")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get global config")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error on updating global config", func(t *testing.T) {
		// given
		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(config.CreateGlobalConfig(config.Entries{}), nil)
		globalConfigRepoMock.EXPECT().Update(testCtx, mock.Anything).RunAndReturn(func(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error) {
			get, b := globalConfig.Get("maintenance")
			require.True(t, b)
			assert.Equal(t, "{\"title\":\"title\",\"text\":\"text\"}", get.String())

			return globalConfig, assert.AnError
		})
		sut := defaultMaintenanceModeSwitch{globalConfigRepo: globalConfigRepoMock}

		// when
		err := sut.ActivateMaintenanceMode(testCtx, "title", "text")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed update global config for key \"maintenance\" value {\"title\":\"title\",\"text\":\"text\"}")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_defaultMaintenanceModeSwitch_DeactivateMaintenanceMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(config.CreateGlobalConfig(config.Entries{"maintenance": "value"}), nil)
		globalConfigRepoMock.EXPECT().Update(testCtx, mock.Anything).RunAndReturn(func(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error) {
			get, b := globalConfig.Get("maintenance")
			require.False(t, b)
			assert.Equal(t, "", get.String())

			return globalConfig, nil
		})
		sut := defaultMaintenanceModeSwitch{globalConfigRepo: globalConfigRepoMock}

		// when
		err := sut.DeactivateMaintenanceMode(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error getting global config", func(t *testing.T) {
		// given
		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(config.CreateGlobalConfig(config.Entries{"maintenance": "value"}), assert.AnError)
		sut := defaultMaintenanceModeSwitch{globalConfigRepo: globalConfigRepoMock}

		// when
		err := sut.DeactivateMaintenanceMode(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get globalConfig")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error on error updating global config", func(t *testing.T) {
		// given
		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(config.CreateGlobalConfig(config.Entries{"maintenance": "value"}), nil)
		globalConfigRepoMock.EXPECT().Update(testCtx, mock.Anything).RunAndReturn(func(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error) {
			get, b := globalConfig.Get("maintenance")
			require.False(t, b)
			assert.Equal(t, "", get.String())

			return globalConfig, assert.AnError
		})
		sut := defaultMaintenanceModeSwitch{globalConfigRepo: globalConfigRepoMock}

		// when
		err := sut.DeactivateMaintenanceMode(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to update global config for key \"maintenance\"")
		assert.ErrorIs(t, err, assert.AnError)
	})
}
