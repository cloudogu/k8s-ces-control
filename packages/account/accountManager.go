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
	serviceName       string
	hostConfiguration registryContext
	keyProvider       keyProvider
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
		serviceName:       serviceName,
		keyProvider:       keyProvider,
		hostConfiguration: registry.HostConfig(hostConfigServiceName),
	}, nil
}

func (manager *ServiceAccountManager) SetServiceName(serviceName string) *ServiceAccountManager {
	manager.serviceName = serviceName
	return manager
}

func (manager *ServiceAccountManager) SetKeyProvider(keyProvider keyProvider) *ServiceAccountManager {
	manager.keyProvider = keyProvider
	return manager
}

func (manager *ServiceAccountManager) SetHostConfiguration(configRegistry registryContext) *ServiceAccountManager {
	manager.hostConfiguration = configRegistry
	return manager
}

func (manager *ServiceAccountManager) GetHostConfiguration() registryContext {
	return manager.hostConfiguration
}

func (manager *ServiceAccountManager) GetServiceAccountData(ctx context.Context) (ServiceAccountData, error) {
	username, err := manager.GetHostConfiguration().Get(manager.serviceName + "/username")
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to get username for service account '%s': %w", manager.serviceName, err)
	}
	encryptedPassword, err := manager.GetHostConfiguration().Get(manager.serviceName + "/password")
	if err != nil {
		return ServiceAccountData{}, fmt.Errorf("failed to get password for service account '%s': %w", manager.serviceName, err)
	}
	password, err := decrypt(encryptedPassword, ssl.CertificateKeyFilePath, manager.keyProvider)
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

func (accountData *ServiceAccountData) String() string {
	return fmt.Sprintf("username:%s\npassword:%s", accountData.Username, accountData.Password)
}

// Create creates a new service account.
func (manager *ServiceAccountManager) Create(context context.Context) (ServiceAccountData, error) {
	log.FromContext(context).Info("create service account for service '%s'", manager.serviceName)
	hostConfig := manager.hostConfiguration
	consumerName := manager.serviceName
	userData := generateUsernamePassword(consumerName)
	passwordEncrypted, err := encrypt(userData.Password, ssl.CertificateKeyFilePath, manager.keyProvider)
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

func encrypt(rawPassword, certificatePath string, provider keyProvider) (string, error) {
	keyPair, err := provider.FromPrivateKeyPath(certificatePath)
	if err != nil {
		return "", fmt.Errorf("failed to load key pair data: %w", err)
	}
	publicKey := keyPair.Public()
	return publicKey.Encrypt(rawPassword)
}

func decrypt(encryptedValue, certificatePath string, provider keyProvider) (string, error) {
	keyPair, err := provider.FromPrivateKeyPath(certificatePath)
	if err != nil {
		return "", fmt.Errorf("failed to load key pair data: %w", err)
	}
	return keyPair.Private().Decrypt(encryptedValue)
}
