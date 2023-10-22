package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Calculator(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	resultStr := "Error!"
	if len(parts) != 4  {
		resultStr = "Error!"
		fmt.Fprintf(w, resultStr)
		return
	}

	operator := parts[1]

	num1, err1 := strconv.Atoi(parts[2])
	num2, err2 := strconv.Atoi(parts[3])

	if err1 != nil || err2 != nil {
		resultStr = "Error!"
		fmt.Fprintf(w, resultStr)
		return
	}

	result := 0
	symbol := ""
	reminder := ""
	switch operator {
	case "add":
		result = num1 + num2
		symbol = "+"
	case "sub":
		result = num1 - num2
		symbol = "-"
	case "mul":
		result = num1 * num2
		symbol = "*"
	case "div":
		if num2 == 0 {
			resultStr = "Error!"
			fmt.Fprintf(w, resultStr)
			return
		}
		result = num1 / num2
		reminder = ", reminder = " + strconv.Itoa(num1%num2)
		symbol = "/"
	default:
		resultStr = "Error!"
		fmt.Fprintf(w, resultStr)
		return
	}

	resultStr = fmt.Sprintf("%d %s %d = %d%s", num1, symbol, num2, result, reminder)
	fmt.Fprintf(w, resultStr)
}
func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
