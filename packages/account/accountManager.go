package account

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/k8s-ces-control/packages/ssl"
)

const (
	hostConfigServiceName    = "k8s-ces-control"
	passwordLength           = 16
	randUsernameSuffixLength = 8
)

// ServiceAccountManager provides methods to create or Delete service accounts.
type ServiceAccountManager struct {
	serviceName         string
	hostConfiguration   registryContext
	globalConfiguration registryContext
	keyProvider         keyProvider
}

// ServiceAccountData holds the raw data for a service account.
type ServiceAccountData struct {
	Username string
	Password string
}

// NewServiceAccountManager creates a new instance of the ServiceAccountManager for the specified service.
func NewServiceAccountManager(serviceName string, registry configRegistry) (*ServiceAccountManager, error) {
	keyProvider, err := getKeyProvider(registry.GlobalConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create ServiceAccountManager: %w", err)
	}

	return &ServiceAccountManager{
		serviceName:         serviceName,
		keyProvider:         keyProvider,
		hostConfiguration:   registry.HostConfig(hostConfigServiceName),
		globalConfiguration: registry.GlobalConfig(),
	}, nil
}

// GetServiceAccountData returns credentials for the actual configured service.
func (manager *ServiceAccountManager) GetServiceAccountData(ctx context.Context) (ServiceAccountData, error) {
	username, err := manager.hostConfiguration.Get(manager.serviceName + "/username")
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to get username for service account '%s': %w", manager.serviceName, err)
	}

	encryptedPassword, err := manager.hostConfiguration.Get(manager.serviceName + "/password")
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to get password for service account '%s': %w", manager.serviceName, err)
	}

	privateKey, err := manager.globalConfiguration.Get(ssl.CertificateKeyRegistryKey)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to get private key from global config: %w", err)
	}

	password, err := decrypt(encryptedPassword, privateKey, manager.keyProvider)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to decrypt password for service account '%s': %w", manager.serviceName, err)
	}

	log.FromContext(ctx).Info("getting credentials for serviceaccount %s ...", manager.serviceName)
	return ServiceAccountData{Username: username, Password: password}, nil
}

func getKeyProvider(globalConfig registryContext) (keyProvider, error) {
	providerStr, err := globalConfig.Get("key_provider")
	if err != nil {
		return nil, fmt.Errorf("failed to get key provider from global config: %w", err)
	}

	keyProvider, err := keys.NewKeyProvider(providerStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create key provider: %w", err)
	}

	return keyProvider, nil
}

// String returns the string presentation of the ServiceAccountData object.
func (accountData ServiceAccountData) String() string {
	return fmt.Sprintf("username:%s\npassword:%s", accountData.Username, accountData.Password)
}

// Create creates a new service account and returns it.
func (manager *ServiceAccountManager) Create(context context.Context) (ServiceAccountData, error) {
	log.FromContext(context).Info("create service account for service '%s'", manager.serviceName)
	hostConfig := manager.hostConfiguration
	consumerName := manager.serviceName
	userData := generateUsernamePassword(consumerName)

	privateKey, err := manager.globalConfiguration.Get(ssl.CertificateKeyRegistryKey)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to get private key from global config: %w", err)
	}

	passwordEncrypted, err := encrypt(userData.Password, privateKey, manager.keyProvider)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to encrypt the users password for service account '%s': %w", consumerName, err)
	}

	err = hostConfig.Set(consumerName+"/username", userData.Username)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to write username to registry for service account '%s': %w", consumerName, err)
	}

	err = hostConfig.Set(consumerName+"/password", passwordEncrypted)
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to write password to registry for service account '%s': %w", consumerName, err)
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

func encrypt(rawPassword, privateKey string, provider keyProvider) (string, error) {
	keyPair, err := provider.FromPrivateKey([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("failed to load key pair data: %w", err)
	}

	publicKey := keyPair.Public()
	return publicKey.Encrypt(rawPassword)
}

func decrypt(encryptedValue, privateKey string, provider keyProvider) (string, error) {
	keyPair, err := provider.FromPrivateKey([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("failed to load key pair data: %w", err)
	}

	return keyPair.Private().Decrypt(encryptedValue)
}
