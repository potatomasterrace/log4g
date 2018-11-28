# Log4g
Logging abstraction for golang.
Current Supported outputs : File, Directory or RamBuffer or Console.
# Install 
    go get github.com/potatomasterrace/log4g
# QuickStart
```Go
package main

import (
	"fmt"
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
	for f := 2; f < squareRoot; f++ {
		if isFactor(n, f, logger) {
			logger(DEBUG, "is not prime factor", f, "found")
			return false
		}
	}
	logger(INFO, "confirmed prime")
	return false
}

func main() {
	// Getting the logger
	logger, buffer := NewInMemoryLogger()
	// Filtering a level
	logger = logger.Filter(DEBUG)
	// passing it to function
	isPrime(41, logger)
	isPrime(103, logger)
	// Outputing the buffer
	for _, line := range *buffer {
		fmt.Println(line)
	}
	// optionnal : overriding the buffer to make it available for GC
	*buffer = nil
}
```

## Output 
```
    [[INFO]   -> isPrime [41] :  square root 6]
    [[TRACE]  -> isPrime [41] :   -> isFactor [41 2] :  false]
    [[TRACE]  -> isPrime [41] :   -> isFactor [41 3] :  false]
    [[TRACE]  -> isPrime [41] :   -> isFactor [41 4] :  false]
    [[TRACE]  -> isPrime [41] :   -> isFactor [41 5] :  false]
    [[INFO]   -> isPrime [41] :  confirmed prime]
    [[INFO]   -> isPrime [103] :  square root 10]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 2] :  false]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 3] :  false]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 4] :  false]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 5] :  false]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 6] :  false]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 7] :  false]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 8] :  false]
    [[TRACE]  -> isPrime [103] :   -> isFactor [103 9] :  false]
    [[INFO]   -> isPrime [103] :  confirmed prime]
```
# 
