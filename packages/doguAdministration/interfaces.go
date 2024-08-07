package doguAdministration

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/blueprintcr/v1"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
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
