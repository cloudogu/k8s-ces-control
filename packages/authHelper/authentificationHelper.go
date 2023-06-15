package authHelper

import (
	"context"

	"github.com/cloudogu/k8s-ces-control/packages/account"
	"github.com/cloudogu/k8s-ces-control/packages/config"
)

// AuthenticationFunc should retrieve credentials for a given service name.
type AuthenticationFunc func(ctx context.Context, serviceName string) (account.ServiceAccountData, error)

// GetServiceAccountCredentials return credentials for a given service name from the ces registry.
func GetServiceAccountCredentials(ctx context.Context, serviceName string) (account.ServiceAccountData, error) {
	cesRegistry, err := config.GetCesRegistry()
	if err != nil {
		return account.ServiceAccountData{}, err
	}

	manager, err := account.NewServiceAccountManager(serviceName, cesRegistry)
	if err != nil {
		return account.ServiceAccountData{}, err
	}

	serviceAccount, err := manager.GetServiceAccountData(ctx)
	if err != nil {
		return account.ServiceAccountData{}, err
	}
	return serviceAccount, nil
}
