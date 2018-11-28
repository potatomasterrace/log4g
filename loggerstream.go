package log4g

import (
	"fmt"
	"runtime"
	"strings"
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
		ls(level, append([]interface{}{time}, values...)...)
	}
}

func (ls LoggerStream) Prepend(prependedMsgs ...interface{}) LoggerStream {
	return func(level string, values ...interface{}) {
		ls(level, append(prependedMsgs, values...)...)
	}
}

func (ls LoggerStream) FunctionCall(args ...interface{}) LoggerStream {
	// get Caller name pointer
	fpcs := make([]uintptr, 1)
	runtime.Callers(2, fpcs)
	// get Caller func
	fun := runtime.FuncForPC(fpcs[0])
	// format func name
	funcName := fun.Name()
	if strings.Contains(funcName, ".") {
		if afterPoint := strings.Split(funcName, ".")[1]; len(afterPoint) > 0 {
			funcName = afterPoint
		}
	}
	header := fmt.Sprintf("- %s %s :", funcName, args)
	return ls.Prepend(header)
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
