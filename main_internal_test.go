package main

import (
	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"os"
	"testing"
)

func Test_startCesControl(tt *testing.T) {
	tt.Run("Error on missing namespace environment variable", func(t *testing.T) {
		// given

		// when
		err := startCesControl()

		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "found invalid value for namespace []: namespace cannot be empty: set valid value with environment variable [NAMESPACE]")
	})

	tt.Run("Should succeed on help command", func(t *testing.T) {
		// given
		t.Setenv("NAMESPACE", "mynamespace")
		os.Args = []string{"k8s-ces-control", "help"}

		// when
		err := startCesControl()

		// then
		require.NoError(t, err)
	})
}

func Test_registerServices(tt *testing.T) {
	tt.Run("Should fail on error in registry", func(t *testing.T) {
		// given
		mockGrpcServerRegistrar := &mockServiceRegistrar{registeredServices: []string{}}
		config.CurrentNamespace = "%31$:://:../dir%25"

		// when
		err := registerServices(nil, mockGrpcServerRegistrar)

		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create CES registry: failed to create etcd client: parse \"http://etcd.%31$:://:../dir%25.svc.cluster.local:4001\"")
	})

	tt.Run("Should success to register Services", func(t *testing.T) {
		// given
		mockGrpcServerRegistrar := &mockServiceRegistrar{registeredServices: []string{}}
		config.CurrentNamespace = "ecosystem"
		clientSetMock := newMockClusterClient(t)
		coreV1Mock := newMockCoreV1Interface(t)
		appsV1Mock := newMockAppv1Interface(t)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		configMapInterfaceMock := newMockConfigMapInterface(t)
		deploymentInterfaceMock := newMockDeploymentInterface(t)
		coreV1Mock.EXPECT().ConfigMaps(config.CurrentNamespace).Return(configMapInterfaceMock)
		appsV1Mock.EXPECT().Deployments(config.CurrentNamespace).Return(deploymentInterfaceMock)

		// when
		err := registerServices(clientSetMock, mockGrpcServerRegistrar)

		// then
		require.NoError(t, err)
		assert.Equal(t, 5, len(mockGrpcServerRegistrar.registeredServices))
		assert.Contains(t, mockGrpcServerRegistrar.registeredServices, "logging.DoguLogMessages")
		assert.Contains(t, mockGrpcServerRegistrar.registeredServices, "doguAdministration.DoguAdministration")
		assert.Contains(t, mockGrpcServerRegistrar.registeredServices, "health.DoguHealth")
		assert.Contains(t, mockGrpcServerRegistrar.registeredServices, "maintenance.DebugMode")
		assert.Contains(t, mockGrpcServerRegistrar.registeredServices, "grpc.health.v1.Health")
	})
}

type mockServiceRegistrar struct {
	registeredServices []string
}

func (sr *mockServiceRegistrar) RegisterService(desc *grpc.ServiceDesc, _ interface{}) {
	sr.registeredServices = append(sr.registeredServices, desc.ServiceName)
}
