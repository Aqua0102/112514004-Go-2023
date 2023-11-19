// wasm/main.go

package main

import (
	"fmt"
	"math/big"
	"strconv"
	"syscall/js"
)

func CheckPrime(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return js.ValueOf(false) // Invalid argument count
	}

	// Convert the argument to a Go integer
	num, err := strconv.Atoi(args[0].String())
	if err != nil {
		return js.ValueOf(false) // Invalid argument type
	}

	// Check if the number is prime
	result := isPrime(num)

	// Return the result to JavaScript
	return js.ValueOf(result)
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	// Use math/big package to check primality
	bigNum := big.NewInt(int64(n))
	return bigNum.ProbablyPrime(0)
}

func registerCallbacks() {
	// Register the CheckPrime function to be accessible from JavaScript
	js.Global().Set("CheckPrime", js.FuncOf(CheckPrime))
}

func main() {
	fmt.Println("Golang main function executed")
	registerCallbacks()

	// Need to block the main thread forever
	select {}
}
