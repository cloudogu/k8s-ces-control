package logging

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func Test_buildLokiQueryUrl(t *testing.T) {
	t.Run("should return Loki LogQL query", func(t *testing.T) {
		// fix the time to the value 6795969378871345152, see below
		mockClock := newMockNowClock(t)
		parsedTime, _ := time.Parse(time.RFC3339, "2022-02-22T22:22:22Z0000")
		mockClock.EXPECT().Now().Return(parsedTime)

		// when
		actual, err := buildLokiQueryUrl("le-pod", 1, mockClock)

		// then
		require.NoError(t, err)
		assert.Equal(t, "http://loki-gateway.monitoring.svc.cluster.local:80/loki/api/v1/query_range?direction=backward&query=%7Bpod%3D~%22le-pod.%2A%22%7D&limit=1&start=-6795969378871345152", actual)
	})
}

func Test_doLokiHttpQuery(t *testing.T) {
	t.Run("should successfully call loki server with correct credentials", func(t *testing.T) {
		// given
		secret := map[string][]byte{"username": []byte("admin"), "password": []byte("admin123")}
		mockK8sSecretGet := NewMockSecretInterface(t)
		mockK8sSecretGet.EXPECT().Get(mock.Anything, "loki-credentials", mock.Anything).Return(&v1.Secret{Data: secret}, nil)
		mockCoreV1 := NewMockCoreV1Interface(t)
		mockCoreV1.EXPECT().Secrets("monitoring").Return(mockK8sSecretGet)
		mockClient := newMockClusterClient(t)
		mockClient.EXPECT().CoreV1().Return(mockCoreV1)
		mockLokiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			assert.True(t, ok)
		}))
		defer mockLokiServer.Close()

		// when
		actual, err := doLokiHttpQuery(mockClient, mockLokiServer.URL)

		// then
		require.NoError(t, err)
		assert.Equal(t, "200 OK", actual.Status)
	})
	t.Run("should error on missing secret", func(t *testing.T) {
		// given
		mockK8sSecretGet := NewMockSecretInterface(t)
		mockK8sSecretGet.EXPECT().Get(mock.Anything, "loki-credentials", mock.Anything).Return(nil, assert.AnError)
		mockCoreV1 := NewMockCoreV1Interface(t)
		mockCoreV1.EXPECT().Secrets("monitoring").Return(mockK8sSecretGet)
		mockClient := newMockClusterClient(t)
		mockClient.EXPECT().CoreV1().Return(mockCoreV1)
		mockLokiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			assert.True(t, ok)
		}))
		defer mockLokiServer.Close()

		// when
		_, err := doLokiHttpQuery(mockClient, mockLokiServer.URL)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, assert.AnError.Error())
	})
	t.Run("should succeed on loki HTTP error", func(t *testing.T) {
		// given
		mockK8sSecretGet := NewMockSecretInterface(t)
		mockK8sSecretGet.EXPECT().Get(mock.Anything, "loki-credentials", mock.Anything).Return(&v1.Secret{}, nil)
		mockCoreV1 := NewMockCoreV1Interface(t)
		mockCoreV1.EXPECT().Secrets("monitoring").Return(mockK8sSecretGet)
		mockClient := newMockClusterClient(t)
		mockClient.EXPECT().CoreV1().Return(mockCoreV1)
		mockLokiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer mockLokiServer.Close()

		// when
		actual, err := doLokiHttpQuery(mockClient, mockLokiServer.URL)

		// then
		require.NoError(t, err)
		assert.Equal(t, "400 Bad Request", actual.Status)
	})
}
