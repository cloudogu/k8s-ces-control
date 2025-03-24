package support

import (
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sClient interface {
	client.Client
}

type discoveryInterface interface {
	discovery.DiscoveryInterface
}
