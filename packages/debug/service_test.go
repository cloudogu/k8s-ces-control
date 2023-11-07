package debug

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDebugModeService(t *testing.T) {
	// service := NewDebugModeService()

	// require.NotNil(t, service)
}

func Test_debugModeService_Disable(t *testing.T) {
	t.Run("should throw error because it is not implemented", func(t *testing.T) {
		// given
		sut := debugModeService{}

		// when
		_, err := sut.Disable(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "this service is not yet implemented")
	})
}

func Test_debugModeService_Enable(t *testing.T) {
	t.Run("should throw error because it is not implemented", func(t *testing.T) {
		// given
		sut := debugModeService{}

		// when
		_, err := sut.Enable(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "this service is not yet implemented")
	})
}

func Test_debugModeService_Status(t *testing.T) {
	t.Run("should throw error because it is not implemented", func(t *testing.T) {
		// given
		sut := debugModeService{}

		// when
		_, err := sut.Status(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "this service is not yet implemented")
	})
}
