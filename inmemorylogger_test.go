package log4g

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryLogger(t *testing.T) {
	loggerStream, buffer := NewInMemoryLogger()
	loggerStream.PrependTime()(WARN, "hello", "1")
	loggerStream(TRACE, "world", "2")
	assert.True(t, strings.Contains(buffer.toString("%s", ",", "\r\n"), "CET,hello,1"))
	assert.True(t, strings.Contains(buffer.toString("%s", ",", "\r\n"), ",\r\n[TRACE],world,2,\r\n"))
}
