package doguinteraction

import (
	"context"
	common "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-dogu-operator/v2/api/ecoSystem"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type DoguInterface interface {
	ecoSystem.DoguInterface
}

type DoguRestartInterface interface {
	ecoSystem.DoguRestartInterface
}

type doguDescriptorGetter interface {
	// GetCurrentOfAll retrieves the specs of all dogus' currently installed versions.
	GetCurrentOfAll(ctx context.Context) ([]*core.Dogu, error)
}

type doguConfigRepository interface {
	Get(context.Context, common.SimpleName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
}
