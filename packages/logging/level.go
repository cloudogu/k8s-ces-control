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
	LevelUnknown LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

// String converts LogLevel type to a string
func (l LogLevel) String() string {
	switch l {
	case LevelUnknown:
		return "UNKNOWN"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "WARN"
	}
}

// CreateLogLevelFromProto maps protobuf log level to an internal log level used in application
func CreateLogLevelFromProto(pLevel pb.LogLevel) (LogLevel, error) {
	switch pLevel {
	case pb.LogLevel_DEBUG:
		return LevelDebug, nil
	case pb.LogLevel_INFO:
		return LevelInfo, nil
	case pb.LogLevel_WARN:
		return LevelWarn, nil
	case pb.LogLevel_ERROR:
		return LevelError, nil
	default:
		return LevelUnknown, fmt.Errorf("unknown log level: %v", pLevel)
	}
}

// CreateLogLevelFromString maps a string to an internal log level used in application
func CreateLogLevelFromString(sLevel string) (LogLevel, error) {
	sLevelUpper := strings.ToUpper(sLevel)

	switch sLevelUpper {
	case LevelError.String():
		return LevelError, nil
	case LevelWarn.String():
		return LevelWarn, nil
	case LevelInfo.String():
		return LevelInfo, nil
	case LevelDebug.String():
		return LevelDebug, nil
	default:
		return LevelUnknown, errors.New("unknown log level")
	}
}
