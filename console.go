package log4g

import (
	"fmt"
	"os"
)

// NewConsoleLogger creates a Logger that outputs to console.
func NewConsoleLogger() Logger {
	return func(level string, values ...interface{}) {
		var file *os.File
		if level == FATAL || level == ERROR {
			file = os.Stderr
		} else {
			file = os.Stdout
		}
		fmt.Fprintf(file, "%s : %v\r\n", level, values)
	}
}
