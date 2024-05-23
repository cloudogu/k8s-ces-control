package logging

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
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

var (
	errMissingDoguName = errors.New("dogu name should not be empty")
)

const loggingKey = "logging/root"

type doguLogMessagesServer interface {
	pb.DoguLogMessages_GetForDoguServer
}

type configProvider interface {
	DoguConfig(dogu string) registry.ConfigurationContext
}

type doguRestarter interface {
	RestartDogu(ctx context.Context, doguName string) error
}

type doguDescriptionGetter interface {
	GetCurrent(ctx context.Context, simpleDoguName string) (*core.Dogu, error)
}

// NewLoggingService creates a new logging service.
func NewLoggingService(provider logProvider, cp configProvider, restarter doguRestarter, descriptionGetter doguDescriptionGetter) *loggingService {
	return &loggingService{
		logProvider:          provider,
		configProvider:       cp,
		doguRestarter:        restarter,
		doguDescriptorGetter: descriptionGetter,
	}
}

type loggingService struct {
	pb.UnimplementedDoguLogMessagesServer
	logProvider          logProvider
	configProvider       configProvider
	doguRestarter        doguRestarter
	doguDescriptorGetter doguDescriptionGetter
}

// QueryForDogu writes dogu log messages into the stream of the given server.
func (s *loggingService) QueryForDogu(request *pb.DoguLogMessageQueryRequest, server pb.DoguLogMessages_QueryForDoguServer) error {
	doguName := request.DoguName
	if doguName == "" {
		return createInternalErr(errMissingDoguName, codes.InvalidArgument)
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
func (s *loggingService) GetForDogu(request *pb.DoguLogMessageRequest, server pb.DoguLogMessages_GetForDoguServer) error {
	linesCount := int(request.LineCount)
	doguName := request.DoguName
	// delegate to an orderly named method because GetForDogu is misleading but cannot be renamed due to the
	// distributed nature of GRPC definitions
	return writeLogLinesToStream(s.logProvider, doguName, linesCount, server)
}

// ApplyLogLevelWithRestart sets the log level for a specific dogu and restarts the dogu if the log level was changed.
func (s *loggingService) ApplyLogLevelWithRestart(ctx context.Context, req *pb.LogLevelRequest) (res *emptypb.Empty, err error) {
	doguName := req.DoguName

	createInternalErrWithCtx := wrapCreateInternalErrWithContext(fmt.Sprintf("error occurred in ApplyLogLevelWithRestart for dogu \"%s\"", doguName))

	if strings.TrimSpace(doguName) == "" {
		return nil, createInternalErrWithCtx(errMissingDoguName, codes.InvalidArgument)
	}

	lLevel, err := CreateLogLevelFromProto(req.GetLogLevel())
	if err != nil {
		return nil, createInternalErrWithCtx(fmt.Errorf("unable to map log level from proto message: %w", err), codes.InvalidArgument)
	}

	restart, err := s.setLogLevel(ctx, doguName, lLevel)
	if err != nil {
		return nil, createInternalErrWithCtx(fmt.Errorf("unable to set log level: %w", err), codes.Internal)
	}

	logrus.Debugf("restart needed for log level change: %v", restart)

	if !restart {
		return &emptypb.Empty{}, nil
	}

	if lErr := s.doguRestarter.RestartDogu(context.WithoutCancel(ctx), doguName); lErr != nil {
		return nil, createInternalErrWithCtx(fmt.Errorf("unable to restart dogu %s after setting new log level: %w", doguName, lErr), codes.Internal)
	}

	logrus.Debugf("Restarted dogu %s", doguName)

	return &emptypb.Empty{}, nil
}

func (s *loggingService) setLogLevel(ctx context.Context, doguName string, l LogLevel) (bool, error) {
	doguConfig := s.configProvider.DoguConfig(doguName)

	currentLogLevel, err := s.getLogLevel(ctx, doguName, doguConfig)
	if err != nil {
		return false, fmt.Errorf("could not get current log level: %w", err)
	}

	if currentLogLevel == l {
		return false, nil
	}

	if lErr := s.writeLogLevel(ctx, doguConfig, l); lErr != nil {
		return false, fmt.Errorf("could not change log level from %s to %s: %w", currentLogLevel, l.String(), err)
	}

	logrus.Debugf("written new log level %s for dogu %s", l.String(), doguName)

	return true, nil
}

// GetLogLevel provides the log level currently set for a specific dogu.
// For this the dogu config is checked as well as the dogu description as fall back.
// When there is no value set for the log level LevelUnknown is returned without an error.
// An error is only returned in case dogu config or dogu description cannot be read.
func (s *loggingService) GetLogLevel(ctx context.Context, doguName string) (LogLevel, error) {
	dConfig := s.configProvider.DoguConfig(doguName)

	return s.getLogLevel(ctx, doguName, dConfig)
}

func (s *loggingService) getLogLevel(ctx context.Context, doguName string, doguConfig registry.ConfigurationContext) (LogLevel, error) {
	currentLogLevelStr, err := s.getConfigLogLevel(ctx, doguConfig)
	if err != nil {
		return LevelUnknown, fmt.Errorf("could not get log level from config: %w", err)
	}

	if currentLogLevelStr == "" {
		logrus.Debugf("config log level is empty, try to get default log level from dogu descrption")

		currentLogLevelStr, err = s.getDefaultLogLevel(ctx, doguName)
		if err != nil {
			return LevelUnknown, fmt.Errorf("could not get default log level from dogu description: %w", err)
		}
	}

	if currentLogLevelStr == "" {
		logrus.Warnf("log level for dogu %s is neither set in config nor description", doguName)
		return LevelUnknown, nil
	}

	logrus.Debugf("current log level from dogu %s is %s", doguName, currentLogLevelStr)

	currentLogLevel, err := CreateLogLevelFromString(currentLogLevelStr)
	if err != nil {
		logrus.Warnf("invalid log level set for dogu %s: %s", doguName, currentLogLevelStr)

		return LevelUnknown, nil
	}

	return currentLogLevel, nil
}

func (s *loggingService) getConfigLogLevel(_ context.Context, dConfig registry.ConfigurationContext) (string, error) {
	_, configLevelStr, err := dConfig.GetOrFalse(loggingKey)
	if err != nil {
		return "", fmt.Errorf("could not receive value from config key: %w", err)
	}

	return configLevelStr, nil
}

func (s *loggingService) getDefaultLogLevel(ctx context.Context, doguName string) (string, error) {
	doguDescription, err := s.doguDescriptorGetter.GetCurrent(ctx, doguName)
	if err != nil {
		return "", fmt.Errorf("could not get dogu description for dogu %s", doguName)
	}

	var defaultLevelStr string

	for _, cfgValue := range doguDescription.Configuration {
		if cfgValue.Name == loggingKey {
			defaultLevelStr = cfgValue.Default
			break
		}
	}

	return defaultLevelStr, nil
}

func (s *loggingService) writeLogLevel(_ context.Context, dConfig registry.ConfigurationContext, l LogLevel) error {
	err := dConfig.Set(loggingKey, l.String())
	if err != nil {
		return fmt.Errorf("could not write to dogu config: %w", err)
	}

	return nil
}

func writeLogLinesToStream(logProvider logProvider, doguName string, linesCount int, server doguLogMessagesServer) error {
	if doguName == "" {
		return createInternalErr(errMissingDoguName, codes.InvalidArgument)
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

func wrapCreateInternalErrWithContext(errCtx string) func(error, codes.Code) error {
	return func(err error, code codes.Code) error {
		contextError := fmt.Errorf("%s : %w", errCtx, err)

		return createInternalErr(contextError, code)
	}
}
