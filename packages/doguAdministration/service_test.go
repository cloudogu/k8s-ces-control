package doguAdministration

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-ces-control/generated/doguAdministration"
	"github.com/cloudogu/k8s-ces-control/generated/types"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var testCtx = context.TODO()

func TestNewDoguAdministrationServer(t *testing.T) {
	t.Run("server should not be empty", func(t *testing.T) {
		// given
		clientMock := newMockClusterClient(t)
		doguRegMock := newMockDoguRegistry(t)
		regMock := newMockCesRegistry(t)
		regMock.EXPECT().DoguRegistry().Return(doguRegMock)

		// when
		actual := NewDoguAdministrationServer(clientMock, regMock)

		// then
		assert.NotEmpty(t, actual)
		assert.Equal(t, clientMock, actual.client)
		assert.Equal(t, doguRegMock, actual.doguRegistry)
	})
}

func Test_server_GetDoguList(t *testing.T) {
	t.Run("should fail to list dogus", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(nil, assert.AnError)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClientMock)
		doguRegMock := newMockDoguRegistry(t)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}

		// when
		actual, err := sut.GetDoguList(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return empty response for empty dogu list", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(&v1.DoguList{}, nil)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClientMock)
		doguRegMock := newMockDoguRegistry(t)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}

		// when
		actual, err := sut.GetDoguList(context.TODO(), nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, &doguAdministration.DoguListResponse{}, actual)
	})
	t.Run("should fail to get one dogu.json from registry", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		doguList := &v1.DoguList{Items: []v1.Dogu{
			{ObjectMeta: metav1.ObjectMeta{Name: "will-succeed"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "will-fail"}},
		}}
		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClientMock)
		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get("will-succeed").Return(&core.Dogu{}, nil)
		doguRegMock.EXPECT().Get("will-fail").Return(nil, assert.AnError)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}

		// when
		actual, err := sut.GetDoguList(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "1 error occurred")
	})
	t.Run("should fail to get two dogu.jsons from registry", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		doguList := &v1.DoguList{Items: []v1.Dogu{
			{ObjectMeta: metav1.ObjectMeta{Name: "will-fail"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "will-fail-too"}},
		}}
		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClientMock)
		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get("will-fail").Return(nil, assert.AnError)
		doguRegMock.EXPECT().Get("will-fail-too").Return(nil, assert.AnError)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}

		// when
		actual, err := sut.GetDoguList(context.TODO(), nil)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "2 errors occurred")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		previousNamespaceVar := config.CurrentNamespace
		defer func() { config.CurrentNamespace = previousNamespaceVar }()
		config.CurrentNamespace = "ecosystem"

		doguList := &v1.DoguList{Items: []v1.Dogu{
			{ObjectMeta: metav1.ObjectMeta{Name: "will-succeed"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "will-succeed-too"}},
		}}
		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().List(context.TODO(), metav1.ListOptions{}).Return(doguList, nil)
		clientMock := newMockClusterClient(t)
		clientMock.EXPECT().Dogus("ecosystem").Return(doguClientMock)
		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get("will-succeed").Return(&core.Dogu{
			Name:        "ns1/will-succeed",
			DisplayName: "WillSucceed",
			Version:     "1.2.3-2",
			Description: "asdf",
			Tags:        []string{"example"},
		}, nil)
		doguRegMock.EXPECT().Get("will-succeed-too").Return(&core.Dogu{
			Name:        "ns2/will-succeed-too",
			DisplayName: "WillSucceedToo",
			Version:     "3.2.1-1",
			Description: "qwert",
			Tags:        []string{"example", "banana"},
		}, nil)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}

		// when
		actual, err := sut.GetDoguList(context.TODO(), nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, &doguAdministration.DoguListResponse{
			Dogus: []*doguAdministration.Dogu{
				{
					Name:        "will-succeed",
					DisplayName: "WillSucceed",
					Version:     "1.2.3-2",
					Description: "asdf",
					Tags:        []string{"example"},
				},
				{
					Name:        "will-succeed-too",
					DisplayName: "WillSucceedToo",
					Version:     "3.2.1-1",
					Description: "qwert",
					Tags:        []string{"example", "banana"},
				},
			},
		}, actual)
	})
}

