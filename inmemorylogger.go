package log4g

import (
	"bytes"
	"fmt"
	"sync"
)

type InMemoryLogs [][]interface{}

func (logs InMemoryLogs) toString(valueFormat string, valueDelimiter string, callDelimiter string) string {
	var buffer bytes.Buffer
	for _, logValues := range logs {
		for _, logValue := range logValues {
			buffer.WriteString(fmt.Sprintf(valueFormat, logValue))
			buffer.WriteString(valueDelimiter)
		}
		buffer.WriteString(callDelimiter)
	}
	return buffer.String()
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
