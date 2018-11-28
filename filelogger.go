package log4g

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

// FileWritingContext stores the data for writing logged values.
type FileWritingContext struct {
	sync.Mutex
	File             *os.File
	PrependTime      string
	FormatingFunc    func(value interface{}) string
	LoggerStream     LoggerStream
	CallDelimiter    string
	ValuesDelimiters string
	Path             string
}

// FormatValues format the logger values into a line to write on the log file
func (fwc FileWritingContext) FormatValues(level string, values ...interface{}) string {
	var buffer bytes.Buffer
	buffer.WriteString(level)
	for _, value := range values {
		buffer.WriteString(fwc.ValuesDelimiters)
		buffer.WriteString(fwc.FormatingFunc(value))
	}
	buffer.WriteString(fwc.CallDelimiter)
	return buffer.String()
}

// Close the underlying file.
func (wc *FileWritingContext) Close() error {
	if wc.File != nil {
		return wc.File.Close()
	} else {
		return fmt.Errorf("trying to close already close log file %s", wc.Path)
	}
	wc.File = nil
	return nil
}

// GetFileLogger returns the logger stream for the file.
func (wc *FileWritingContext) Init() error {
	file, err := os.OpenFile(wc.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	wc.File = file
	lock := sync.Mutex{}
	wc.LoggerStream = func(level string, values ...interface{}) {
		byts := wc.FormatValues(level, values...)
		lock.Lock()
		defer lock.Unlock()
		_, err := file.WriteString(byts)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

type DirLogger struct {
	DirContext FileWritingContext
	OpenFiles  []FileWritingContext
}

func (dirLogger DirLogger) topicToPath(topic string) string {
	return fmt.Sprintf("%s/%s", dirLogger.DirContext.Path, topic)
}

func (dirLogger DirLogger) find(topic string) *FileWritingContext {
	for _, openFile := range dirLogger.OpenFiles {
		if openFile.Path == dirLogger.topicToPath(topic) {
			return &openFile
		}
	}
	return nil
}

func (dirLogger *DirLogger) Get(topic string) *FileWritingContext {
	currentLogger := dirLogger.find(topic)
	if currentLogger != nil {
		return currentLogger
	}
	fwc := dirLogger.DirContext
	fwc.Path = dirLogger.topicToPath(topic)
	err := fwc.Init()
	if err != nil {
		panic(err)
	}
	dirLogger.OpenFiles = append(dirLogger.OpenFiles, fwc)
	return &fwc
}

func (dirLogger *DirLogger) Close() error {
	if dirLogger.OpenFiles == nil || len(dirLogger.OpenFiles) == 0 {
		return fmt.Errorf("no files to close")
	}
	for _, openFile := range dirLogger.OpenFiles {
		openFile.Close()
	}
	dirLogger.OpenFiles = nil
	return nil
}
func (dirLogger *DirLogger) GetLoggerFactory() LoggerFactory {
	lock := sync.Mutex{}
	return func(topic string) LoggerStream {
		lock.Lock()
		defer lock.Unlock()
		return dirLogger.Get(topic).LoggerStream
	}
}

// NewDirLogger returns a logger that dispatch topic in a folder files.
func NewDirLogger(dirPath string, dirContext FileWritingContext) DirLogger {
	dirLogger := DirLogger{
		DirContext: dirContext,
		OpenFiles:  make([]FileWritingContext, 0),
	}
	return dirLogger
}
