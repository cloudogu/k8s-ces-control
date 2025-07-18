package logging

import (
	"errors"
	"fmt"
	pb "github.com/cloudogu/ces-control-api/generated/logging"
	"strings"
)

// LogLevel is the log level that can be defined for a dogu.
type LogLevel int

const (
	LevelErrorUnspecified LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

// String converts LogLevel type to a string
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "LOG_LEVEL_DEBUG"
	case LevelInfo:
		return "LOG_LEVEL_INFO"
	case LevelWarn:
		return "LOG_LEVEL_WARN"
	case LevelErrorUnspecified:
		return "LOG_LEVEL_ERROR_UNSPECIFIED"
	default:
		return "LOG_LEVEL_WARN"
	}
}

// CreateLogLevelFromProto maps protobuf log level to an internal log level used in application
func CreateLogLevelFromProto(pLevel pb.LogLevel) (LogLevel, error) {
	switch pLevel {
	case pb.LogLevel_LOG_LEVEL_DEBUG:
		return LevelDebug, nil
	case pb.LogLevel_LOG_LEVEL_INFO:
		return LevelInfo, nil
	case pb.LogLevel_LOG_LEVEL_WARN:
		return LevelWarn, nil
	case pb.LogLevel_LOG_LEVEL_ERROR_UNSPECIFIED:
		return LevelErrorUnspecified, nil
	default:
		return LevelErrorUnspecified, fmt.Errorf("unknown log level: %v", pLevel)
	}
}

// CreateLogLevelFromString maps a string to an internal log level used in application
func CreateLogLevelFromString(sLevel string) (LogLevel, error) {
	sLevelUpper := strings.ToUpper(sLevel)

	switch sLevelUpper {
	case LevelErrorUnspecified.String():
		return LevelErrorUnspecified, nil
	case LevelWarn.String():
		return LevelWarn, nil
	case LevelInfo.String():
		return LevelInfo, nil
	case LevelDebug.String():
		return LevelDebug, nil
	default:
		return LevelErrorUnspecified, errors.New("unknown log level")
	}
}
