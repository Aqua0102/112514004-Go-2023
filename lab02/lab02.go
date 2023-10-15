package main

import (
	"fmt"
	"strconv"
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
	answer := ""
	separator := ""

	for i := int64(1); i <= n; i++ {
		if i%7 != 0 {
			sum += i
			answer += separator + strconv.Itoa(int(i))
			separator = "+"
		}
	}

	return answer + fmt.Sprintf("=%d", sum)
}
