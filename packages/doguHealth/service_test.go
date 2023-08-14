package doguHealth

import (
	"context"
	"github.com/cloudogu/k8s-ces-control/generated/health"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	doguv1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
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
		require.Error(t, err)
		assert.Equal(t, expectedResponse, actual)
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

func Test_server_GetByNames(t *testing.T) {
	t.Run("should fail to get deployment", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthListRequest{Dogus: []string{"will-fail", "will-succeed"}}
		deploymentMock := newMockDeploymentClient(t)
		deploymentMock.EXPECT().Get(context.TODO(), "will-fail", metav1.GetOptions{}).Return(nil, assert.AnError)
		healthyDeployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		deploymentMock.EXPECT().Get(context.TODO(), "will-succeed", metav1.GetOptions{}).Return(healthyDeployment, nil)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthMapResponse{
			AllHealthy: false,
			Results: map[string]*health.DoguHealthResponse{
				"will-fail": {
					FullName:    "will-fail",
					ShortName:   "will-fail",
					DisplayName: "will-fail",
					Healthy:     false,
					Results:     []*health.DoguHealthCheck{},
				},
				"will-succeed": {
					FullName:    "will-succeed",
					ShortName:   "will-succeed",
					DisplayName: "will-succeed",
					Healthy:     true,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a deployment contains at least one ready replica.",
					}},
				},
			},
		}

		// when
		actual, err := sut.GetByNames(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Equal(t, expectedResponse, actual)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail to get multiple deployments", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthListRequest{Dogus: []string{"will-fail", "will-fail-too"}}
		deploymentMock := newMockDeploymentClient(t)
		deploymentMock.EXPECT().Get(context.TODO(), "will-fail", metav1.GetOptions{}).Return(nil, assert.AnError)
		deploymentMock.EXPECT().Get(context.TODO(), "will-fail-too", metav1.GetOptions{}).Return(nil, assert.AnError)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthMapResponse{
			AllHealthy: false,
			Results: map[string]*health.DoguHealthResponse{
				"will-fail": {
					FullName:    "will-fail",
					ShortName:   "will-fail",
					DisplayName: "will-fail",
					Healthy:     false,
					Results:     []*health.DoguHealthCheck{},
				},
				"will-fail-too": {
					FullName:    "will-fail-too",
					ShortName:   "will-fail-too",
					DisplayName: "will-fail-too",
					Healthy:     false,
					Results:     []*health.DoguHealthCheck{},
				},
			},
		}

		// when
		actual, err := sut.GetByNames(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Equal(t, expectedResponse, actual)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should not all be healthy if one is unhealthy", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthListRequest{Dogus: []string{"healthy", "unhealthy"}}
		deploymentMock := newMockDeploymentClient(t)
		healthyDeployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		deploymentMock.EXPECT().Get(context.TODO(), "healthy", metav1.GetOptions{}).Return(healthyDeployment, nil)
		unhealthyDeployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 0}}
		deploymentMock.EXPECT().Get(context.TODO(), "unhealthy", metav1.GetOptions{}).Return(unhealthyDeployment, nil)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthMapResponse{
			AllHealthy: false,
			Results: map[string]*health.DoguHealthResponse{
				"healthy": {
					FullName:    "healthy",
					ShortName:   "healthy",
					DisplayName: "healthy",
					Healthy:     true,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a deployment contains at least one ready replica.",
					}},
				},
				"unhealthy": {
					FullName:    "unhealthy",
					ShortName:   "unhealthy",
					DisplayName: "unhealthy",
					Healthy:     false,
					Results:     []*health.DoguHealthCheck{},
				},
			},
		}

		// when
		actual, err := sut.GetByNames(context.TODO(), request)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, actual)
	})
	t.Run("should all be healthy", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthListRequest{Dogus: []string{"healthy1", "healthy2"}}
		deploymentMock := newMockDeploymentClient(t)
		healthy1Deployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		deploymentMock.EXPECT().Get(context.TODO(), "healthy1", metav1.GetOptions{}).Return(healthy1Deployment, nil)
		healthy2Deployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		deploymentMock.EXPECT().Get(context.TODO(), "healthy2", metav1.GetOptions{}).Return(healthy2Deployment, nil)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthMapResponse{
			AllHealthy: true,
			Results: map[string]*health.DoguHealthResponse{
				"healthy1": {
					FullName:    "healthy1",
					ShortName:   "healthy1",
					DisplayName: "healthy1",
					Healthy:     true,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a deployment contains at least one ready replica.",
					}},
				},
				"healthy2": {
					FullName:    "healthy2",
					ShortName:   "healthy2",
					DisplayName: "healthy2",
					Healthy:     true,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a deployment contains at least one ready replica.",
					}},
				},
			},
		}

		// when
		actual, err := sut.GetByNames(context.TODO(), request)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, actual)
	})
}

