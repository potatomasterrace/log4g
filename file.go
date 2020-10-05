package log4g

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
)

// FileWritingContext stores the data for writing logged values.
type FileWritingContext struct {
	Logger
	File             *os.File
	writer           *bufio.Writer
	FormatingFunc    func(value interface{}) string
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
		if fwc.FormatingFunc == nil {
			buffer.WriteString(fmt.Sprint(value))
		} else {
			buffer.WriteString(fwc.FormatingFunc(value))
		}
	}
	buffer.WriteString(fwc.CallDelimiter)
	return buffer.String()
}

// Close the underlying file.
func (fwc *FileWritingContext) Close() error {
	if fwc.writer == nil || fwc.File == nil {
		return fmt.Errorf("trying to close already close log file %s", fwc.Path)
	}
	fwc.writer.Flush()
	err := fwc.File.Close()
	fwc.writer = nil
	fwc.File = nil
	return err
}

// Init initialises the output file.
func (fwc *FileWritingContext) Init() error {
	file, err := os.OpenFile(fwc.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	fwc.File = file
	fwc.writer = bufio.NewWriter(file)
	fwc.Logger = func(level string, values ...interface{}) {
		byts := fwc.FormatValues(level, values...)
		_, err := fwc.writer.WriteString(byts)
		if err != nil {
			panic(err)
		}
		err = fwc.writer.Flush()
		if err != nil {
			panic(err)
		}
	}
	// adds a lock for keeping logs consistent.
	fwc.Logger = fwc.Logger.WithLock()
	return nil
}

// DirLogger is a struct for keeping that of files open in the same folder.
type DirLogger struct {
	DirContext FileWritingContext
	OpenFiles  []FileWritingContext
}

// topicToPath convert a topic to a file path.
func (dirLogger DirLogger) topicToPath(topic string) string {
	return fmt.Sprintf("%s/%s", dirLogger.DirContext.Path, topic)
}

// find returns the open stream for the file.
// returns nil on file not open.
func (dirLogger DirLogger) find(topic string) *FileWritingContext {
	for _, openFile := range dirLogger.OpenFiles {
		if openFile.Path == dirLogger.topicToPath(topic) {
			return &openFile
		}
	}
	return nil
}

// Get return a filewritingcontext generating filename from the topic.
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

// Close the Directory files.
// Can panic.
func (dirLogger *DirLogger) Close() []error {
	errors := make([]error, 0)
	for _, openFile := range dirLogger.OpenFiles {
		err := openFile.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}
	dirLogger.OpenFiles = nil
	return errors
}

// GetLoggerFactory opens a directory for writing logs by topic.
func (dirLogger *DirLogger) GetLoggerFactory() LoggerFactory {
	lock := &sync.Mutex{}
	return func(topic string) Logger {
		lock.Lock()
		defer lock.Unlock()
		return dirLogger.Get(topic).Logger
	}
}

// NewDirLogger returns a logger that dispatch topic in a folder files.
func NewDirLogger(dirContext FileWritingContext) (*DirLogger, error) {
	err := os.Mkdir(dirContext.Path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	dirLogger := DirLogger{
		DirContext: dirContext,
		OpenFiles:  make([]FileWritingContext, 0),
	}
	return &dirLogger, nil
}
