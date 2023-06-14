package ssl

import (
	context "context"
	"crypto/tls"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/ssl"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
)

const (
	certificateLegacyRegistryKey = "certificate/cesappd/server.crt"
	certificateRegistryKey       = "certificate/k8s-ces-control/server.crt"
	CertificateKeyRegistryKey    = "certificate/k8s-ces-control/server.key"
)

type clusterClient interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

type manager struct {
	globalRegistry configurationContext
	certGenerator  sslGenerator
	client         clusterClient
}

// NewManager returns a new manager instance.
func NewManager(client clusterClient, globalRegistry configurationContext) *manager {
	return &manager{
		globalRegistry: globalRegistry,
		certGenerator:  ssl.NewSSLGenerator(),
		client:         client,
	}
}

// GetCertificateCredentials returns the certificate from the ces registry.
// If no certificate is found this routine generate a new self-signed certificate and writes it to the ces registry.
func (r *manager) GetCertificateCredentials(ctx context.Context) (credentials.TransportCredentials, error) {
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

		err = setCertificateToRegistry(r.globalRegistry, cert, key)
		if err != nil {
			return nil, err
		}

		err = createCertificateSecret(ctx, config.CurrentNamespace, r.client, cert, key)
		if err != nil {
			return nil, err
		}
	}

	cert, err := r.createCertFromRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create cert from registry: %w", err)
	}

	return credentials.NewServerTLSFromCert(cert), nil
}

func createCertificateSecret(ctx context.Context, namespace string, client clusterClient, cert, key string) error {
	data := map[string]string{corev1.TLSCertKey: cert, corev1.TLSPrivateKeyKey: key}
	const secretName = "k8s-ces-control-server-certificate"
	var updateOpts metav1.UpdateOptions
	var getOpts metav1.GetOptions
	var createOpts metav1.CreateOptions

	creds := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: data,
		Type:       corev1.SecretTypeTLS,
	}

	_, err := client.CoreV1().Secrets(config.CurrentNamespace).Get(ctx, secretName, getOpts)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("error while checking whether certificate secret %s already exists: %w", secretName, err)
		}

		logrus.Info("did not found a certificate secret, creating one now")
		_, err = client.CoreV1().Secrets(config.CurrentNamespace).Create(ctx, &creds, createOpts)
		if err != nil {
			return fmt.Errorf("could not patch certificate secret %s: %w", secretName, err)
		}
	}

	_, err = client.CoreV1().Secrets(config.CurrentNamespace).Update(ctx, &creds, updateOpts)
	if err != nil {
		return fmt.Errorf("could not update certificate secret %s: %w", secretName, err)
	}

	logrus.Info("successfully update (even the shortly created) certificate as secret")
	return nil
}

func setCertificateToRegistry(globalReg configurationContext, cert string, key string) error {
	err := globalReg.Set(certificateRegistryKey, cert)
	if err != nil {
		return fmt.Errorf("failed to set certificate in registry: %w", err)
	}

	err = globalReg.Set(certificateLegacyRegistryKey, cert)
	if err != nil {
		return fmt.Errorf("failed to set certificate in registry legacy location: %w", err)
	}

	err = globalReg.Set(CertificateKeyRegistryKey, key)
	if err != nil {
		return fmt.Errorf("failed to set certificate key in registry: %w", err)
	}

	return nil
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
