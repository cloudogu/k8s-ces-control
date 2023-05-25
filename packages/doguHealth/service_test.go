package doguHealth

import (
	"context"
	"github.com/cloudogu/k8s-ces-control/generated/health"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewDoguHealthService(t *testing.T) {
	t.Run("server should not be empty", func(t *testing.T) {
		// given
		clientMock := newMockClusterClient(t)

		// when
		actual := NewDoguHealthService(clientMock)

		// then
		assert.NotEmpty(t, actual)
		assert.Equal(t, clientMock, actual.client)
	})
}

func Test_server_GetByName(t *testing.T) {
	t.Run("should fail for empty dogu name", func(t *testing.T) {
		// given
		request := &health.DoguHealthRequest{DoguName: ""}
		clientMock := newMockClusterClient(t)
		sut := &server{client: clientMock}

		// when
		actual, err := sut.GetByName(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = dogu name is empty")
	})
	t.Run("should fail to get deployment", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthRequest{DoguName: "my-dogu"}
		deploymentMock := newMockDeploymentClient(t)
		deploymentMock.EXPECT().Get(context.TODO(), "my-dogu", metav1.GetOptions{}).Return(nil, assert.AnError)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		sut := &server{client: clientMock}

		// when
		actual, err := sut.GetByName(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return unhealthy for unhealthy dogu", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthRequest{DoguName: "my-dogu"}
		deploymentMock := newMockDeploymentClient(t)
		unhealthyDeployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 0}}
		deploymentMock.EXPECT().Get(context.TODO(), "my-dogu", metav1.GetOptions{}).Return(unhealthyDeployment, nil)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthResponse{
			FullName:    "my-dogu",
			ShortName:   "my-dogu",
			DisplayName: "my-dogu",
			Healthy:     false,
			Results:     []*health.DoguHealthCheck{},
		}

		// when
		actual, err := sut.GetByName(context.TODO(), request)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, actual)
	})
	t.Run("should return healthy for healthy dogu", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthRequest{DoguName: "my-dogu"}
		deploymentMock := newMockDeploymentClient(t)
		healthyDeployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		deploymentMock.EXPECT().Get(context.TODO(), "my-dogu", metav1.GetOptions{}).Return(healthyDeployment, nil)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthResponse{
			FullName:    "my-dogu",
			ShortName:   "my-dogu",
			DisplayName: "my-dogu",
			Healthy:     true,
			Results: []*health.DoguHealthCheck{{
				Type:    "container",
				Success: true,
				Message: "Check whether a deployment contains at least one ready replica.",
			}},
		}

		// when
		actual, err := sut.GetByName(context.TODO(), request)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, actual)
	})
}
