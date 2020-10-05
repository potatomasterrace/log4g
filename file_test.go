package log4g

import (
	"fmt"
	"os"
	"testing"

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
		fo := fwc.Logger
		// stays opens
		fo("hello", "world", "1")
		fo.Prepend("hello")("world", "2")
		err := fo.NoPanic("hello", "world", "3")
		assert.Nil(t, err)
		// closing file
		err = fwc.Close()
		assert.Nil(t, err)
		err = fwc.Close()
		assert.NotNil(t, err)
		err = fo.NoPanic("hello", "world", "4")
		assert.NotNil(t, err)
	})
	t.Run("inputstream", func(t *testing.T) {
		is, err := NewFileInput(path)
		is = is.WithLock()
		assert.Nil(t, err)
		lines := make([]string, 0)
		for line := is(); line != nil; line = is() {
			lines = append(lines, *line)
		}
		assert.Equal(t, []string{"hello | world | 1", "world | hello | 2", "hello | world | 3"}, lines)
		is = nil
		_, err = is.NoPanic()
		assert.NotNil(t, err)
	})
}

func TestDirLogger(t *testing.T) {
	folderpath := "./testdata/logs/"
	os.RemoveAll(folderpath)
	dirLogger, err := NewDirLogger(FileWritingContext{
		FormatingFunc:    func(v interface{}) string { return fmt.Sprintf("%s", v) },
		CallDelimiter:    "\r\n",
		ValuesDelimiters: " , ",
		Path:             folderpath,
	})
	assert.Nil(t, err)
	loggerFactory := dirLogger.GetLoggerFactory()
	// doesn't panic on multiple access
	loggerFactory("file1")
	fo1 := loggerFactory("file1")
	fo2 := loggerFactory("file2")
	fo1("hello", "world", "1")
	fo2("hello", "world", "2")
	loggerFactory.Logger()("file1", "foo", "1")
	loggerFactory.Logger()("file2", "foo", "2")
	errs := dirLogger.Close()
	assert.Equal(t, len(errs), 0)
	dirLogger.OpenFiles = []FileWritingContext{
		FileWritingContext{
			Path: "unexisting",
		},
	}
	errs = dirLogger.Close()
	assert.Equal(t, len(errs), 1)
	err = fo1.NoPanic("hello", "world", "4")
	assert.NotNil(t, err)
	t.Run("inputstream", func(t *testing.T) {
		t.Run("file1", func(t *testing.T) {
			is, err := NewFileInput("./testdata/logs/file1")
			assert.Nil(t, err)
			lines := make([]string, 0)
			for line := is(); line != nil; line = is() {
				lines = append(lines, *line)
			}
			assert.Equal(t, []string{"hello , world , 1", "file1 , foo , 1"}, lines)
		})
		t.Run("file2", func(t *testing.T) {
			is, err := NewFileInput("./testdata/logs/file2")
			assert.Nil(t, err)
			lines := make([]string, 0)
			for line := is(); line != nil; line = is() {
				lines = append(lines, *line)
			}
			assert.Equal(t, []string{"hello , world , 2", "file2 , foo , 2"}, lines)
		})
	})

}
