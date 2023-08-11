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
		require.Contains(t, err.Error(), "failed to create CES registry: failed to create etcd client: parse")
	})

	tt.Run("Should success to register Services", func(t *testing.T) {
		// given
		mockGrpcServerRegistrar := &mockServiceRegistrar{registeredServices: []string{}}
		config.CurrentNamespace = "ecosystem"

		// when
		err := registerServices(nil, mockGrpcServerRegistrar)

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

func (sr *mockServiceRegistrar) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	sr.registeredServices = append(sr.registeredServices, desc.ServiceName)
}
