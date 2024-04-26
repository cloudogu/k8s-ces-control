package logging

import (
	"archive/zip"
	"bytes"
	"fmt"
	"time"

	pb "github.com/cloudogu/ces-control-api/generated/logging"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	responseMessageMissingDoguname = "Dogu name should not be empty"
)

type doguLogMessagesServer interface {
	pb.DoguLogMessages_GetForDoguServer
}

// NewLoggingService creates a new logging service.
func NewLoggingService(provider logProvider) *loggingService {
	return &loggingService{logProvider: provider}
}

type loggingService struct {
	pb.UnimplementedDoguLogMessagesServer
	logProvider logProvider
}

// GetForDogu writes dogu log messages into the stream of the given server.
func (s *loggingService) GetForDogu(request *pb.DoguLogMessageRequest, server pb.DoguLogMessages_GetForDoguServer) error {
	linesCount := int(request.LineCount)
	doguName := request.DoguName
	// delegate to an orderly named method because GetForDogu is misleading but cannot be renamed due to the
	// distributed nature of GRPC definitions
	return writeLogLinesToStream(s.logProvider, doguName, linesCount, server)
}

// GetForDoguWithDate writes dogu log messages into the stream of the given server.
func (s *loggingService) GetForDoguWithDate(request *pb.DoguLogMessageWithDateRequest, server pb.DoguLogMessages_GetForDoguWithDateServer) error {
	doguName := request.DoguName
	startDate := request.StartDate
	endDate := request.EndDate
	filter := request.Filter
	// delegate to an orderly named method because GetForDoguWithDate is misleading but cannot be renamed due to the
	// distributed nature of GRPC definitions
	return writeLogLinesToStreamWithDate(s.logProvider, doguName, startDate, endDate, filter, server)
}

func writeLogLinesToStream(logProvider logProvider, doguName string, linesCount int, server doguLogMessagesServer) error {
	if doguName == "" {
		return status.Error(codes.InvalidArgument, responseMessageMissingDoguname)
	}
	logrus.Debugf("retrieving %d line(s) of log messages for dogu '%s'", linesCount, doguName)

	logLines, err := logProvider.getLogs(doguName, linesCount)
	if err != nil {
		logrus.Errorf("error reading logs: %v", err)
		return createInternalErr(err, codes.InvalidArgument)
	}

	compressedMessagesBytes, err := compressMessages(doguName, logLines)
	if err != nil {
		logrus.Errorf("error compressing message: %v", err)
		return createInternalErr(err, codes.Internal)
	}

	err = stream.WriteToStream(compressedMessagesBytes, server)
	if err != nil {
		logrus.Errorf("error writing logs to stream: %v", err)
		return createInternalErr(err, codes.Internal)
	}

	return nil
}

func writeLogLinesToStreamWithDate(logProvider logProvider, doguName string, startDate string, endDate string, filter string, server doguLogMessagesServer) error {
	if doguName == "" {
		return status.Error(codes.InvalidArgument, responseMessageMissingDoguname)
	}
	logrus.Debugf("retrieving log messages from %s to %s for dogu '%s' with filter %s", startDate, endDate, doguName, filter)

	logLines, err := logProvider.getLogsInRange(doguName, startDate, endDate, filter)
	if err != nil {
		logrus.Errorf("error reading logs: %v", err)
		return createInternalErr(err, codes.InvalidArgument)
	}

	compressedMessagesBytes, err := compressMessages(doguName, logLines)
	if err != nil {
		logrus.Errorf("error compressing message: %v", err)
		return createInternalErr(err, codes.Internal)
	}

	err = stream.WriteToStream(compressedMessagesBytes, server)
	if err != nil {
		logrus.Errorf("error writing logs to stream: %v", err)
		return createInternalErr(err, codes.Internal)
	}

	return nil
}

func compressMessages(doguName string, logLines []logLine) ([]byte, error) {
	if len(logLines) <= 0 {
		return []byte{}, nil
	}

	compressedMessages := bytes.NewBuffer(nil)
	zipWriter := zip.NewWriter(compressedMessages)
	logFileMetadata := zip.FileHeader{
		Name:     fmt.Sprintf("%s.log", doguName),
		Modified: time.Now(),
		Method:   zip.Deflate,
	}

	writer, err := zipWriter.CreateHeader(&logFileMetadata)
	if err != nil {
		return nil, err
	}

	totalWrittenBytes := 0
	for _, line := range logLines {
		writtenBytes, err := writer.Write([]byte(line.value + "\n"))
		if err != nil {
			return nil, err
		}
		totalWrittenBytes += writtenBytes
	}

	_ = zipWriter.Close()
	logrus.Debugf("wrote %d byte(s) to archive", totalWrittenBytes)
	return compressedMessages.Bytes(), nil
}

func createInternalErr(err error, code codes.Code) error {
	logrus.Error(err)
	return status.Error(code, err.Error())
}
