package logging

import (
	"context"
	common "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type doguConfigRepository interface {
	Get(context.Context, common.SimpleName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
}
