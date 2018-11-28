package log4g

import (
	"time"

	"github.com/potatomasterrace/catch"
)

const (
	FATAL = "[FATAL]"
	ERROR = "[ERROR]"
	WARN  = "[WARN] "
	INFO  = "[INFO] "
	DEBUG = "[DEBUG]"
	TRACE = "[TRACE]"
	ALL   = "[ALL]  "
)

type LoggerStream func(level string, values ...interface{})

func (ls LoggerStream) PrependTime() LoggerStream {
	return func(level string, values ...interface{}) {
		time := time.Now().Format(time.RFC1123)
		ls.Prepend(time)(level, values...)
	}
}

func (ls LoggerStream) Prepend(prependedMsgs ...interface{}) LoggerStream {
	return func(level string, values ...interface{}) {
		ls(level, append(prependedMsgs, values...)...)
	}
}

func (ls LoggerStream) Append(appendedMsgs ...interface{}) LoggerStream {
	return func(level string, values ...interface{}) {
		ls(level, append(values, appendedMsgs...)...)
	}
}
func (ls LoggerStream) NoPanic(level string, values ...interface{}) error {
	return catch.Error(func() {
		ls(level, values...)
	})
}

func (ls LoggerStream) Filter(filteredLevels ...string) LoggerStream {
	return func(level string, values ...interface{}) {
		for _, filteredLevel := range filteredLevels {
			if filteredLevel == level {
				return
			}
		}
		ls(level, values...)
	}
}

// LoggerFactory dispatch messages and organizes them in topics.
type LoggerFactory func(topic string) LoggerStream

func (lf LoggerFactory) NoPanic(topic string) (LoggerStream, error) {
	var loggerStream LoggerStream
	err := catch.Error(
		func() {
			loggerStream = lf(topic)
		})
	return loggerStream, err
}
