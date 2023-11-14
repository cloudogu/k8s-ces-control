package doguAdministration

import (
	"context"
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

type doguInterActor interface {
	// StartDogu starts the specified dogu
	StartDogu(ctx context.Context, doguName string) error
	// StopDogu stops the specified dogu
	StopDogu(ctx context.Context, doguName string) error
	// RestartDogu restarts the specified dogu
	RestartDogu(ctx context.Context, doguName string) error
}
