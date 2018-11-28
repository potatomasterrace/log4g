package log4g

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type loggerCall struct {
	level  string
	values []interface{}
}

func testFunction2(loggerStream LoggerStream) {
	loggerStream = loggerStream.FunctionCall("func2Arg2", "func2Arg1")
	loggerStream(TRACE, "hello world", "3")
}
func testFunction(loggerStream LoggerStream) {
	loggerStream = loggerStream.FunctionCall("func1Arg2", "func1Arg1")
	testFunction2(loggerStream)
}
func TestLoggerStream(t *testing.T) {
	loggerCalls := make([]loggerCall, 0)
	loggerStream := LoggerStream(func(level string, values ...interface{}) {
		loggerCalls = append(loggerCalls, loggerCall{
			level:  level,
			values: values,
		})
	}).PrependTime().Prepend("prepend").Append("append").Filter(ALL)
	loggerStream(ERROR, "hello world", "1")
	loggerStream(INFO, "hello world", "2")
	testFunction(loggerStream)
	loggerStream(ALL, "hello world", "4")
	expectedCalls := []loggerCall{
		loggerCall{level: "[ERROR]",
			values: []interface{}{"Wed, 28 Nov 2018 23:42:26 CET", "prepend", "hello world", "1", "append"}},
		loggerCall{level: "[INFO] ",
			values: []interface{}{"Wed, 28 Nov 2018 23:42:26 CET", "prepend", "hello world", "2", "append"}},
		loggerCall{level: "[TRACE]",
			values: []interface{}{"Wed, 28 Nov 2018 23:42:26 CET", "prepend", "- _/Users/redabourial/Documents/GitHub/log4g.testFunction [func1Arg2 func1Arg1] :", "- _/Users/redabourial/Documents/GitHub/log4g.testFunction2 [func2Arg2 func2Arg1] :", "hello world", "3", "append"}}}

	assert.Equal(t, len(expectedCalls), len(loggerCalls))
	for i, expectedCall := range expectedCalls {
		assert.Equal(t, expectedCall.level, loggerCalls[i].level)
		assert.Equal(t, expectedCall.values[1:], loggerCalls[i].values[1:])
	}
	fmt.Printf("\r\n%#v", loggerCalls)
}

func TestLoggerStreamPanicHandle(t *testing.T) {
	loggerStream := LoggerStream(nil)
	err := loggerStream.NoPanic(WARN, "hello")
	assert.NotNil(t, err)
	loggerFactory := LoggerFactory(nil)
	loggerStream, err = loggerFactory.NoPanic("hello")
	assert.NotNil(t, err)
	assert.Nil(t, loggerStream)
}
