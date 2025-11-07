package doguAdministration

import (
	"context"

	common "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	v3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-registry-lib/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BlueprintLister interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v3.BlueprintList, error)
}

type doguDescriptorGetter interface {
	// GetCurrentOfAll retrieves the specs of all dogus' currently installed versions.
	GetCurrentOfAll(ctx context.Context) ([]*core.Dogu, error)
}
type doguInterActor interface {
	// StartDogu starts the specified dogu
	StartDogu(ctx context.Context, doguName string) error
	// StopDogu stops the specified dogu
	StopDogu(ctx context.Context, doguName string) error
	// RestartDogu restarts the specified dogu
	RestartDogu(ctx context.Context, doguName string) error
}

//nolint:unused
//goland:noinspection GoUnusedType
type doguConfigRepository interface {
	Get(context.Context, common.SimpleName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
}
