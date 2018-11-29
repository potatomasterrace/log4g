package log4g

import (
	"bufio"
	"os"
	"sync"

	"github.com/potatomasterrace/catch"
)

// InputStream in an shortcup for reading data line by line.
// It can panic (see NoPanic method)
type InputStream func() *string

// WithLock adds a lock for concurrent reads.
func (is InputStream) WithLock() InputStream {
	lock := &sync.Mutex{}
	return func() *string {
		var line *string
		err := catch.Error(func() {
			lock.Lock()
			defer lock.Unlock()
			line = is()
		})
		if err != nil {
			panic(err)
		}
		return line
	}
}

// NoPanic intercept an eventual panic and returns it as an error.
func (is InputStream) NoPanic() (*string, error) {
	var line *string
	err := catch.Error(func() {
		line = is()
	})
	return line, err
}

// NewFileInput creates a new InputStream.
func NewFileInput(path string) (InputStream, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	return func() (nextLine *string) {
		if scanner.Scan() {
			text := scanner.Text()
			return &text
		}
		file.Close()
		return nil
	}, nil
}
