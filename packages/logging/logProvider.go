package logging

import "time"

type logLine struct {
	timestamp time.Time
	value     string
}

type logProvider interface {
	getLogs(doguName string, linesCount int) ([]logLine, error)
	queryLogs(doguName string, startDate time.Time, endDate time.Time, filter string) ([]logLine, error)
}
