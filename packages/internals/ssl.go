package internals

import "google.golang.org/grpc/credentials"

// SSLManager is used to get the certificate.
type SSLManager interface {
	// GetCertificateCredentials reads the current ssl certificate and returns the transport credentials.
	GetCertificateCredentials() (credentials.TransportCredentials, error)
}
