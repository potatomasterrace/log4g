package log4g

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/potatomasterrace/catch"
)

// Logger is an abstract logger.
type Logger func(level string, values ...interface{})

// MockLogger returns a mock of a logger.
func T() Logger {
	return func(level string, values ...interface{}) {
		return
	}
}

// PrependTime prepends the time of calls to the logger.
func (logger Logger) PrependTime() Logger {
	return func(level string, values ...interface{}) {
		time := time.Now().Format(time.RFC1123)
		logger(level, append([]interface{}{time}, values...)...)
	}
}

// PrependGoRoutines prepends the current number of running goroutines.
func (logger Logger) PrependGoRoutines() Logger {
	return func(level string, values ...interface{}) {
		msg := fmt.Sprint("[ Go routines : ", runtime.NumGoroutine(), " ]")
		logger(level, append([]interface{}{msg}, values...)...)
	}
}

// Prepend the values of loggint to the logger.
func (logger Logger) Prepend(prependValues ...interface{}) Logger {
	return func(level string, values ...interface{}) {
		logger(level, append(prependValues, values...)...)
	}
}

// AppendString append strings to the logger.
func (logger Logger) D(appendedMsgs ...string) Logger {
	appendedValues := make([]interface{}, len(appendedMsgs))
	for i := range appendedMsgs {
		appendedValues[i] = appendedMsgs[i]
	}
	return logger.Append(appendedValues...)
}

// PrependString the strings to the logger.
func (logger Logger) PrependString(prependedMsgs ...string) Logger {
	prependedValues := make([]interface{}, len(prependedMsgs))
	for i := range prependedMsgs {
		prependedValues[i] = prependedMsgs[i]
	}
	return logger.Prepend(prependedValues...)
}

// FunCall prepend the function call info to the logger.
// The function name is prepended automatically.
// Provide the arguments to log as parameters.
func (logger Logger) FunCall(args ...interface{}) Logger {
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
	return logger.Prepend(header)
}

// DetailedFunCall Provide the arguments to log as parameters.
func (logger Logger) DetailedFunCall(args ...interface{}) Logger {
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
	header := fmt.Sprintf(" -> %s %v %d : ", funcName, args, fun.Entry())
	return logger.Prepend(header)
}

// Append values to the logger.
func (logger Logger) Append(appendedValues ...interface{}) Logger {
	return func(level string, values ...interface{}) {
		logger(level, append(values, appendedValues...)...)
	}
}

// AppendString append strings to the logger.
func (logger Logger) AppendString(appendedMsgs ...string) Logger {
	appendedValues := make([]interface{}, len(appendedMsgs))
	for i := range appendedMsgs {
		appendedValues[i] = appendedMsgs[i]
	}
	return logger.Append(appendedValues...)
}

// NoPanic intercept an eventual panic and returns it as an error.
func (logger Logger) NoPanic(level string, values ...interface{}) error {
	return catch.Error(func() {
		logger(level, values...)
	})
}

// Filter the logging level.
func (logger Logger) Filter(filteredLevels ...string) Logger {
	return func(level string, values ...interface{}) {
		for _, filteredLevel := range filteredLevels {
			if filteredLevel == level {
				return
			}
		}
		logger(level, values...)
	}
}

// WithLock adds a lock for concurrent writes.
func (logger Logger) WithLock() Logger {
	lock := &sync.Mutex{}
	return func(level string, values ...interface{}) {
		lock.Lock()
		defer lock.Unlock()
		logger(level, values...)
	}
}

// Async makes the logger asynchronous
func (logger Logger) Async(errorHandler func(error)) Logger {
	return func(level string, values ...interface{}) {
		go func() {
			err := logger.NoPanic(level, values...)
			if err != nil && errorHandler != nil {
				catch.Interface(func() {
					errorHandler(err)
				})
			}
		}()
	}
}

// LoggerFactory dispatch messages and organizes them in topics.
type LoggerFactory func(topic string) Logger

// NoPanic intercept an eventual panic and returns it as an error.
func (lf LoggerFactory) NoPanic(topic string) (Logger, error) {
	var Logger Logger
	err := catch.Error(
		func() {
			Logger = lf(topic)
		})
	return Logger, err
}

// Logger transforms a logger factory to logger.
func (lf LoggerFactory) Logger() Logger {
	return func(level string, values ...interface{}) {
		lf(level)(level, values...)
	}
}

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

// LoggerStream is an abstraction for logging to a level
type LoggerStream func(values ...interface{})

// Fatal is pretty self explanatory and repetitive
func (logger Logger) Fatal(values ...interface{}) {
	logger(FATAL, values...)
}

// Error repetetive did you say ?
func (logger Logger) Error(values ...interface{}) {
	logger(ERROR, values...)
}

// Warn really, how repetitive ?
func (logger Logger) Warn(values ...interface{}) {
	logger(WARN, values...)
}

// Info not enough to keep me from writing non-sense
func (logger Logger) Info(values ...interface{}) {
	logger(INFO, values...)
}

// Debug just to avoid warnings
func (logger Logger) Debug(values ...interface{}) {
	logger(DEBUG, values...)
}

// Trace now this is getting boring
func (logger Logger) Trace(values ...interface{}) {
	logger(TRACE, values...)
}

// All and that's all !
func (logger Logger) All(values ...interface{}) {
	logger(ALL, values...)
}
