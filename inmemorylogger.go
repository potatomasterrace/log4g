package log4g

import (
	"bytes"
	"fmt"
	"sync"
)

type InMemoryLogs [][]interface{}

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
func NewInMemoryLogger() (loggerStream LoggerStream, buffer *InMemoryLogs) {
	lock := sync.Mutex{}
	var logBuffer InMemoryLogs = make([][]interface{}, 0)
	return func(level string, values ...interface{}) {
		lock.Lock()
		defer lock.Unlock()
		data := append([]interface{}{level}, values...)
		logBuffer = append(logBuffer, data)
	}, &logBuffer
}
