package ssl

import (
	"context"
	_ "embed"
	"k8s.io/apimachinery/pkg/runtime"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	faketesting "k8s.io/client-go/testing"
	"testing"

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

	t.Run("should fail because an error occur getting the secret", func(t *testing.T) {
		oldCurrentNamespace := config.CurrentNamespace
		config.CurrentNamespace = "test-namespace"
		defer func() { config.CurrentNamespace = oldCurrentNamespace }()
		// given
		fakeClient := fake.NewSimpleClientset()
		fakeClient.CoreV1().(*fakecorev1.FakeCoreV1).PrependReactor("get", "secrets", func(action faketesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, assert.AnError
		})
		clientMock := &mockClusterClient{fakeClient}
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return(string(validCertBytes), nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return(string(validCertKeyBytes), nil)
		sut := &manager{
			globalRegistry: registryMock,
			client:         clientMock,
		}

		// when
		_, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while checking whether certificate secret k8s-ces-control-server-certificate already exists")
	})

	t.Run("should fail because an error occur creating the secret", func(t *testing.T) {
		oldCurrentNamespace := config.CurrentNamespace
		config.CurrentNamespace = "test-namespace"
		defer func() { config.CurrentNamespace = oldCurrentNamespace }()
		// given
		fakeClient := fake.NewSimpleClientset()
		fakeClient.CoreV1().(*fakecorev1.FakeCoreV1).PrependReactor("create", "secrets", func(action faketesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, assert.AnError
		})
		clientMock := &mockClusterClient{fakeClient}
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return(string(validCertBytes), nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return(string(validCertKeyBytes), nil)
		sut := &manager{
			globalRegistry: registryMock,
			client:         clientMock,
		}

		// when
		_, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not create certificate secret k8s-ces-control-server-certificate")
	})

	t.Run("should fail because an error occur update the secret", func(t *testing.T) {
		oldCurrentNamespace := config.CurrentNamespace
		config.CurrentNamespace = "test-namespace"
		defer func() { config.CurrentNamespace = oldCurrentNamespace }()
		// given
		fakeClient := fake.NewSimpleClientset()
		fakeClient.CoreV1().(*fakecorev1.FakeCoreV1).PrependReactor("update", "secrets", func(action faketesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, assert.AnError
		})
		clientMock := &mockClusterClient{fakeClient}
		registryMock := newMockConfigurationContext(t)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return(string(validCertBytes), nil)
		registryMock.EXPECT().Get("certificate/k8s-ces-control/server.key").Return(string(validCertKeyBytes), nil)
		sut := &manager{
			globalRegistry: registryMock,
			client:         clientMock,
		}

		// when
		_, err := sut.GetCertificateCredentials(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not update certificate secret k8s-ces-control-server-certificate")
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

func (m *mockClusterClient) Dogus(_ string) ecoSystem.DoguInterface {
	panic("implement me")
}

func Test_setCertificateToRegistry(t *testing.T) {
	t.Run("should return nil if no error occurs", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Set("certificate/k8s-ces-control/server.crt", "server.crt").Return(nil)
		globalConfigMock.EXPECT().Set("certificate/cesappd/server.crt", "server.crt").Return(nil)
		globalConfigMock.EXPECT().Set("certificate/k8s-ces-control/server.key", "server.key").Return(nil)

		// when
		err := setCertificateToRegistry(globalConfigMock, "server.crt", "server.key")

		// then
		require.NoError(t, err)
	})
}
