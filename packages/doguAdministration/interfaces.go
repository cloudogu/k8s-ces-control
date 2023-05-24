package doguAdministration

import (
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"k8s.io/client-go/kubernetes"
)

type clusterClient interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

type doguRegistry interface {
	registry.DoguRegistry
}

type cesRegistry interface {
	registry.Registry
}
