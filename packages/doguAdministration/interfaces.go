package doguAdministration

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/blueprintcr/v1"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type BlueprintLister interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1.BlueprintList, error)
}

type clusterClient interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	BlueprintLister
	kubernetes.Interface
}

type doguRegistry interface {
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

type doguConfigRepository interface {
	Get(context.Context, config.SimpleDoguName) (config.DoguConfig, error)
	Delete(context.Context, config.SimpleDoguName) error
	Create(context.Context, config.DoguConfig) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
	SaveOrMerge(context.Context, config.DoguConfig) (config.DoguConfig, error)
	Watch(ctx context.Context, dName config.SimpleDoguName, filters ...config.WatchFilter) (<-chan repository.DoguConfigWatchResult, error)
}
