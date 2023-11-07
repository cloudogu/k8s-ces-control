package doguinteraction

import (
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"k8s.io/client-go/kubernetes"
)

type clusterClientSet interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}
