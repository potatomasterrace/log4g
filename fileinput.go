package log4g

import (
	"bufio"
	"os"
	"sync"
)

// NewFileInput creates a new InputStream.
func NewFileInput(path string) (func() *string, error) {
	lock := sync.Mutex{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	return func() (nextLine *string) {
		lock.Lock()
		defer lock.Unlock()
		if scanner.Scan() {
			text := scanner.Text()
			return &text
		}
		file.Close()
		return nil
	}, nil
}
