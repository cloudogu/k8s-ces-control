package logging

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/registry"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"

	pb "github.com/cloudogu/ces-control-api/generated/logging"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	responseMessageMissingDoguname = "dogu name should not be empty"
	loggingKey                     = "logging/root"
)

type doguLogMessagesServer interface {
	pb.DoguLogMessages_GetForDoguServer
}

type configProvider interface {
	DoguConfig(dogu string) registry.ConfigurationContext
}

type doguRestarter interface {
	RestartDogu(ctx context.Context, doguName string) error
}

type logLevel int

const (
	levelDebug logLevel = iota
	levelInfo
	levelWarn
	levelError
)

func (l logLevel) String() string {
	switch l {
	case levelDebug:
		return "DEBUG"
	case levelInfo:
		return "INFO"
	case levelWarn:
		return "WARN"
	case levelError:
		return "ERROR"
	default:
		return "WARN"
	}
}

// NewLoggingService creates a new logging service.
func NewLoggingService(provider logProvider, cp configProvider, restarter doguRestarter) *LoggingService {
	return &LoggingService{
		logProvider:    provider,
		configProvider: cp,
		doguRestarter:  restarter,
	}
}

type LoggingService struct {
	pb.UnimplementedDoguLogMessagesServer
	logProvider    logProvider
	configProvider configProvider
	doguRestarter  doguRestarter
}

// QueryForDogu writes dogu log messages into the stream of the given server.
func (s *LoggingService) QueryForDogu(request *pb.DoguLogMessageQueryRequest, server pb.DoguLogMessages_QueryForDoguServer) error {
	doguName := request.DoguName
	if doguName == "" {
		return status.Error(codes.InvalidArgument, responseMessageMissingDoguname)
	}

	var filter = ""
	if request.Filter != nil {
		filter = request.GetFilter()
	}

	var startDate time.Time
	if request.GetStartDate() != nil {
		startDate = request.GetStartDate().AsTime()
	}

	var endDate time.Time
	if request.GetEndDate() != nil {
		endDate = request.GetEndDate().AsTime()
	}

	logrus.Debugf("retrieving log messages from %s to %s for dogu '%s' with filter %s", startDate, endDate, doguName, filter)

	logLines, err := s.logProvider.queryLogs(doguName, startDate, endDate, filter)
	if err != nil {
		logrus.Errorf("error reading logs: %v", err)
		return createInternalErr(err, codes.InvalidArgument)
	}

	for _, line := range logLines {
		err := server.Send(&pb.DoguLogMessage{
			Timestamp: timestamppb.New(line.timestamp),
			Message:   line.value,
		})
		if err != nil {
			logrus.Errorf("error writing log-lines to stream: %v", err)
			return createInternalErr(err, codes.InvalidArgument)
		}
	}

	return nil
}

// GetForDogu writes dogu log messages into the stream of the given server.
func (s *LoggingService) GetForDogu(request *pb.DoguLogMessageRequest, server pb.DoguLogMessages_GetForDoguServer) error {
	linesCount := int(request.LineCount)
	doguName := request.DoguName
	// delegate to an orderly named method because GetForDogu is misleading but cannot be renamed due to the
	// distributed nature of GRPC definitions
	return writeLogLinesToStream(s.logProvider, doguName, linesCount, server)
}

func (s *LoggingService) SetLogLevel(ctx context.Context, req *pb.LogLevelRequest) (*emptypb.Empty, error) {
	doguName := req.DoguName

	if strings.TrimSpace(doguName) == "" {
		return nil, createInternalErr(errors.New(responseMessageMissingDoguname), codes.InvalidArgument)
	}

	lLevel, err := mapLogLevel(req.GetLogLevel())
	if err != nil {
		return nil, createInternalErr(fmt.Errorf("unable to map log level from proto message: %w", err), codes.InvalidArgument)
	}

	restart, err := s.setLogLevel(ctx, doguName, lLevel)
	if err != nil {
		return nil, createInternalErr(fmt.Errorf("unable to set log level: %w", err), codes.Internal)
	}

	if !restart {
		return &emptypb.Empty{}, nil
	}

	if lErr := s.doguRestarter.RestartDogu(ctx, doguName); lErr != nil {
		return nil, createInternalErr(fmt.Errorf("unable to restart dogu %s after setting new log level: %w", doguName, lErr), codes.Internal)
	}

	return &emptypb.Empty{}, nil
}

func (s *LoggingService) GetLogLevel(doguName string) (string, error) {
	dConfig := s.configProvider.DoguConfig(doguName)

	currentLevel, err := dConfig.Get(loggingKey)
	if err != nil {
		return "", fmt.Errorf("could not get current log level: %w", err)
	}

	return currentLevel, nil
}

func (s *LoggingService) setLogLevel(_ context.Context, doguName string, l logLevel) (bool, error) {
	dConfig := s.configProvider.DoguConfig(doguName)

	currentLevel, err := dConfig.Get(loggingKey)
	if err != nil {
		return false, fmt.Errorf("could not get current log level: %w", err)
	}

	if strings.EqualFold(currentLevel, l.String()) {
		return false, nil
	}

	err = dConfig.Set(loggingKey, l.String())
	if err != nil {
		return false, fmt.Errorf("could not change log level from %s to %s: %w", currentLevel, l.String(), err)
	}

	return true, nil
}

func mapLogLevel(pLevel pb.LogLevel) (logLevel, error) {
	switch pLevel {
	case pb.LogLevel_DEBUG:
		return levelDebug, nil
	case pb.LogLevel_INFO:
		return levelInfo, nil
	case pb.LogLevel_WARN:
		return levelWarn, nil
	case pb.LogLevel_ERROR:
		return levelError, nil
	default:
		return 0, fmt.Errorf("unknown loglevel: %v", pLevel)
	}
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
