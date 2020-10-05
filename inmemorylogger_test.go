package log4g

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryLogger(t *testing.T) {
	Logger, buffer := NewInMemoryLogger()
	Logger.PrependTime()(WARN, "hello", "1")
	Logger(TRACE, "world", "2")
	loggedLines := buffer.StringArray(" ")
	assert.Equal(t, len(loggedLines), 2)
	assert.Contains(t, loggedLines[0], "[WARN]  ")
	assert.Contains(t, loggedLines[0], time.Now().Format(time.RFC1123)[:10])
	assert.Contains(t, loggedLines[0], "hello 1")
	assert.Contains(t, loggedLines[1], "world 2")
}
