package logging

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	pb "github.com/cloudogu/ces-control-api/generated/logging"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	"testing"
	"time"
)

func TestNewLoggingService(t *testing.T) {
	t.Run("should create query clock", func(t *testing.T) {
		// given
		llp := &LokiLogProvider{}

		// when
		sut := NewLoggingService(llp, newMockConfigProvider(t), newMockDoguRestarter(t), newMockDoguDescriptionGetter(t), newMockDeploymentGetter(t))

		// then
		require.NotNil(t, sut)
		assert.NotNil(t, sut.logProvider)
		assert.Equal(t, llp, sut.logProvider)
	})
}

func Test_writeLogLinesToStream(t *testing.T) {
	t.Run("should write logs to stream", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesServer(t)

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().getLogs("my-dogu", 222).Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(mock.Anything).Return(nil)

		// when
		err := writeLogLinesToStream(mockedLogProvider, "my-dogu", 222, mockedDoguLogServer)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail for empty dogu-name", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesServer(t)

		// when
		err := writeLogLinesToStream(mockedLogProvider, "", 222, mockedDoguLogServer)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = dogu name should not be empty")
	})

	t.Run("should fail for error in log provider", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesServer(t)

		mockedLogProvider.EXPECT().getLogs("my-dogu", 222).Return(nil, assert.AnError)

		// when
		err := writeLogLinesToStream(mockedLogProvider, "my-dogu", 222, mockedDoguLogServer)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = assert.AnError general error for testing")
	})

	t.Run("should fail for error in grpc-send", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesServer(t)

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().getLogs("my-dogu", 222).Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(mock.Anything).Return(assert.AnError)

		// when
		err := writeLogLinesToStream(mockedLogProvider, "my-dogu", 222, mockedDoguLogServer)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = Internal desc = assert.AnError general error for testing")
	})
}

func Test_compressMessages(t *testing.T) {
	t.Run("should compress log lines", func(t *testing.T) {
		// given
		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}

		// when
		actualMyDoguLogZip, err := compressMessages("my-dogu", logLines)

		// then
		require.NoError(t, err)
		assert.NotNil(t, actualMyDoguLogZip)

		// test text equality
		zipreader, err := zip.NewReader(bytes.NewReader(actualMyDoguLogZip), int64(len(actualMyDoguLogZip)))
		require.NoError(t, err)
		assert.Len(t, zipreader.File, 1)

		for _, zipfile := range zipreader.File {
			assert.Equal(t, "my-dogu.log", zipfile.Name)

			fc, err := zipfile.Open()
			require.NoError(t, err)
			defer func() { _ = fc.Close() }()
			actualFileContent, err := io.ReadAll(fc)
			require.NoError(t, err)

			assert.Len(t, actualFileContent, 333)
		}
	})

	t.Run("should not compress for empty log lines", func(t *testing.T) {
		// given

		// when
		actualMyDoguLogZip, err := compressMessages("my-dogu", []logLine{})

		// then
		require.NoError(t, err)
		assert.NotNil(t, actualMyDoguLogZip)
		assert.Empty(t, actualMyDoguLogZip)
	})
}

