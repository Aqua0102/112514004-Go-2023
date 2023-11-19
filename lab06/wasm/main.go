package main

import (
	"math/big"
	"strconv"
	"syscall/js"
)

func main() {
	js.Global().Set("CheckPrime", js.FuncOf(CheckPrime))

	select {}
}

func CheckPrime(this js.Value, p []js.Value) interface{} {
	if len(p) != 1 {
		return js.ValueOf("Invalid argument count")
	}

	// Convert the argument to a Go integer
	num, err := strconv.Atoi(p[0].String())
	if err != nil {
		return js.ValueOf("Invalid argument type")
	}

	// Check if the number is prime
	result := isPrime(num)

	// Return the result to JavaScript
	if result {
		return js.ValueOf("It's prime")
	} else {
		return js.ValueOf("It's not prime")
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	// Use math/big package to check primality
	bigNum := big.NewInt(int64(n))
	return bigNum.ProbablyPrime(0)
}
