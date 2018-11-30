package log4g

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsoleLogger(t *testing.T) {
	t.Run("stdout", func(t *testing.T) {
		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		NewConsoleLogger()(INFO, "I", "AM", "FIRST")
		NewConsoleLogger()(INFO, "I", "AM", "SECOND")
		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		// back to normal state
		w.Close()
		os.Stdout = old // restoring the real stdout
		out := <-outC
		expectedOutput := []string{
			"I AM FIRST",
			"I AM SECOND",
		}
		for lineNumber, line := range strings.Split(out, "\r\n")[:2] {
			assert.True(t, strings.Contains(line, expectedOutput[lineNumber]))
		}
	})
	t.Run("stderr", func(t *testing.T) {
		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stderr = w
		NewConsoleLogger()(ERROR, "I", "AM", "FIRST")
		NewConsoleLogger()(ERROR, "I", "AM", "SECOND")
		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		// back to normal state
		w.Close()
		os.Stderr = old // restoring the real stdout
		out := <-outC
		expectedOutput := []string{
			"I AM FIRST",
			"I AM SECOND",
		}
		for lineNumber, line := range strings.Split(out, "\r\n")[:2] {
			assert.True(t, strings.Contains(line, expectedOutput[lineNumber]))
		}
	})
}
