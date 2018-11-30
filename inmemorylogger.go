package log4g

import (
	"bytes"
	"fmt"
)

// InMemoryLogs is an abstraction for logs in memory
type InMemoryLogs [][]interface{}

// StringArray transforms the log values to string seperated by valueDelimiter.
func (logs InMemoryLogs) StringArray(valueDelimiter string) []string {
	lines := make([]string, len(logs))
	for i, logValues := range logs {
		var buffer bytes.Buffer
		for _, logValue := range logValues {
			buffer.WriteString(fmt.Sprint(logValue))
			buffer.WriteString(valueDelimiter)
		}
		lines[i] = buffer.String()
	}
	return lines
}

// NewInMemoryLogger create a logger that outputs values to buffer.
func NewInMemoryLogger() (Logger Logger, buffer *InMemoryLogs) {
	var logBuffer InMemoryLogs = make([][]interface{}, 0)
	return func(level string, values ...interface{}) {
		data := append([]interface{}{level}, values...)
		logBuffer = append(logBuffer, data)
	}, &logBuffer
}
