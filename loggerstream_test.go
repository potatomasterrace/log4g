package log4g

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type loggerCall struct {
	level  string
	values []interface{}
}

func TestLoggerStream(t *testing.T) {
	loggerCalls := make([]loggerCall, 0)
	loggerStream := LoggerStream(func(level string, values ...interface{}) {
		loggerCalls = append(loggerCalls, loggerCall{
			level:  level,
			values: values,
		})
	}).PrependTime().Prepend("prepend").Append("append").Filter(TRACE, ALL)
	loggerStream(ERROR, "hello world", "1")
	loggerStream(INFO, "hello world", "2")
	loggerStream(TRACE, "hello world", "3")
	loggerStream(ALL, "hello world", "4")
	expectedCalls := []loggerCall{
		loggerCall{level: "[ERROR]", values: []interface{}{
			"Wed, 28 Nov 2018 19:33:08 CET", "prepend", "hello world", "1", "append"}},
		loggerCall{level: "[INFO]", values: []interface{}{
			"Wed, 28 Nov 2018 19:33:08 CET", "prepend", "hello world", "2", "append"}}}

	assert.Equal(t, len(expectedCalls), len(loggerCalls))
	for i, expectedCall := range expectedCalls {
		assert.Equal(t, expectedCall.level, loggerCalls[i].level)
		assert.Equal(t, expectedCall.values[1:], loggerCalls[i].values[1:])
	}
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
