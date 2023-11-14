package main

import (
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//nolint:unused
//goland:noinspection GoUnusedType
type clusterClientSet interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type configMapInterface interface {
	v1.ConfigMapInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type coreV1Interface interface {
	v1.CoreV1Interface
}
