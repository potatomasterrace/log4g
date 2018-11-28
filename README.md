# Log4g
Logging abstraction for golang.
Current Supported outputs : File, Directory or ram or Console.
# Install 
    go get github.com/potatomasterrace/log4g
# QuickStart
## Code
```Go
package main

import (
	"math"

	. "github.com/potatomasterrace/log4g"
)

func isFactor(n int, f int, logger LoggerStream) bool {
	logger = logger.FunctionCall(n, f)
	isfactor := n%f == 0
	logger(TRACE, isfactor)
	return isfactor
}
func isPrime(n int, logger LoggerStream) bool {
	// Declaring a function call
	logger = logger.FunctionCall(n)
	squareRoot := int(math.Sqrt(float64(n)))
	// Logging stuff
	logger(INFO, "square root", squareRoot)
	for f := 3; f < squareRoot; f+=2 {
		if isFactor(n, f, logger) {
			logger(DEBUG, "is not prime factor", f, "found")
			return false
		}
	}
	logger(INFO, "is prime")
	return false
}

func main() {
	// Getting the logger
	logger := NewConsoleLogger()

	// Filtering a level
	logger = logger.Filter(DEBUG)

	// prepending time to logger
	logger = logger.PrependTime()

	// // Possible to append text
	// logger= logger.Append()

	// // Possible to prepend text
	// logger= logger.Prepend()

	// passing it to function
	isPrime(41, logger)
	isPrime(103, logger)
}

```
## Console output 
```
[INFO]  : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [41] :  square root 6]
[TRACE] : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [41] :   -> isFactor [41 3] :  false]
[TRACE] : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [41] :   -> isFactor [41 5] :  false]
[INFO]  : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [41] :  is prime]
[INFO]  : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [103] :  square root 10]
[TRACE] : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [103] :   -> isFactor [103 3] :  false]
[TRACE] : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [103] :   -> isFactor [103 5] :  false]
[TRACE] : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [103] :   -> isFactor [103 7] :  false]
[TRACE] : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [103] :   -> isFactor [103 9] :  false]
[INFO]  : [Thu, 29 Nov 2018 00:38:11 CET  -> isPrime [103] :  is prime]
```
# Using file for logging 
``` Go
	// Getting the logger
	fwc := FileWritingContext{
        // FilePath
		Path:             "./logs"
        // Function called to convert a value to string
		FormatingFunc: func(v interface{}) string {
			return fmt.Sprintf("%s", v)
        },
        // Value to separate between function calls
		CallDelimiter:    "\r\n",
        // Value to separate between function args
        ValuesDelimiters: " ",}
    // Initializes the LoggerStream
	err := fwc.Init()
	defer fwc.Close()
	if err != nil {
		panic(err)
    }
    // logger can be used like the quickstart
	logger := fwc.LoggerStream
```
## Output 
same data as the quickstart written in file ./logs :
```
[INFO]  Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [41] :  square root %!s(int=6)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [41] :   -> isFactor [41 2] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [41] :   -> isFactor [41 3] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [41] :   -> isFactor [41 4] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [41] :   -> isFactor [41 5] :  %!s(bool=false)
[INFO]  Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [41] :  is prime
[INFO]  Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :  square root %!s(int=10)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 2] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 3] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 4] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 5] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 6] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 7] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 8] :  %!s(bool=false)
[TRACE] Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :   -> isFactor [103 9] :  %!s(bool=false)
[INFO]  Thu, 29 Nov 2018 00:42:16 CET  -> isPrime [103] :  is prime 
```
# Using a directory for logging
## Code 
```Golang
	// Create the LoggerFactory
	folderpath := "./logs"
	dirLogger,err := NewDirLogger(FileWritingContext{
		Path:             folderpath,
		FormatingFunc:    func(v interface{}) string { return fmt.Sprintf("%s", v) },
		CallDelimiter:    "\r\n",
		ValuesDelimiters: " ",
	})
    loggerFactory := dirLogger.GetLoggerFactory()
    defer loggerFactory.Close()
    // this logger writes to ./logs/file1
	logger1 := loggerFactory("file1")
    // this logger writes to ./logs/file2
	logger2 := loggerFactory("file2")
```
# Using ram for logging
```Golang
	logger,buffer:= NewInMemoryLogger()
	// buffer is type *[][]interface{} and contains all the logged data.
```

# Intercept panic called inside logger
Call method No Panic of the logger
```Golang 
	level := WARN
	msgs := []interface{}{
		"hello","world"
	}
	err := logger.NoPanic(level,msgs...)
```