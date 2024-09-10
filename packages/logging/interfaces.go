package logging

import (
	"context"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type doguConfigRepository interface {
	Get(context.Context, config.SimpleDoguName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
}
