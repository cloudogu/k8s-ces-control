package logging

import (
	"context"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

type doguConfigRepository interface {
	Get(context.Context, config.SimpleDoguName) (config.DoguConfig, error)
	Delete(context.Context, config.SimpleDoguName) error
	Create(context.Context, config.DoguConfig) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
	SaveOrMerge(context.Context, config.DoguConfig) (config.DoguConfig, error)
	Watch(ctx context.Context, dName config.SimpleDoguName, filters ...config.WatchFilter) (<-chan repository.DoguConfigWatchResult, error)
}
