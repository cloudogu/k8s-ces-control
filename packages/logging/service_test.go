package logging

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_buildLokiQueryUrl(t *testing.T) {
	t.Run("should return Loki LogQL query", func(t *testing.T) {
		// fix the time to the value 6795969378871345152, see below
		mockClock := newMockNowClock(t)

		// when
		actual, err := buildLokiQueryUrl("le-pod", 1, mockClock)

		// then
		require.NoError(t, err)
		assert.Equal(t, "http://loki-gateway.monitoring.svc.cluster.local:80/loki/api/v1/query_range?direction=backward&query=%7Bpod%3D~%22le-pod.%2A%22%7D&limit=1&start=-6795969378871345152", actual)
	})
}

type mockNowClock struct {
	mock.Mock
}

func newMockNowClock(t *testing.T) *mockNowClock {
	mock := &mockNowClock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

func (m *mockNowClock) Now() time.Time {
	parse, _ := time.Parse(time.RFC3339, "2022-02-22T22:22:22Z0000")
	return parse
}
