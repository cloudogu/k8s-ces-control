package logging

import "time"

type logLine struct {
	timestamp time.Time
	value     string
}

type logProvider interface {
	getLogs(doguName string, linesCount int) ([]logLine, error)
	getLogsInRange(doguName string, startDate string, endDate string, filter string) ([]logLine, error)
}
