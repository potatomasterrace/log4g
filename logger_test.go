package log4g

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type loggerCall struct {
	level  string
	values []interface{}
}

func testFunction2(Logger Logger) {
	Logger = Logger.FunCall("func2Arg2", "func2Arg1")
	Logger(TRACE, "fc1", "fc2")
}
func testFunction(Logger Logger) {
	Logger = Logger.FunCall("func1Arg2", "func1Arg1")
	testFunction2(Logger)
}
func TestLogger(t *testing.T) {
	loggerCalls := make([]loggerCall, 0)
	Logger := Logger(func(level string, values ...interface{}) {
		loggerCalls = append(loggerCalls, loggerCall{
			level:  level,
			values: values,
		})
	}).PrependTime().Prepend("prepend").Append("append").Filter(ALL)
	Logger(ERROR, "msg1", "mgs2")
	Logger.AppendString("a2", "a3").PrependString("p1", "p2")(INFO, "msg3", "msg4")
	testFunction(Logger)
	Logger(ALL, "f1", "f2")
	expectedCalls := []loggerCall{
		loggerCall{level: "[ERROR]",
			values: []interface{}{"Thu, 29 Nov 2018 21:53:07 CET", "prepend", "msg1", "mgs2", "append"}},
		loggerCall{level: "[INFO] ",
			values: []interface{}{"Thu, 29 Nov 2018 21:53:07 CET", "prepend", "p1", "p2", "msg3", "msg4", "a2", "a3", "append"}},
		loggerCall{level: "[TRACE]",
			values: []interface{}{"Thu, 29 Nov 2018 21:53:07 CET", "prepend", " -> testFunction [func1Arg2 func1Arg1] : ", " -> testFunction2 [func2Arg2 func2Arg1] : ", "fc1", "fc2", "append"}}}
	assert.Equal(t, len(expectedCalls), len(loggerCalls))
	for i, expectedCall := range expectedCalls {
		assert.Equal(t, expectedCall.level, loggerCalls[i].level)
		assert.Equal(t, expectedCall.values[1:], loggerCalls[i].values[1:])
	}
	// fmt.Printf("\r\n%#v", loggerCalls)
}

func TestAsync(t *testing.T) {
	panicHandler := func(err error) {
		assert.NotNil(t, err)
	}
	Logger(nil).Async(panicHandler)
	Logger(nil).Async(nil)
}
func TestWithLock(t *testing.T) {
	start := time.Now()
	waitime := 100
	logger := Logger(func(string, ...interface{}) {
		time.Sleep(1 * time.Millisecond)
	}).WithLock()
	var wg sync.WaitGroup
	wg.Add(waitime)
	for i := 0; i < waitime; i++ {
		go func() {
			defer wg.Done()
			logger(INFO, i)
		}()
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(start))
	assert.True(t, time.Now().Sub(start) > time.Duration(waitime)*time.Millisecond)
}
func TestLoggerPanicHandle(t *testing.T) {
	Logger := Logger(nil)
	err := Logger.NoPanic(WARN, "hello")
	assert.NotNil(t, err)
	loggerFactory := LoggerFactory(nil)
	Logger, err = loggerFactory.NoPanic("hello")
	assert.NotNil(t, err)
	assert.Nil(t, Logger)
}
