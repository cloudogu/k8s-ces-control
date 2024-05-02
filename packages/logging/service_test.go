package logging

import (
	"archive/zip"
	"bytes"
	pb "github.com/cloudogu/ces-control-api/generated/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"testing"
	"time"
)

func TestNewLoggingService(t *testing.T) {
	t.Run("should create query clock", func(t *testing.T) {
		// given
		llp := &LokiLogProvider{}

		// when
		sut := NewLoggingService(llp)

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
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = Dogu name should not be empty")
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

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().getLogs("my-dogu", 333).Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(mock.Anything).Return(nil)

		sut := NewLoggingService(mockedLogProvider)

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

		sut := NewLoggingService(mockedLogProvider)

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

		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		mockedLogProvider.EXPECT().queryLogs("my-dogu", time.Time{}, time.Time{}, "").Return(logLines, nil)

		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[0].timestamp), Message: logLines[0].value}).Return(nil)
		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[1].timestamp), Message: logLines[1].value}).Return(nil)
		mockedDoguLogServer.EXPECT().Send(&pb.DoguLogMessage{Timestamp: timestamppb.New(logLines[2].timestamp), Message: logLines[2].value}).Return(nil)

		sut := NewLoggingService(mockedLogProvider)

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

		sut := NewLoggingService(mockedLogProvider)

		// when
		request := &pb.DoguLogMessageQueryRequest{
			DoguName: "",
		}
		err := sut.QueryForDogu(request, mockedDoguLogServer)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = Dogu name should not be empty")
	})

	t.Run("should fail to get logs for error querying the log provider", func(t *testing.T) {
		// given
		mockedLogProvider := newMockLogProvider(t)
		mockedDoguLogServer := newMockDoguLogMessagesQueryServer(t)

		start := time.Unix(1711616504, 0).UTC()
		end := time.Unix(1712131304, 0).UTC()
		filter := "foo=bar"

		mockedLogProvider.EXPECT().queryLogs("my-dogu", start, end, filter).Return(nil, assert.AnError)

		sut := NewLoggingService(mockedLogProvider)

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

		sut := NewLoggingService(mockedLogProvider)

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
