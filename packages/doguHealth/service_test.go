package doguHealth

import (
	"context"
	"github.com/cloudogu/ces-control-api/generated/health"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	doguv2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewDoguHealthService(t *testing.T) {
	t.Run("server should not be empty", func(t *testing.T) {
		// given
		clientMock := newMockDoguClient(t)

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
		clientMock := newMockDoguClient(t)
		sut := &server{client: clientMock}

		// when
		actual, err := sut.GetByName(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = dogu name is empty")
	})
	t.Run("should fail to get dogu", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthRequest{DoguName: "my-dogu"}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "my-dogu", metav1.GetOptions{}).Return(nil, assert.AnError)
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
		var replicas = new(int32)
		*replicas = 1
		unhealthyDogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.UnavailableHealthStatus}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "my-dogu", metav1.GetOptions{}).Return(unhealthyDogu, nil)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthResponse{
			FullName:    "my-dogu",
			ShortName:   "my-dogu",
			DisplayName: "my-dogu",
			Healthy:     false,
			Results: []*health.DoguHealthCheck{
				{
					Type:     "container",
					Success:  true,
					Message:  "Check whether a dogu is not stopped (ready or not).",
					Critical: false,
				},
			},
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
		var replicas = new(int32)
		*replicas = 1
		healthyDogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "my-dogu", metav1.GetOptions{}).Return(healthyDogu, nil)
		sut := &server{client: clientMock}
		expectedResponse := &health.DoguHealthResponse{
			FullName:    "my-dogu",
			ShortName:   "my-dogu",
			DisplayName: "my-dogu",
			Healthy:     true,
			Results: []*health.DoguHealthCheck{{
				Type:    "container",
				Success: true,
				Message: "Check whether a dogu is not stopped (ready or not).",
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
	t.Run("should fail to get dogu", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthListRequest{Dogus: []string{"will-fail", "will-succeed"}}
		var replicas = new(int32)
		*replicas = 1
		healthyDogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "will-fail", metav1.GetOptions{}).Return(nil, assert.AnError)
		clientMock.EXPECT().Get(context.TODO(), "will-succeed", metav1.GetOptions{}).Return(healthyDogu, nil)
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
						Type:     "container",
						Success:  true,
						Message:  "Check whether a dogu is not stopped (ready or not).",
						Critical: false,
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
	t.Run("should fail to get multiple dogus", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthListRequest{Dogus: []string{"will-fail", "will-fail-too"}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "will-fail", metav1.GetOptions{}).Return(nil, assert.AnError)
		clientMock.EXPECT().Get(context.TODO(), "will-fail-too", metav1.GetOptions{}).Return(nil, assert.AnError)
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
		var replicasHealthy = new(int32)
		*replicasHealthy = 1
		healthyDogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}
		unhealthyDogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.UnavailableHealthStatus}, Spec: doguv2.DoguSpec{Stopped: true}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "healthy", metav1.GetOptions{}).Return(healthyDogu, nil)
		clientMock.EXPECT().Get(context.TODO(), "unhealthy", metav1.GetOptions{}).Return(unhealthyDogu, nil)
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
						Type:     "container",
						Success:  true,
						Message:  "Check whether a dogu is not stopped (ready or not).",
						Critical: false,
					}},
				},
				"unhealthy": {
					FullName:    "unhealthy",
					ShortName:   "unhealthy",
					DisplayName: "unhealthy",
					Healthy:     false,
					Results: []*health.DoguHealthCheck{{
						Type:     "container",
						Success:  false,
						Message:  "Check whether a dogu is not stopped (ready or not).",
						Critical: false,
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
	t.Run("should all be healthy", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		request := &health.DoguHealthListRequest{Dogus: []string{"healthy1", "healthy2"}}

		healthy1Dogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}
		healthy2Dogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}

		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "healthy1", metav1.GetOptions{}).Return(healthy1Dogu, nil)
		clientMock.EXPECT().Get(context.TODO(), "healthy2", metav1.GetOptions{}).Return(healthy2Dogu, nil)
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
						Message: "Check whether a dogu is not stopped (ready or not).",
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
						Message: "Check whether a dogu is not stopped (ready or not).",
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

		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(nil, assert.AnError)
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

		doguList := &doguv2.DoguList{Items: []doguv2.Dogu{}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
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

		healthyDogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}
		unhealthyDogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.UnavailableHealthStatus}}

		doguList := &doguv2.DoguList{Items: []doguv2.Dogu{
			{ObjectMeta: metav1.ObjectMeta{Name: "healthy"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "unhealthy"}},
		}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
		clientMock.EXPECT().Get(context.TODO(), "healthy", metav1.GetOptions{}).Return(healthyDogu, nil)
		clientMock.EXPECT().Get(context.TODO(), "unhealthy", metav1.GetOptions{}).Return(unhealthyDogu, nil)

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
						Message: "Check whether a dogu is not stopped (ready or not).",
					}},
				},
				"unhealthy": {
					FullName:    "unhealthy",
					ShortName:   "unhealthy",
					DisplayName: "unhealthy",
					Healthy:     false,
					Results: []*health.DoguHealthCheck{{
						Type:    "container",
						Success: true,
						Message: "Check whether a dogu is not stopped (ready or not).",
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
	t.Run("should all be healthy", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		healthy1Dogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}
		healthy2Dogu := &doguv2.Dogu{Status: doguv2.DoguStatus{Health: doguv2.AvailableHealthStatus}}

		doguList := &doguv2.DoguList{Items: []doguv2.Dogu{
			{ObjectMeta: metav1.ObjectMeta{Name: "healthy1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "healthy2"}},
		}}
		clientMock := newMockDoguClient(t)
		clientMock.EXPECT().Get(context.TODO(), "healthy1", metav1.GetOptions{}).Return(healthy1Dogu, nil)
		clientMock.EXPECT().Get(context.TODO(), "healthy2", metav1.GetOptions{}).Return(healthy2Dogu, nil)
		clientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)

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
						Message: "Check whether a dogu is not stopped (ready or not).",
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
						Message: "Check whether a dogu is not stopped (ready or not).",
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
