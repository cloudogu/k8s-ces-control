package account

import (
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
)

// registryContext provides functions to access the configuration registry.
type registryContext interface {
	registry.ConfigurationContext
}

// keyProvider provides functions to access public and private keys of the system.
type keyProvider interface {
	// FromPrivateKeyPath reads the keypair from the private key file path
	FromPrivateKeyPath(path string) (*keys.KeyPair, error)
}

type configRegistry interface {
	registry.Registry
}
