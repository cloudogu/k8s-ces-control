package doguHealth

import (
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"k8s.io/client-go/kubernetes"
)

type clusterClient interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type doguClient interface {
	ecoSystem.DoguInterface
}
