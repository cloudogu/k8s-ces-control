package ssl

import (
	"context"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
)

var testCtx = context.Background()

//go:embed testdata/valid_server.crt
var validCertBytes []byte

//go:embed testdata/valid_server.key
var validCertKeyBytes []byte

func TestNewManager(t *testing.T) {
	t.Run("should not be empty", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)

		// when
		actual := NewManager(nil, globalConfigMock)

		// then
		assert.NotEmpty(t, actual)
		assert.NotNil(t, actual.certGenerator)
		assert.Equal(t, globalConfigMock, actual.globalRegistry)
	})
}

func Test_manager_GetCertificateCredentials(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		oldCurrentNamespace := config.CurrentNamespace
		config.CurrentNamespace = "test-namespace"
		defer func() { config.CurrentNamespace = oldCurrentNamespace }()
		// given
		fakeClient := fake.NewSimpleClientset()
		clientMock := &mockClusterClient{fakeClient}
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return(string(validCertBytes), nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return(string(validCertKeyBytes), nil)
		sut := &manager{
			globalRegistry: registryMock,
			client:         clientMock,
		}

		// when
		cert, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.NoError(t, err)
		assert.NotNil(t, cert)
		// retrieve created secret and test it
		var getOpts metav1.GetOptions
		actual, err := fakeClient.CoreV1().Secrets("test-namespace").Get(testCtx, secretName, getOpts)
		require.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Equal(t, v1.SecretTypeTLS, actual.Type)
		assert.Equal(t, string(validCertBytes), actual.StringData["tls.crt"])
		assert.Equal(t, string(validCertKeyBytes), actual.StringData["tls.key"])
	})

	t.Run("should return error on certificate creation error", func(t *testing.T) {
		// given
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return(string(validCertBytes), nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return("", assert.AnError)
		sut := &manager{
			globalRegistry: registryMock,
		}

		// when
		_, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to create cert from registry")
	})

	t.Run("should fail to get certificate", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("", assert.AnError)
		sut := &manager{globalRegistry: globalConfigMock}

		// when
		actual, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check if certificate exists")
	})
	t.Run("should fail on cert generation", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		globalConfigMock := newMockConfigurationContext(t)
		notFoundErr := client.Error{Code: client.ErrorCodeKeyNotFound}
		globalConfigMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("", notFoundErr)

		certGenMock := newMockSslGenerator(t)
		certGenMock.EXPECT().GenerateSelfSignedCert(
			"k8s-ces-control",
			"k8s-ces-control",
			24000,
			"DE",
			"Lower Saxony",
			"Brunswick",
			[]string{"k8s-ces-control.ecosystem.svc.cluster.local", "localhost"},
		).Return("", "", assert.AnError)

		sut := &manager{
			globalRegistry: globalConfigMock,
			certGenerator:  certGenMock,
		}

		// when
		actual, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to generate self-signed certificate")
	})
	t.Run("should fail to set certificate in registry", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("", nil)
		globalConfigMock.EXPECT().Set("certificate/k8s-ces-control/server.crt", "some-cert").Return(assert.AnError)

		certGenMock := newMockSslGenerator(t)
		certGenMock.EXPECT().GenerateSelfSignedCert(
			"k8s-ces-control",
			"k8s-ces-control",
			24000,
			"DE",
			"Lower Saxony",
			"Brunswick",
			[]string{"k8s-ces-control.ecosystem.svc.cluster.local", "localhost"},
		).Return("some-cert", "some-key", nil)

		sut := &manager{
			globalRegistry: globalConfigMock,
			certGenerator:  certGenMock,
		}

		// when
		actual, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set certificate in registry")
	})
	t.Run("should fail to set legacy certificate in registry", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("", nil)
		globalConfigMock.EXPECT().Set("certificate/k8s-ces-control/server.crt", "some-cert").Return(nil)
		globalConfigMock.EXPECT().Set("certificate/cesappd/server.crt", "some-cert").Return(assert.AnError)

		certGenMock := newMockSslGenerator(t)
		certGenMock.EXPECT().GenerateSelfSignedCert(
			"k8s-ces-control",
			"k8s-ces-control",
			24000,
			"DE",
			"Lower Saxony",
			"Brunswick",
			[]string{"k8s-ces-control.ecosystem.svc.cluster.local", "localhost"},
		).Return("some-cert", "some-key", nil)

		sut := &manager{
			globalRegistry: globalConfigMock,
			certGenerator:  certGenMock,
		}

		// when
		actual, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set certificate in registry legacy location")
	})
	t.Run("should fail to set certificate key in registry", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("", nil)
		globalConfigMock.EXPECT().Set("certificate/k8s-ces-control/server.crt", "some-cert").Return(nil)
		globalConfigMock.EXPECT().Set("certificate/cesappd/server.crt", "some-cert").Return(nil)
		globalConfigMock.EXPECT().Set("certificate/k8s-ces-control/server.key", "some-key").Return(assert.AnError)

		certGenMock := newMockSslGenerator(t)
		certGenMock.EXPECT().GenerateSelfSignedCert(
			"k8s-ces-control",
			"k8s-ces-control",
			24000,
			"DE",
			"Lower Saxony",
			"Brunswick",
			[]string{"k8s-ces-control.ecosystem.svc.cluster.local", "localhost"},
		).Return("some-cert", "some-key", nil)

		sut := &manager{
			globalRegistry: globalConfigMock,
			certGenerator:  certGenMock,
		}

		// when
		actual, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set certificate key in registry")
	})
	// TODO write tests for certificate creation
}

func Test_manager_createCertFromRegistry(t *testing.T) {
	t.Run("should return error on registry cert read error", func(t *testing.T) {
		// given
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("", assert.AnError)
		sut := &manager{
			globalRegistry: registryMock,
		}

		// when
		_, _, _, err := sut.createCertFromRegistry()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error on registry key read error", func(t *testing.T) {
		// given
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("server.crt", nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return("", assert.AnError)
		sut := &manager{
			globalRegistry: registryMock,
		}

		// when
		_, _, _, err := sut.createCertFromRegistry()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error on certificate creation with invalid pem data", func(t *testing.T) {
		// given
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("server.crt", nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return("server.key", nil)
		sut := &manager{
			globalRegistry: registryMock,
		}

		// when
		_, _, _, err := sut.createCertFromRegistry()

		// then
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		// given
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return(string(validCertBytes), nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return(string(validCertKeyBytes), nil)
		sut := &manager{
			globalRegistry: registryMock,
		}

		// when
		certPEM, keyPEM, tlsCert, err := sut.createCertFromRegistry()

		// then
		require.NoError(t, err)
		assert.NotNil(t, tlsCert)
		assert.Equal(t, validCertBytes, []byte(certPEM))
		assert.Equal(t, validCertKeyBytes, []byte(keyPEM))
	})
}

type mockClusterClient struct {
	*fake.Clientset
}

func (m *mockClusterClient) Dogus(namespace string) ecoSystem.DoguInterface {
	// TODO implement me
	panic("implement me")
}
