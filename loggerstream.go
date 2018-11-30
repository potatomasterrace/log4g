package log4g

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/potatomasterrace/catch"
)

const (
	// FATAL logging level
	FATAL = "[FATAL]"
	// ERROR logging level
	ERROR = "[ERROR]"
	// WARN logging level
	WARN = "[WARN] "
	// INFO logging level
	INFO = "[INFO] "
	// DEBUG logging level
	DEBUG = "[DEBUG]"
	// TRACE logging level
	TRACE = "[TRACE]"
	// ALL logging level
	ALL = "[ALL]  "
)

// LoggerStream is an abstract logger.
type LoggerStream func(level string, values ...interface{})

// PrependTime prepends the time of calls to the logger.
func (ls LoggerStream) PrependTime() LoggerStream {
	return func(level string, values ...interface{}) {
		time := time.Now().Format(time.RFC1123)
		ls(level, append([]interface{}{time}, values...)...)
	}
}

// Prepend the values of loggint to the logger.
func (ls LoggerStream) Prepend(prependValues ...interface{}) LoggerStream {
	return func(level string, values ...interface{}) {
		ls(level, append(prependValues, values...)...)
	}
}

// PrependString the strings to the logger.
func (ls LoggerStream) PrependString(prependedMsgs ...string) LoggerStream {
	prependedValues := make([]interface{}, len(prependedMsgs))
	for i := range prependedMsgs {
		prependedValues[i] = prependedMsgs[i]
	}
	return ls.Prepend(prependedValues...)
}

// FunCall prepend the function call info to the logger.
// The function name is prepended automatically.
// Provide the arguments to log as parameters.
func (ls LoggerStream) FunCall(args ...interface{}) LoggerStream {
	// get Caller name pointer
	fpcs := make([]uintptr, 1)
	runtime.Callers(2, fpcs)
	// get Caller func
	fun := runtime.FuncForPC(fpcs[0])
	// format func name
	funcName := fun.Name()
	// Removing filePath
	if strings.Contains(funcName, ".") {
		if nbPoint := strings.Count(funcName, "."); nbPoint > 0 {
			parts := strings.Split(funcName, ".")[1:]
			funcName = strings.Join(parts, ".")
		}
	}
	header := fmt.Sprintf(" -> %s %v : ", funcName, args)
	return ls.Prepend(header)
}

// Append values to the logger.
func (ls LoggerStream) Append(appendedValues ...interface{}) LoggerStream {
	return func(level string, values ...interface{}) {
		ls(level, append(values, appendedValues...)...)
	}
}

// AppendString append strings to the logger.
func (ls LoggerStream) AppendString(appendedMsgs ...string) LoggerStream {
	appendedValues := make([]interface{}, len(appendedMsgs))
	for i := range appendedMsgs {
		appendedValues[i] = appendedMsgs[i]
	}
	return ls.Append(appendedValues...)
}

// NoPanic intercept an eventual panic and returns it as an error.
func (ls LoggerStream) NoPanic(level string, values ...interface{}) error {
	return catch.Error(func() {
		ls(level, values...)
	})
}

// Filter the logging level.
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

// WithLock adds a lock for concurrent writes.
func (ls LoggerStream) WithLock() LoggerStream {
	lock := &sync.Mutex{}
	return func(level string, values ...interface{}) {
		lock.Lock()
		defer lock.Unlock()
		ls(level, values...)
	}
}

// Async makes the logger asynchronous
func (ls LoggerStream) Async(errorHandler func(error)) LoggerStream {
	return func(level string, values ...interface{}) {
		go func() {
			err := ls.NoPanic(level, values...)
			if err != nil && errorHandler != nil {
				catch.Interface(func() {
					errorHandler(err)
				})
			}
		}()
	}
}

// LoggerFactory dispatch messages and organizes them in topics.
type LoggerFactory func(topic string) LoggerStream

// NoPanic intercept an eventual panic and returns it as an error.
func (lf LoggerFactory) NoPanic(topic string) (LoggerStream, error) {
	var loggerStream LoggerStream
	err := catch.Error(
		func() {
			loggerStream = lf(topic)
		})
	return loggerStream, err
}

// NoPanic intercept an eventual panic and returns it as an error.
func (lf LoggerFactory) LoggerStream() LoggerStream {
	return func(level string, values ...interface{}) {
		lf(level)(level, values...)
	}
}
