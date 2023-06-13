package logging

import (
	"archive/zip"
	"bytes"
	"io"
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

func TestNewLoggingService(t *testing.T) {
	t.Run("should create query clock", func(t *testing.T) {
		// given

		// when
		sut := NewLoggingService(nil)

		// then
		require.NotNil(t, sut)
		assert.NotNil(t, sut.clock)
	})
}

func Test_realClock_Now(t *testing.T) {
	sut := new(realClock)
	actual := sut.Now()
	assert.IsType(t, actual, time.Now())
}

func Test_compressMessages(t *testing.T) {
	t.Run("decompressed message should be equal to the input message", func(t *testing.T) {
		// given
		const unicodeText = `
The ASCII compatible UTF-8 encoding of ISO 10646 and Unicode
plain-text files is defined in RFC 2279 and in ISO 10646-1 Annex R.

Using Unicode/UTF-8, you can write in emails and source code things such as

Mathematics and Sciences:
  ∮ E⋅da = Q,  n → ∞, ∑ f(i) = ∏ g(i), ∀x∈ℝ: ⌈x⌉ = −⌊−x⌋, α ∧ ¬β = ¬(¬α ∨ β),
  ℕ ⊆ ℕ₀ ⊂ ℤ ⊂ ℚ ⊂ ℝ ⊂ ℂ, ⊥ < a ≠ b ≡ c ≤ d ≪ ⊤ ⇒ (A ⇔ B),
  2H₂ + O₂ ⇌ 2H₂O, R = 4.7 kΩ, ⌀ 200 mm

Linguistics and dictionaries:
  ði ıntəˈnæʃənəl fəˈnɛtık əsoʊsiˈeıʃn
  Y [ˈʏpsilɔn], Yen [jɛn], Yoga [ˈjoːgɑ]

APL:
  ((V⍳V)=⍳⍴V)/V←,V    ⌷←⍳→⍴∆∇⊃‾⍎⍕⌈
`
		input := []byte(unicodeText)

		// when
		actual, err := compressMessages("my-dogu", input)

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
		zipreader, err := zip.NewReader(bytes.NewReader(actual), int64(len(actual)))
		require.NoError(t, err)
		for _, zipfile := range zipreader.File {
			assert.Equal(t, "my-dogu.log", zipfile.Name)
			fc, err := zipfile.Open()
			require.NoError(t, err)
			defer fc.Close()
			actualFileContent, err := io.ReadAll(fc)
			assert.Equal(t, []byte(unicodeText), actualFileContent)
		}
	})
}
