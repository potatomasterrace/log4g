package log4g

import "sync"

func NewInMemoryLogger() (loggerStream LoggerStream, buffer *[][]interface{}) {
	lock := sync.Mutex{}
	logBuffer := make([][]interface{}, 0)
	return func(level string, values ...interface{}) {
		lock.Lock()
		defer lock.Unlock()
		data := append([]interface{}{level}, values...)
		logBuffer = append(logBuffer, data)
	}, &logBuffer
}