func Test_GetForDogu(t *testing.T) {
	t.Run("should get logs for Dogu", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesServer(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().getLogs("my-dogu", 333).Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(mock.Anything).Return(nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		// when
		request := &pb.DoguLogMessageRequest{
			DoguName:  "my-dogu",
			LineCount: 333,
		}
		err := sut.GetForDogu(request, mockedDoguLogServer)

		// then
		require.NoError(t, err)
	})
}

func Test_QueryForDogu(t *testing.T) {
	t.Run("should get logs for Dogu", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesQueryServer(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		start := time.Unix(1711616504, 0).UTC()
		end := time.Unix(1712131304, 0).UTC()
		filter := "foo=bar"

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().queryLogs("my-dogu", start, end, filter).Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[0].timestamp), Message: logLines[0].value}).Return(nil)
		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[1].timestamp), Message: logLines[1].value}).Return(nil)
		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[2].timestamp), Message: logLines[2].value}).Return(nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		// when
		request := &pb.DoguLogMessageQueryRequest{
			DoguName:  "my-dogu",
			StartDate: timestamppb.New(start),
			EndDate:   timestamppb.New(end),
			Filter:    &filter,
		}
		err := sut.QueryForDogu(request, mockedDoguLogServer)

		// then
		require.NoError(t, err)
	})

	t.Run("should get logs for Dogu without optional parameters", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesQueryServer(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().queryLogs("my-dogu", time.Time{}, time.Time{}, "").Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[0].timestamp), Message: logLines[0].value}).Return(nil)
		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[1].timestamp), Message: logLines[1].value}).Return(nil)
		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[2].timestamp), Message: logLines[2].value}).Return(nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		// when
		request := &pb.DoguLogMessageQueryRequest{
			DoguName:  "my-dogu",
			StartDate: nil,
			EndDate:   nil,
			Filter:    nil,
		}
		err := sut.QueryForDogu(request, mockedDoguLogServer)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to get logs for empty doguname", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesQueryServer(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		// when
		request := &pb.DoguLogMessageQueryRequest{
			DoguName: "",
		}
		err := sut.QueryForDogu(request, mockedDoguLogServer)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = dogu name should not be empty")
	})

	t.Run("should fail to get logs for error querying the log provider", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesQueryServer(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		start := time.Unix(1711616504, 0).UTC()
		end := time.Unix(1712131304, 0).UTC()
		filter := "foo=bar"

		mockedLogProvider.EXPECT().queryLogs("my-dogu", start, end, filter).Return(nil, assert.AnError)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		// when
		request := &pb.DoguLogMessageQueryRequest{
			DoguName:  "my-dogu",
			StartDate: timestamppb.New(start),
			EndDate:   timestamppb.New(end),
			Filter:    &filter,
		}
		err := sut.QueryForDogu(request, mockedDoguLogServer)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, assert.AnError.Error())
	})

	t.Run("should fail to get logs for error while sending", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesQueryServer(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		start := time.Unix(1711616504, 0).UTC()
		end := time.Unix(1712131304, 0).UTC()
		filter := "foo=bar"

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().queryLogs("my-dogu", start, end, filter).Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[0].timestamp), Message: logLines[0].value}).Return(nil)
		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[1].timestamp), Message: logLines[1].value}).Return(assert.AnError)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		// when
		request := &pb.DoguLogMessageQueryRequest{
			DoguName:  "my-dogu",
			StartDate: timestamppb.New(start),
			EndDate:   timestamppb.New(end),
			Filter:    &filter,
		}
		err := sut.QueryForDogu(request, mockedDoguLogServer)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, assert.AnError.Error())
	})
}

