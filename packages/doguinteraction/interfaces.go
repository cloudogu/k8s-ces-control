package doguinteraction

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"github.com/cloudogu/k8s-registry-lib/config"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type clusterClientSet interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type appsV1Interface interface {
	appsv1.AppsV1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type coreV1Interface interface {
	corev1.CoreV1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type deploymentInterface interface {
	appsv1.DeploymentInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type podInterface interface {
	corev1.PodInterface
}

type doguRegistry interface {
	// GetCurrentOfAll retrieves the specs of all dogus' currently installed versions.
	GetCurrentOfAll(ctx context.Context) ([]*core.Dogu, error)
}

//nolint:unused
//goland:noinspection GoUnusedType
type configurationContext interface {
	registry.ConfigurationContext
}

type doguConfigRepository interface {
	Get(context.Context, config.SimpleDoguName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
}
