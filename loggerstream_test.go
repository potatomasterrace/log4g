package log4g

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type loggerCall struct {
	level  string
	values []interface{}
}

func testFunction2(loggerStream LoggerStream) {
	loggerStream = loggerStream.FunCall("func2Arg2", "func2Arg1")
	loggerStream(TRACE, "fc1", "fc2")
}
func testFunction(loggerStream LoggerStream) {
	loggerStream = loggerStream.FunCall("func1Arg2", "func1Arg1")
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
	loggerStream(ERROR, "msg1", "mgs2")
	loggerStream.AppendString("a2", "a3").PrependString("p1", "p2")(INFO, "msg3", "msg4")
	testFunction(loggerStream)
	loggerStream(ALL, "f1", "f2")
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

func TestLoggerStreamPanicHandle(t *testing.T) {
	loggerStream := LoggerStream(nil)
	err := loggerStream.NoPanic(WARN, "hello")
	assert.NotNil(t, err)
	loggerFactory := LoggerFactory(nil)
	loggerStream, err = loggerFactory.NoPanic("hello")
	assert.NotNil(t, err)
	assert.Nil(t, loggerStream)
}
