package doguHealth

import (
	"github.com/cloudogu/k8s-dogu-operator/v2/api/ecoSystem"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type clusterClient interface {
	ecoSystem.EcoSystemV2Interface
	kubernetes.Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type doguClient interface {
	ecoSystem.DoguInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type deploymentClient interface {
	v1.DeploymentInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type appsV1Client interface {
	v1.AppsV1Interface
}