func TestLoggingService_SetLogLevel(t *testing.T) {
	tests := []struct {
		name          string
		req           *pb.LogLevelRequest
		xResponse     bool
		xResponseCode codes.Code
	}{
		{
			name: "Set LogLevel DEBUG",
			req: &pb.LogLevelRequest{
				DoguName: "test",
				LogLevel: pb.LogLevel_DEBUG,
			},
			xResponse:     true,
			xResponseCode: codes.OK,
		},
		{
			name: "Set LogLevel INFO",
			req: &pb.LogLevelRequest{
				DoguName: "test",
				LogLevel: pb.LogLevel_INFO,
			},
			xResponse:     true,
			xResponseCode: codes.OK,
		},
		{
			name: "Set LogLevel WARN",
			req: &pb.LogLevelRequest{
				DoguName: "test",
				LogLevel: pb.LogLevel_WARN,
			},
			xResponse:     true,
			xResponseCode: codes.OK,
		},
		{
			name: "Set LogLevel ERROR",
			req: &pb.LogLevelRequest{
				DoguName: "test",
				LogLevel: pb.LogLevel_ERROR,
			},
			xResponse:     true,
			xResponseCode: codes.OK,
		},
		{
			name: "Set wrong LogLevel",
			req: &pb.LogLevelRequest{
				DoguName: "test",
				LogLevel: 4,
			},
			xResponse:     false,
			xResponseCode: codes.InvalidArgument,
		},
		{
			name: "Empty dogu name",
			req: &pb.LogLevelRequest{
				DoguName: "",
				LogLevel: pb.LogLevel_DEBUG,
			},
			xResponse:     false,
			xResponseCode: codes.InvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockedLogProvider := newMockLogProvider(t)
			mockedConfigProvider := newMockConfigProvider(t)
			mockedDoguRestarter := newMockDoguRestarter(t)
			mockedConfigurationContext := NewMockConfigurationContext(t)
			mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
			mockedDeploymentGetter := newMockDeploymentGetter(t)

			if tc.xResponse {
				mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)

				mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(true, "TEST", nil)
				mockedConfigurationContext.EXPECT().Set(mock.Anything, mock.Anything).Return(nil)

				replica := int32(1)
				mockedDeploymentGetter.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(&appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: &replica,
					},
				}, nil)

				mockedDoguRestarter.EXPECT().RestartDogu(mock.Anything, mock.Anything).Return(nil)
			}

			sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

			resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), tc.req)

			assert.Equal(t, tc.xResponse, resp != nil)
			assert.Equal(t, tc.xResponseCode, status.Code(err))

			if !tc.xResponse {
				mockedConfigProvider.AssertNotCalled(t, "DoguConfig", mock.Anything)
				mockedConfigurationContext.AssertNotCalled(t, "GetOrFalse", mock.Anything)
				mockedConfigurationContext.AssertNotCalled(t, "Set", mock.Anything, mock.Anything)
				mockedDoguRestarter.AssertNotCalled(t, "RestartDogu", mock.Anything, mock.Anything)
			}
		})
	}

	t.Run("Get log level from default config", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)
		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(false, "", nil)
		mockedDescriptionGetter.EXPECT().GetCurrent(mock.Anything, mock.Anything).Return(&core.Dogu{
			Name: "test",
			Configuration: []core.ConfigurationField{
				{
					Name:    "logging/root",
					Default: "WARN",
				},
			},
		}, nil)

		mockedConfigurationContext.EXPECT().Set(mock.Anything, mock.Anything).Return(nil)

		replica := int32(1)
		mockedDeploymentGetter.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(&appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: &replica,
			},
		}, nil)

		mockedDoguRestarter.EXPECT().RestartDogu(mock.Anything, mock.Anything).Return(nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_DEBUG,
		})

		assert.True(t, resp != nil)
		assert.Equal(t, codes.OK, status.Code(err))
	})

	t.Run("Set log level even when no current log level is found", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)
		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(false, "", nil)
		mockedDescriptionGetter.EXPECT().GetCurrent(mock.Anything, mock.Anything).Return(&core.Dogu{
			Name: "test",
			Configuration: []core.ConfigurationField{
				{
					Name: "logging/root",
				},
			},
		}, nil)

		mockedConfigurationContext.EXPECT().Set(mock.Anything, mock.Anything).Return(nil)

		replica := int32(1)
		mockedDeploymentGetter.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(&appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: &replica,
			},
		}, nil)

		mockedDoguRestarter.EXPECT().RestartDogu(mock.Anything, mock.Anything).Return(nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_DEBUG,
		})

		assert.True(t, resp != nil)
		assert.Equal(t, codes.OK, status.Code(err))
	})

	t.Run("No restart when dogu is stopped", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)

		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(true, "DEBUG", nil)

		mockedConfigurationContext.EXPECT().Set(mock.Anything, mock.Anything).Return(nil)

		replica := int32(0)
		mockedDeploymentGetter.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(&appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: &replica,
			},
		}, nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_WARN,
		})

		assert.True(t, resp != nil)
		assert.Equal(t, codes.OK, status.Code(err))

		mockedDoguRestarter.AssertNotCalled(t, "RestartDogu", mock.Anything, mock.Anything)
	})

	t.Run("No restart as log level already set", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)

		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(true, "DEBUG", nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_DEBUG,
		})

		assert.True(t, resp != nil)
		assert.Equal(t, codes.OK, status.Code(err))

		mockedConfigurationContext.AssertNotCalled(t, "Set", mock.Anything)
		mockedDoguRestarter.AssertNotCalled(t, "RestartDogu", mock.Anything, mock.Anything)
	})

	t.Run("Error getting config log level", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)

		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(false, "", errors.New("testError"))

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_DEBUG,
		})

		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))

		mockedDeploymentGetter.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		mockedConfigurationContext.AssertNotCalled(t, "Set", mock.Anything)
		mockedDoguRestarter.AssertNotCalled(t, "RestartDogu", mock.Anything, mock.Anything)
	})

	t.Run("Error getting default log level", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)
		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(false, "", nil)
		mockedDescriptionGetter.EXPECT().GetCurrent(mock.Anything, mock.Anything).Return(nil, errors.New("testError"))

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_DEBUG,
		})

		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))

		mockedDeploymentGetter.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		mockedConfigurationContext.AssertNotCalled(t, "Set", mock.Anything)
		mockedDoguRestarter.AssertNotCalled(t, "RestartDogu", mock.Anything, mock.Anything)
	})

	t.Run("Error setting new log level", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)

		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(false, "", nil)
		mockedConfigurationContext.EXPECT().Set(mock.Anything, mock.Anything).Return(errors.New("testError"))

		mockedDescriptionGetter.EXPECT().GetCurrent(mock.Anything, mock.Anything).Return(&core.Dogu{
			Name: "test",
			Configuration: []core.ConfigurationField{
				{
					Name:    "logging/root",
					Default: "WARN",
				},
			},
		}, nil)

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_DEBUG,
		})

		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))

		mockedDeploymentGetter.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		mockedDoguRestarter.AssertNotCalled(t, "RestartDogu", mock.Anything, mock.Anything)
	})

	t.Run("Error restarting dogu", func(t *testing.T) {
		mockedLogProvider := newMockLogProvider(t)
		mockedConfigProvider := newMockConfigProvider(t)
		mockedDoguRestarter := newMockDoguRestarter(t)
		mockedConfigurationContext := NewMockConfigurationContext(t)
		mockedDescriptionGetter := newMockDoguDescriptionGetter(t)
		mockedDeploymentGetter := newMockDeploymentGetter(t)

		mockedDescriptionGetter.EXPECT().GetCurrent(mock.Anything, mock.Anything).Return(&core.Dogu{
			Name: "test",
			Configuration: []core.ConfigurationField{
				{
					Name:    "logging/root",
					Default: "WARN",
				},
			},
		}, nil)

		mockedConfigProvider.EXPECT().DoguConfig(mock.Anything).Return(mockedConfigurationContext)

		mockedConfigurationContext.EXPECT().GetOrFalse(mock.Anything).Return(false, "", nil)
		mockedConfigurationContext.EXPECT().Set(mock.Anything, mock.Anything).Return(nil)

		replica := int32(1)
		mockedDeploymentGetter.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(&appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: &replica,
			},
		}, nil)

		mockedDoguRestarter.EXPECT().RestartDogu(mock.Anything, mock.Anything).Return(errors.New("testError"))

		sut := NewLoggingService(mockedLogProvider, mockedConfigProvider, mockedDoguRestarter, mockedDescriptionGetter, mockedDeploymentGetter)

		resp, err := sut.ApplyLogLevelWithRestart(context.TODO(), &pb.LogLevelRequest{
			DoguName: "test",
			LogLevel: pb.LogLevel_DEBUG,
		})

		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

}
