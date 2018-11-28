package log4g

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInputStream(t *testing.T) {
	path := "./testdata/testfile"
	os.Remove(path)
	t.Run("outputstream", func(t *testing.T) {
		fwc := FileWritingContext{
			FormatingFunc: func(value interface{}) string {
				return fmt.Sprintf("%s", value)
			},
			CallDelimiter:    "\r\n",
			ValuesDelimiters: " | ",
			Path:             "./testdata/testfile",
		}
		t.Run("initialisation", func(t *testing.T) {
			err := fwc.Init()
			assert.Nil(t, err)
		})
		fo := fwc.LoggerStream
		// stays opens
		t.Run("writing", func(t *testing.T) {
			fo("hello", "world", "1")
			fo.Prepend("hello")("world", "2")
			err := fo.NoPanic("hello", "world", "3")
			assert.Nil(t, err)
		})
		t.Run("Closing", func(t *testing.T) {
			// closing file
			err := fwc.Close()
			assert.Nil(t, err)
			err = fwc.Close()
			assert.NotNil(t, err)
		})
		t.Run("error handling", func(t *testing.T) {
			time.Sleep(200 * time.Millisecond)
			err := fo.NoPanic("hello", "world", "4")
			assert.NotNil(t, err)
		})
	})
	t.Run("inputstream", func(t *testing.T) {
		is, err := NewFileInput(path)
		assert.Nil(t, err)
		lines := make([]string, 0)
		for line := is(); line != nil; line = is() {
			lines = append(lines, *line)
		}
		assert.Equal(t, []string{"hello | world | 1", "world | hello | 2", "hello | world | 3"}, lines)
	})
}

func TestDirLogger(t *testing.T) {
	folderpath := "./testdata/logs/"
	os.RemoveAll(folderpath)
	os.Mkdir(folderpath, os.ModePerm)
	dirLogger := NewDirLogger(folderpath, FileWritingContext{
		FormatingFunc:    func(v interface{}) string { return fmt.Sprintf("%s", v) },
		CallDelimiter:    "\r\n",
		ValuesDelimiters: " , ",
		Path:             folderpath,
	})
	loggerFactory := dirLogger.GetLoggerFactory()
	loggerFactory("file1")
	fo1 := loggerFactory("file1")
	fo2 := loggerFactory("file2")
	// stays opens
	t.Run("writing", func(t *testing.T) {
		fo1("hello", "world", "1")
		fo1("foo", "1")
		fo2("hello", "world", "2")
		fo2("foo", "2")
	})
	t.Run("Closing", func(t *testing.T) {
		// closing file
		err := dirLogger.Close()
		assert.Nil(t, err)
		err = dirLogger.Close()
		assert.NotNil(t, err)
	})
	t.Run("error handling", func(t *testing.T) {
		time.Sleep(200 * time.Millisecond)
		err := fo1.NoPanic("hello", "world", "4")
		assert.NotNil(t, err)
	})
	t.Run("inputstream", func(t *testing.T) {
		t.Run("file1", func(t *testing.T) {
			is, err := NewFileInput("./testdata/logs/file1")
			assert.Nil(t, err)
			lines := make([]string, 0)
			for line := is(); line != nil; line = is() {
				lines = append(lines, *line)
			}
			assert.Equal(t, []string{"hello , world , 1", "foo , 1"}, lines)
		})
		t.Run("file2", func(t *testing.T) {
			is, err := NewFileInput("./testdata/logs/file2")
			assert.Nil(t, err)
			lines := make([]string, 0)
			for line := is(); line != nil; line = is() {
				lines = append(lines, *line)
			}
			assert.Equal(t, []string{"hello , world , 2", "foo , 2"}, lines)
		})
	})

}
