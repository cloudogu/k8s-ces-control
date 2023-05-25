package ssl

import (
	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v2"
	"testing"
)

func TestNewManager(t *testing.T) {
	t.Run("should not be empty", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)

		// when
		actual := NewManager(globalConfigMock)

		// then
		assert.NotEmpty(t, actual)
		assert.NotNil(t, actual.certGenerator)
		assert.Equal(t, globalConfigMock, actual.globalRegistry)
	})
}

func Test_manager_GetCertificateCredentials(t *testing.T) {
	t.Run("should fail to get certificate", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Get("certificate/k8s-ces-control/server.crt").Return("", assert.AnError)
		sut := &manager{globalRegistry: globalConfigMock}

		// when
		actual, err := sut.GetCertificateCredentials()

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
		actual, err := sut.GetCertificateCredentials()

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
		actual, err := sut.GetCertificateCredentials()

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
		actual, err := sut.GetCertificateCredentials()

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
		actual, err := sut.GetCertificateCredentials()

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set certificate key in registry")
	})
	// TODO write tests for certificate creation
}
