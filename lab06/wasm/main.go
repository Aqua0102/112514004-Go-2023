package main

import (
	"fmt"
	"math/big"
	"syscall/js"
)

func registerCallBacks() {
	js.Global().Set("CheckPrime", js.FuncOf(CheckPrime))
}

func CheckPrime(this js.Value, p []js.Value) interface{} {
	js.Global().Get("answer").Set("innerText", "test")

	numStr := js.Global().Get("document").Call("getElementById", "value").Get("value").String()
	num, ok := new(big.Int).SetString(numStr, 10)
	if !ok {
		updateAnswer("Invalid argument type")
		return js.ValueOf("Invalid argument type")
	}
	fmt.Printf(numStr)
	result := isPrime(num)

	if result {
		updateAnswer("It's prime")
	} else {
		updateAnswer("It's not prime")
	}

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

func updateAnswer(answer string) {
	js.Global().Get("document").Call("getElementById", "answer").Set("innerText", answer)
}

func main() {
	registerCallBacks()

	select {}
}
