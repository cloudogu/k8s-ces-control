package support

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSupportArchiveService(t *testing.T) {
	t.Run("should create query clock", func(t *testing.T) {
		// given
		k8sClientMock := mock.AnythingOfType("k8sClient")
		discoveryClientMock := mock.AnythingOfType("discoveryInterface")

		// when
		sut := NewSupportArchiveService(k8sClientMock, discoveryClientMock)

		// then
		require.NotNil(t, sut)
	})
}
