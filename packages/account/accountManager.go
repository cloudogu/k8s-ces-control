package account

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-control/packages/ssl"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	hostConfigServiceName    = "k8s-ces-control"
	passwordLength           = 16
	randUsernameSuffixLength = 8
)

// ConfigurationRegistry provides functions to access the configuration registry.
type ConfigurationRegistry interface {
	// Set sets a configuration value in current context
	Set(key, value string) error
	// Get returns a configuration value from the current context
	Get(key string) (string, error)
}

// KeyProvider provides functions to access public and private keys of the system.
type KeyProvider interface {
	// FromPrivateKeyPath reads the keypair from the private key file path
	FromPrivateKeyPath(path string) (*keys.KeyPair, error)
}

// ServiceAccountManager provides methods to create or Delete service accounts.
type ServiceAccountManager struct {
	serviceName       string
	hostConfiguration ConfigurationRegistry
	keyProvider       KeyProvider
}

// ServiceAccountData holds the raw data for a service account.
type ServiceAccountData struct {
	Username string
	Password string
}

// NewServiceAccountManager creates a new instance of the ServiceAccountManager for the specified service.
func NewServiceAccountManager(serviceName string, registry registry.Registry) (ServiceAccountManager, error) {
	keyProvider, err := getKeyProvider(registry)
	if err != nil {
		return ServiceAccountManager{}, fmt.Errorf("failed to create ServiceAccountManager: %w", err)
	}

	return ServiceAccountManager{
		serviceName:       serviceName,
		keyProvider:       keyProvider,
		hostConfiguration: registry.HostConfig(hostConfigServiceName),
	}, nil
}

func (manager *ServiceAccountManager) SetServiceName(serviceName string) *ServiceAccountManager {
	manager.serviceName = serviceName
	return manager
}

func (manager *ServiceAccountManager) SetKeyProvider(keyProvider KeyProvider) *ServiceAccountManager {
	manager.keyProvider = keyProvider
	return manager
}

func (manager *ServiceAccountManager) SetHostConfiguration(configRegistry ConfigurationRegistry) *ServiceAccountManager {
	manager.hostConfiguration = configRegistry
	return manager
}

func (manager *ServiceAccountManager) GetHostConfiguration() ConfigurationRegistry {
	return manager.hostConfiguration
}

func (manager *ServiceAccountManager) GetServiceAccountData(ctx context.Context) (ServiceAccountData, error) {
	username, err := manager.GetHostConfiguration().Get(manager.serviceName + "/username")
	if err != nil {
		return ServiceAccountData{}, err
	}
	encryptedPassword, err := manager.GetHostConfiguration().Get(manager.serviceName + "/password")
	if err != nil {
		return ServiceAccountData{}, err
	}
	password, err := decrypt(encryptedPassword, ssl.CertificateKeyFilePath, manager.keyProvider)
	if err != nil {
		return ServiceAccountData{}, err
	}
	log.FromContext(ctx).Info("getting credentials for serviceaccount %s ...", manager.serviceName)
	return ServiceAccountData{Username: username, Password: password}, nil
}

func getKeyProvider(registry registry.Registry) (KeyProvider, error) {
	providerStr, err := registry.GlobalConfig().Get("key_provider")
	if err != nil {
		return nil, fmt.Errorf("failed to get key provider from global config: %w", err)
	}
	keyProvider, err := keys.NewKeyProvider(providerStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create key provider: %w", err)
	}
	return keyProvider, nil
}

func (accountData ServiceAccountData) String() string {
	return fmt.Sprintf("username:%s\npassword:%s", accountData.Username, accountData.Password)
}

// Create creates a new service account.
func (manager ServiceAccountManager) Create(context context.Context) (ServiceAccountData, error) {
	log.FromContext(context).Info("create service account for service '%s'", manager.serviceName)
	hostConfig := manager.hostConfiguration
	consumerName := manager.serviceName
	userData := generateUsernamePassword(consumerName)
	passwordEncrypted, err := encrypt(userData.Password, ssl.CertificateKeyFilePath, manager.keyProvider)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to encrypt the users password: %w", err)
	}
	err = hostConfig.Set(consumerName+"/username", userData.Username)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to write username to registry: %w", err)
	}
	err = hostConfig.Set(consumerName+"/password", passwordEncrypted)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to write password to registry: %w", err)
	}

	return userData, nil
}

func generateUsernamePassword(service string) ServiceAccountData {
	username := fmt.Sprintf("%s_%s", service, randomString(randUsernameSuffixLength))
	password := randomString(passwordLength)

	return ServiceAccountData{
		Username: username,
		Password: password,
	}
}

func encrypt(rawPassword, certificatePath string, provider KeyProvider) (string, error) {
	keyPair, err := provider.FromPrivateKeyPath(certificatePath)
	if err != nil {
		return "", fmt.Errorf("failed to load key pair data: %w", err)
	}
	publicKey := keyPair.Public()
	return publicKey.Encrypt(rawPassword)
}

func decrypt(encryptedValue, certificatePath string, provider KeyProvider) (string, error) {
	keyPair, err := provider.FromPrivateKeyPath(certificatePath)
	if err != nil {
		return "", fmt.Errorf("failed to load key pair data: %w", err)
	}
	return keyPair.Private().Decrypt(encryptedValue)
}
