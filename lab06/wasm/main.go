package main

import (
	"math/big"
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

	num, ok := new(big.Int).SetString(p[0].String(), 10)
	if !ok {
		return js.ValueOf("Invalid argument type")
	}

	result := isPrime(num)

	if result {
		return js.ValueOf("It's prime")
	} else {
		return js.ValueOf("It's not prime")
	}
}

func isPrime(n *big.Int) bool {
	if n.Cmp(big.NewInt(1)) <= 0 {
		return false
	}

	return n.ProbablyPrime(0)
}
