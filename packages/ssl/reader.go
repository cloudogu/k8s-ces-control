package ssl

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/ssl"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
)

const (
	certificateRegistryKey       = "certificate/k8s-ces-control/server.crt"
	certificateLegacyRegistryKey = "certificate/cesappd/server.crt"
	certificateFilePath          = "/etc/k8s-ces-control/server.crt"
	certificateKeyRegistryKey    = "certificate/k8s-ces-control/server.key"
	CertificateKeyFilePath       = "/etc/k8s-ces-control/server.key"
)

type manager struct {
	globalRegistry registry.ConfigurationContext
}

func NewManager(globalRegistry registry.ConfigurationContext) *manager {
	return &manager{globalRegistry: globalRegistry}
}

func (r manager) GetCertificateCredentials() (credentials.TransportCredentials, error) {
	hasCertificate, err := r.hasCertificate()
	if err != nil {
		return nil, err
	}

	if !hasCertificate {
		logrus.Println("Found no ssl certificate -> generating new one.")

		cert, key, err := ssl.NewSSLGenerator().GenerateSelfSignedCert(
			"k8s-ces-control",
			"k8s-ces-control",
			24000,
			"DE",
			"Lower Saxony",
			"Brunswick",
			[]string{fmt.Sprintf("k8s-ces-control.%s.svc.cluster.local", config.CurrentNamespace), "localhost"},
		)
		if err != nil {
			return nil, err
		}

		err = r.globalRegistry.Set(certificateRegistryKey, cert)
		if err != nil {
			return nil, err
		}

		err = r.globalRegistry.Set(certificateLegacyRegistryKey, cert)
		if err != nil {
			return nil, err
		}

		err = r.globalRegistry.Set(certificateKeyRegistryKey, key)
		if err != nil {
			return nil, err
		}
	}

	logrus.Println("Found existing SSL certificate.")
	err = r.copyFromRegistryToFile(certificateRegistryKey, certificateFilePath)
	if err != nil {
		return nil, err
	}

	err = r.copyFromRegistryToFile(certificateKeyRegistryKey, CertificateKeyFilePath)
	if err != nil {
		return nil, err
	}

	return credentials.NewServerTLSFromFile(certificateFilePath, CertificateKeyFilePath)
}

func (r manager) hasCertificate() (bool, error) {
	serverCrt, err := r.globalRegistry.Get(certificateRegistryKey)
	if registry.IsKeyNotFoundError(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return serverCrt != "", nil
}

func (r manager) copyFromRegistryToFile(registryKey string, fileName string) error {
	value, err := r.globalRegistry.Get(registryKey)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(value)
	if err != nil {
		return fmt.Errorf("failed to write file [%s]: %w", f.Name(), err)
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
