package ssl

import (
	"crypto/tls"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/ssl"
	"github.com/cloudogu/k8s-ces-control/packages/config"
)

const (
	certificateRegistryKey       = "certificate/k8s-ces-control/server.crt"
	certificateLegacyRegistryKey = "certificate/cesappd/server.crt"
	CertificateKeyRegistryKey    = "certificate/k8s-ces-control/server.key"
)

type manager struct {
	globalRegistry configurationContext
	certGenerator  sslGenerator
}

// NewManager returns a new manager instance.
func NewManager(globalRegistry configurationContext) *manager {
	return &manager{
		globalRegistry: globalRegistry,
		certGenerator:  ssl.NewSSLGenerator(),
	}
}

// GetCertificateCredentials returns the certificate from the ces registry.
// If no certificate is found this routine generate a new self-signed certificate and writes it to the ces registry.
func (r *manager) GetCertificateCredentials() (credentials.TransportCredentials, error) {
	hasCertificate, err := r.hasCertificate()
	if err != nil {
		return nil, fmt.Errorf("failed to check if certificate exists: %w", err)
	}

	if !hasCertificate {
		logrus.Println("Found no ssl certificate -> generating new one.")

		cert, key, err := r.certGenerator.GenerateSelfSignedCert(
			"k8s-ces-control",
			"k8s-ces-control",
			24000,
			ssl.Country,
			ssl.Province,
			ssl.Locality,
			[]string{fmt.Sprintf("k8s-ces-control.%s.svc.cluster.local", config.CurrentNamespace), "localhost"},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate self-signed certificate: %w", err)
		}

		err = r.globalRegistry.Set(certificateRegistryKey, cert)
		if err != nil {
			return nil, fmt.Errorf("failed to set certificate in registry: %w", err)
		}

		err = r.globalRegistry.Set(certificateLegacyRegistryKey, cert)
		if err != nil {
			return nil, fmt.Errorf("failed to set certificate in registry legacy location: %w", err)
		}

		err = r.globalRegistry.Set(CertificateKeyRegistryKey, key)
		if err != nil {
			return nil, fmt.Errorf("failed to set certificate key in registry: %w", err)
		}
	}

	cert, err := r.createCertFromRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create cert from registry: %w", err)
	}

	return credentials.NewServerTLSFromCert(cert), nil
}

func (r *manager) createCertFromRegistry() (*tls.Certificate, error) {
	certPEMBlock, err := r.globalRegistry.Get(certificateRegistryKey)
	if err != nil {
		return nil, err
	}

	keyPEMBlock, err := r.globalRegistry.Get(CertificateKeyRegistryKey)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (r *manager) hasCertificate() (bool, error) {
	serverCrt, err := r.globalRegistry.Get(certificateRegistryKey)
	if registry.IsKeyNotFoundError(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return serverCrt != "", nil
}