func Test_server_StartDogu(t *testing.T) {
	t.Run("should fail if dogu name is empty", func(t *testing.T) {
		// given
		clientMock := newMockClusterClient(t)
		doguRegMock := newMockDoguRegistry(t)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: ""}

		// when
		actual, err := sut.StartDogu(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = dogu name is empty")
	})

	t.Run("should fail to start dogu", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StartDogu(testCtx, "my-dogu").Return(assert.AnError)

		sut := &server{
			doguInterActor: doguInterActorMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: "my-dogu"}

		// when
		actual, err := sut.StartDogu(testCtx, request)

		// then
		require.Error(t, err)
		assert.Equal(t, &types.BasicResponse{}, actual)
		assert.ErrorContains(t, err, assert.AnError.Error())
	})

	t.Run("should succeed", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StartDogu(testCtx, "my-dogu").Return(nil)

		sut := &server{
			doguInterActor: doguInterActorMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: "my-dogu"}

		// when
		actual, err := sut.StartDogu(testCtx, request)

		// then
		require.NoError(t, err)
		assert.Equal(t, &types.BasicResponse{}, actual)
	})
}

func Test_server_StopDogu(t *testing.T) {
	t.Run("should fail if dogu name is empty", func(t *testing.T) {
		// given
		clientMock := newMockClusterClient(t)
		doguRegMock := newMockDoguRegistry(t)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: ""}

		// when
		actual, err := sut.StopDogu(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = dogu name is empty")
	})

	t.Run("should fail to stop dogu", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopDogu(testCtx, "my-dogu").Return(assert.AnError)

		sut := &server{
			doguInterActor: doguInterActorMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: "my-dogu"}

		// when
		actual, err := sut.StopDogu(testCtx, request)

		// then
		require.Error(t, err)
		assert.Equal(t, &types.BasicResponse{}, actual)
		assert.ErrorContains(t, err, assert.AnError.Error())
	})

	t.Run("should succeed", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().StopDogu(testCtx, "my-dogu").Return(nil)

		sut := &server{
			doguInterActor: doguInterActorMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: "my-dogu"}

		// when
		actual, err := sut.StopDogu(testCtx, request)

		// then
		require.NoError(t, err)
		assert.Equal(t, &types.BasicResponse{}, actual)
	})
}

func Test_server_RestartDogu(t *testing.T) {
	t.Run("should fail if dogu name is empty", func(t *testing.T) {
		// given
		clientMock := newMockClusterClient(t)
		doguRegMock := newMockDoguRegistry(t)
		sut := &server{
			doguRegistry: doguRegMock,
			client:       clientMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: ""}

		// when
		actual, err := sut.RestartDogu(context.TODO(), request)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = dogu name is empty")
	})

	t.Run("should fail to restart dogu", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().RestartDogu(testCtx, "my-dogu").Return(assert.AnError)

		sut := &server{
			doguInterActor: doguInterActorMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: "my-dogu"}

		// when
		actual, err := sut.RestartDogu(testCtx, request)

		// then
		require.Error(t, err)
		assert.Equal(t, &types.BasicResponse{}, actual)
		assert.ErrorContains(t, err, assert.AnError.Error())
	})

	t.Run("should succeed", func(t *testing.T) {
		// given
		doguInterActorMock := newMockDoguInterActor(t)
		doguInterActorMock.EXPECT().RestartDogu(testCtx, "my-dogu").Return(nil)

		sut := &server{
			doguInterActor: doguInterActorMock,
		}
		request := &doguAdministration.DoguAdministrationRequest{DoguName: "my-dogu"}

		// when
		actual, err := sut.RestartDogu(testCtx, request)

		// then
		require.NoError(t, err)
		assert.Equal(t, &types.BasicResponse{}, actual)
	})
}
