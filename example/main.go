package main

import (
	"math"

	. "github.com/potatomasterrace/log4g"
)

func isFactor(n int, f int, logger Logger) bool {
	logger = logger.FunCall(n, f)
	isfactor := n%f == 0
	logger(TRACE, isfactor)
	return isfactor
}
func isPrime(n int, logger Logger) bool {
	// Declaring a function call
	logger = logger.FunCall(n)
	squareRoot := int(math.Sqrt(float64(n)))
	// Logging stuff
	logger(INFO, "square root", squareRoot)
	for f := 3; f < squareRoot; f += 2 {
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