func Test_server_GetAll(t *testing.T) {
	t.Run("should fail to list dogus", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		doguClient := newMockDoguClient(t)
		doguClient.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(nil, assert.AnError)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClient)
		sut := &server{client: clientMock}

		// when
		actual, err := sut.GetAll(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should all be healthy for empty dogu list", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		doguList := &doguv1.DoguList{Items: []doguv1.Dogu{}}
		doguClient := newMockDoguClient(t)
		doguClient.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClient)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthMapResponse{
			AllHealthy: true,
			Results:    map[string]*health.DoguHealthResponse{},
		}

		// when
		actual, err := sut.GetAll(context.TODO(), nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, actual)
	})
	t.Run("should not all be healthy", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		deploymentMock := newMockDeploymentClient(t)
		healthyDeployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		unhealthyDeployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 0}}
		deploymentMock.EXPECT().Get(context.TODO(), "healthy", metav1.GetOptions{}).Return(healthyDeployment, nil)
		deploymentMock.EXPECT().Get(context.TODO(), "unhealthy", metav1.GetOptions{}).Return(unhealthyDeployment, nil)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)

		doguList := &doguv1.DoguList{Items: []doguv1.Dogu{
			{ObjectMeta: metav1.ObjectMeta{Name: "healthy"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "unhealthy"}},
		}}
		doguClient := newMockDoguClient(t)
		doguClient.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClient)
		clientMock.EXPECT().AppsV1().Return(appsMock)

		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthMapResponse{
			AllHealthy: false,
			Results: map[string]*health.DoguHealthResponse{
				"healthy": {
					FullName:    "healthy",
					ShortName:   "healthy",
					DisplayName: "healthy",
					Healthy:     true,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a deployment contains at least one ready replica.",
					}},
				},
				"unhealthy": {
					FullName:    "unhealthy",
					ShortName:   "unhealthy",
					DisplayName: "unhealthy",
					Healthy:     false,
					Results:     []*health.DoguHealthCheck{},
				},
			},
		}

		// when
		actual, err := sut.GetAll(context.TODO(), nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, actual)
	})
	t.Run("should all be healthy", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		deploymentMock := newMockDeploymentClient(t)
		healthy1Deployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		healthy2Deployment := &v1.Deployment{Status: v1.DeploymentStatus{ReadyReplicas: 1}}
		deploymentMock.EXPECT().Get(context.TODO(), "healthy1", metav1.GetOptions{}).Return(healthy1Deployment, nil)
		deploymentMock.EXPECT().Get(context.TODO(), "healthy2", metav1.GetOptions{}).Return(healthy2Deployment, nil)
		appsMock := newMockAppsV1Client(t)
		appsMock.EXPECT().Deployments("ecosystem").Return(deploymentMock)

		doguList := &doguv1.DoguList{Items: []doguv1.Dogu{
			{ObjectMeta: metav1.ObjectMeta{Name: "healthy1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "healthy2"}},
		}}
		doguClient := newMockDoguClient(t)
		doguClient.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClient)
		clientMock.EXPECT().AppsV1().Return(appsMock)

		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthMapResponse{
			AllHealthy: true,
			Results: map[string]*health.DoguHealthResponse{
				"healthy1": {
					FullName:    "healthy1",
					ShortName:   "healthy1",
					DisplayName: "healthy1",
					Healthy:     true,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a deployment contains at least one ready replica.",
					}},
				},
				"healthy2": {
					FullName:    "healthy2",
					ShortName:   "healthy2",
					DisplayName: "healthy2",
					Healthy:     true,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a deployment contains at least one ready replica.",
					}},
				},
			},
		}

		// when
		actual, err := sut.GetAll(context.TODO(), nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, actual)
	})
}
