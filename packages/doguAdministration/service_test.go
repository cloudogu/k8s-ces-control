package doguAdministration

import (
	"context"
	"errors"
	"github.com/cloudogu/ces-control-api/generated/doguAdministration"
	"github.com/cloudogu/ces-control-api/generated/types"
	"github.com/cloudogu/cesapp-lib/core"
	blueprintcrv1 "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-ces-control/packages/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

var testCtx = context.TODO()

func TestNewDoguAdministrationServer(t *testing.T) {
	t.Run("server should not be empty", func(t *testing.T) {
		// given
		descriptorGetter := newMockDoguDescriptorGetter(t)
		bluePrintListerMock := NewMockBlueprintLister(t)
		doguInterActorMock := newMockDoguInterActor(t)
		loggingMock := newMockLogService(t)

		// when
		actual := NewDoguAdministrationServer(
			bluePrintListerMock,
			descriptorGetter,
			doguInterActorMock,
			loggingMock,
		)

		// then
		assert.NotEmpty(t, actual)
		assert.Equal(t, bluePrintListerMock, actual.blueprintLister)
		assert.Equal(t, descriptorGetter, actual.doguDescriptorGetter)
		assert.Equal(t, doguInterActorMock, actual.doguInterActor)
		assert.Equal(t, loggingMock, actual.loggingService)
	})
}

func Test_server_GetDoguList(t *testing.T) {
	t.Run("should return empty response for empty dogu list", func(t *testing.T) {
		// given
		descriptorGetter := newMockDoguDescriptorGetter(t)
		bluePrintListerMock := NewMockBlueprintLister(t)
		doguInterActorMock := newMockDoguInterActor(t)
		loggingMock := newMockLogService(t)

		descriptorGetter.EXPECT().GetCurrentOfAll(testCtx).Return(make([]*core.Dogu, 0), nil)

		sut := &server{
			blueprintLister:      bluePrintListerMock,
			doguDescriptorGetter: descriptorGetter,
			doguInterActor:       doguInterActorMock,
			loggingService:       loggingMock,
		}

		// when
		actual, err := sut.GetDoguList(testCtx, nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, &doguAdministration.DoguListResponse{}, actual)
	})
	t.Run("should fail to get dogu.jsons from registry", func(t *testing.T) {
		// given
		descriptorGetter := newMockDoguDescriptorGetter(t)
		bluePrintListerMock := NewMockBlueprintLister(t)
		doguInterActorMock := newMockDoguInterActor(t)
		loggingMock := newMockLogService(t)

		descriptorGetter.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)

		sut := &server{
			blueprintLister:      bluePrintListerMock,
			doguDescriptorGetter: descriptorGetter,
			doguInterActor:       doguInterActorMock,
			loggingService:       loggingMock,
		}

		// when
		actual, err := sut.GetDoguList(testCtx, nil)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get dogu registry")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		descriptorGetter := newMockDoguDescriptorGetter(t)
		bluePrintListerMock := NewMockBlueprintLister(t)
		doguInterActorMock := newMockDoguInterActor(t)
		loggingMock := newMockLogService(t)

		descriptorGetter.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{
			{
				Name:        "ns1/will-succeed",
				DisplayName: "WillSucceed",
				Version:     "1.2.3-2",
				Description: "asdf",
				Tags:        []string{"example"},
			},
			{
				Name:        "ns2/will-succeed-too",
				DisplayName: "WillSucceedToo",
				Version:     "3.2.1-1",
				Description: "qwert",
				Tags:        []string{"example", "banana"},
			},
		}, nil)

		loggingMock.EXPECT().GetLogLevel(mock.Anything, mock.Anything).Return(logging.LevelDebug, nil)

		sut := &server{
			blueprintLister:      bluePrintListerMock,
			doguDescriptorGetter: descriptorGetter,
			doguInterActor:       doguInterActorMock,
			loggingService:       loggingMock,
		}

		// when
		actual, err := sut.GetDoguList(testCtx, nil)

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
					LogLevel:    "LOG_LEVEL_DEBUG",
				},
				{
					Name:        "will-succeed-too",
					DisplayName: "WillSucceedToo",
					Version:     "3.2.1-1",
					Description: "qwert",
					Tags:        []string{"example", "banana"},
					LogLevel:    "LOG_LEVEL_DEBUG",
				},
			},
		}, actual)
	})
}

