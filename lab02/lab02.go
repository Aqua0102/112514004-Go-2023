package main

import (
	"fmt"
)

func main() {
	var n int64

	fmt.Print("Enter a number: ")
	fmt.Scanln(&n)

	result := Sum(n)
	fmt.Println(result)
}

func Sum(n int64) string {
	var sum int64
	sum = 0
	separator := ""

	for i := int64(1); i <= n; i++ {
		if i%7 != 0 {
			sum += i
			fmt.Print(separator, i)
			separator = "+"
		}
	}

	return fmt.Sprintf("=%d", sum)
}
