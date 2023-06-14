package internals

import (
	context "context"

	"google.golang.org/grpc/credentials"
)

// SSLManager is used to get the certificate.
type SSLManager interface {
	// GetCertificateCredentials reads the current ssl certificate and returns the transport credentials.
	GetCertificateCredentials(context.Context) (credentials.TransportCredentials, error)
}