func Test_server_StartDogu(t *testing.T) {
	t.Run("should fail if dogu name is empty", func(t *testing.T) {
		// given
		descriptorGetter := newMockDoguDescriptorGetter(t)
		bluePrintListerMock := NewMockBlueprintLister(t)

		sut := &server{
			doguDescriptorGetter: descriptorGetter,
			blueprintLister:      bluePrintListerMock,
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
		bluePrintListerMock := NewMockBlueprintLister(t)
		descriptorGetter := newMockDoguDescriptorGetter(t)
		sut := &server{
			doguDescriptorGetter: descriptorGetter,
			blueprintLister:      bluePrintListerMock,
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
		bluePrintListerMock := NewMockBlueprintLister(t)
		descriptorGetter := newMockDoguDescriptorGetter(t)
		sut := &server{
			doguDescriptorGetter: descriptorGetter,
			blueprintLister:      bluePrintListerMock,
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

func Test_server_GetBlueprintId(t *testing.T) {
	ctx := context.TODO()

	t.Run("client List returns error", func(t *testing.T) {
		bluePrintListerMock := NewMockBlueprintLister(t)
		bluePrintListerMock.EXPECT().List(ctx, metav1.ListOptions{}).Return(nil, errors.New("testError"))

		sut := &server{
			doguDescriptorGetter: newMockDoguDescriptorGetter(t),
			blueprintLister:      bluePrintListerMock,
			doguInterActor:       newMockDoguInterActor(t),
		}

		request := &doguAdministration.DoguBlueprinitIdRequest{}

		// when
		actual, err := sut.GetBlueprintId(testCtx, request)
		require.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))

		require.Nil(t, actual)
	})

	t.Run("client List returns empty list", func(t *testing.T) {
		bluePrintListerMock := NewMockBlueprintLister(t)
		bluePrintListerMock.EXPECT().List(ctx, metav1.ListOptions{}).
			Return(&blueprintcrv1.BlueprintList{Items: make([]blueprintcrv1.Blueprint, 0)}, nil)

		sut := &server{
			doguDescriptorGetter: newMockDoguDescriptorGetter(t),
			blueprintLister:      bluePrintListerMock,
			doguInterActor:       newMockDoguInterActor(t),
		}

		request := &doguAdministration.DoguBlueprinitIdRequest{}

		// when
		actual, err := sut.GetBlueprintId(testCtx, request)
		require.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))

		require.Nil(t, actual)
	})

	t.Run("client List returns list with one element", func(t *testing.T) {
		bluePrintListerMock := NewMockBlueprintLister(t)
		bluePrintListerMock.EXPECT().List(ctx, metav1.ListOptions{}).
			Return(&blueprintcrv1.BlueprintList{Items: []blueprintcrv1.Blueprint{
				{ObjectMeta: metav1.ObjectMeta{
					Name:              "SIV1",
					CreationTimestamp: metav1.Now(),
				}},
			}}, nil)

		sut := &server{
			doguDescriptorGetter: newMockDoguDescriptorGetter(t),
			blueprintLister:      bluePrintListerMock,
			doguInterActor:       newMockDoguInterActor(t),
		}

		request := &doguAdministration.DoguBlueprinitIdRequest{}

		// when
		actual, err := sut.GetBlueprintId(testCtx, request)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "SIV1", actual.GetBlueprintId())
	})

	t.Run("client List returns list with two elements", func(t *testing.T) {
		bluePrintListerMock := NewMockBlueprintLister(t)
		bluePrintListerMock.EXPECT().List(ctx, metav1.ListOptions{}).
			Return(&blueprintcrv1.BlueprintList{Items: []blueprintcrv1.Blueprint{
				{ObjectMeta: metav1.ObjectMeta{
					Name:              "SIV1",
					CreationTimestamp: metav1.Now(),
				}},
				{ObjectMeta: metav1.ObjectMeta{
					Name:              "SIV2",
					CreationTimestamp: metav1.Now(),
				}},
			}}, nil)

		sut := &server{
			doguDescriptorGetter: newMockDoguDescriptorGetter(t),
			blueprintLister:      bluePrintListerMock,
			doguInterActor:       newMockDoguInterActor(t),
		}

		request := &doguAdministration.DoguBlueprinitIdRequest{}

		// when
		actual, err := sut.GetBlueprintId(testCtx, request)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "SIV2", actual.GetBlueprintId())
	})

	t.Run("client List returns list with two elements (replace order)", func(t *testing.T) {
		bluePrintListerMock := NewMockBlueprintLister(t)
		bluePrintListerMock.EXPECT().List(ctx, metav1.ListOptions{}).
			Return(&blueprintcrv1.BlueprintList{Items: []blueprintcrv1.Blueprint{
				{ObjectMeta: metav1.ObjectMeta{
					Name:              "SIV1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1 * time.Hour)},
				}},
				{ObjectMeta: metav1.ObjectMeta{
					Name:              "SIV2",
					CreationTimestamp: metav1.Now(),
				}},
			}}, nil)

		sut := &server{
			doguDescriptorGetter: newMockDoguDescriptorGetter(t),
			blueprintLister:      bluePrintListerMock,
			doguInterActor:       newMockDoguInterActor(t),
		}

		request := &doguAdministration.DoguBlueprinitIdRequest{}

		// when
		actual, err := sut.GetBlueprintId(testCtx, request)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "SIV1", actual.GetBlueprintId())
	})

}
