package main

import (
	"github.com/cloudogu/k8s-ces-control/packages/doguAdministration"
	ecoSystemV2 "github.com/cloudogu/k8s-dogu-operator/v2/api/ecoSystem"
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//nolint:unused
//goland:noinspection GoUnusedType
type configMapInterface interface {
	coreV1.ConfigMapInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type coreV1Interface interface {
	coreV1.CoreV1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type appsV1Interface interface {
	appsV1.AppsV1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type clusterClient interface {
	ecoSystemV2.EcoSystemV2Interface
	doguAdministration.BlueprintLister
	kubernetes.Interface
}
