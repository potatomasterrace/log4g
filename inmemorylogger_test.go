package log4g

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryLogger(t *testing.T) {
	loggerStream, buffer := NewInMemoryLogger()
	loggerStream.PrependTime()(WARN, "hello", "1")
	loggerStream(TRACE, "world", "2")
	assert.Equal(t, 2, len(*buffer))
	assert.Equal(t, []interface{}([]interface{}{"hello", "1"}), (*buffer)[0][2:])
	assert.Equal(t, []interface{}{"[TRACE]", "world", "2"}, (*buffer)[1])
}
