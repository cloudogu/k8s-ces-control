package logging

import (
	"errors"
	"fmt"
	pbLogging "github.com/cloudogu/ces-control-api/generated/logging"
	"testing"
)

func TestCreateLogLevelFromProto(t *testing.T) {
	tests := []struct {
		name     string
		expected LogLevel
		input    pbLogging.LogLevel
		err      error
	}{
		{"Debug", LevelDebug, pbLogging.LogLevel_LOG_LEVEL_DEBUG, nil},
		{"Info", LevelInfo, pbLogging.LogLevel_LOG_LEVEL_INFO, nil},
		{"Warn", LevelWarn, pbLogging.LogLevel_LOG_LEVEL_WARN, nil},
		{"Error", LevelErrorUnspecified, pbLogging.LogLevel_LOG_LEVEL_ERROR_UNSPECIFIED, nil},
		{"Unknown", LevelErrorUnspecified, 100, fmt.Errorf("unknown log level UNKNOWN")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CreateLogLevelFromProto(tt.input)

			if tt.err != nil {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result != tt.expected {
					t.Errorf("Expected log level '%v', got '%v'", tt.expected, result)
				}
			}
		})
	}
}

func TestCreateLevelFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LogLevel
		err      error
	}{
		{"Debug", "LOG_LEVEL_DEBUG", LevelDebug, nil},
		{"Info", "LOG_LEVEL_INFO", LevelInfo, nil},
		{"Warn", "LOG_LEVEL_WARN", LevelWarn, nil},
		{"Error", "LOG_LEVEL_ERROR_UNSPECIFIED", LevelErrorUnspecified, nil},
		{"Unknown", "LOG_LEVEL_UNKNOWN", LevelErrorUnspecified, fmt.Errorf("unknown log level UNKNOWN")},
		{"Empty", "", 0, errors.New("log level string is empty")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CreateLogLevelFromString(tt.input)

			if tt.err != nil {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result != tt.expected {
					t.Errorf("Expected log level '%v', got '%v'", tt.expected, result)
				}
			}
		})
	}
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name     string
		input    LogLevel
		expected string
	}{
		{"Debug", LevelDebug, "LOG_LEVEL_DEBUG"},
		{"Info", LevelInfo, "LOG_LEVEL_INFO"},
		{"Warn", LevelWarn, "LOG_LEVEL_WARN"},
		{"Error", LevelErrorUnspecified, "LOG_LEVEL_ERROR_UNSPECIFIED"},
		{"Unknown", LevelErrorUnspecified, "LOG_LEVEL_ERROR_UNSPECIFIED"},
		{"Unknown", LogLevel(100), "LOG_LEVEL_WARN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
