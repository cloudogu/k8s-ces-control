package logging

import (
	"context"

	"github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type doguConfigRepository interface {
	Get(context.Context, dogu.SimpleName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
}
