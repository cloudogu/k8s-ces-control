package debug

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_defaultMaintenanceModeSwitch_ActivateMaintenanceMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Set("maintenance", "{\"title\":\"title\",\"text\":\"text\"}").Return(nil)
		sut := defaultMaintenanceModeSwitch{globalConfig: globalConfigMock}

		// when
		err := sut.ActivateMaintenanceMode("title", "text")

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error setting maintenance", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Set("maintenance", "{\"title\":\"title\",\"text\":\"text\"}").Return(assert.AnError)
		sut := defaultMaintenanceModeSwitch{globalConfig: globalConfigMock}

		// when
		err := sut.ActivateMaintenanceMode("title", "text")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set value [{\"title\":\"title\",\"text\":\"text\"}] with key maintenance")
	})
}

func Test_defaultMaintenanceModeSwitch_DeactivateMaintenanceMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Delete("maintenance").Return(nil)
		sut := defaultMaintenanceModeSwitch{globalConfig: globalConfigMock}

		// when
		err := sut.DeactivateMaintenanceMode()

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Delete("maintenance").Return(assert.AnError)
		sut := defaultMaintenanceModeSwitch{globalConfig: globalConfigMock}

		// when
		err := sut.DeactivateMaintenanceMode()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to delete key maintenance")
	})
}
