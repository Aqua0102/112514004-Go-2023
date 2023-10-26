package main

import (
	"fmt"
	"html/template"
	"log"

	"net/http"
	"strconv"
)

type PageData struct {
	Operation  string
	Num1       int
	Num2       int
	Result     int
	Expression string
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
	http.ServeFile(w, r, "error.html")
}

func calculate(op string, num1, num2 int) (int, string) {
	symbol := ""
	switch op {
	case "add":
		symbol = " + "
		return num1 + num2, symbol
	case "sub":
		symbol = " - "
		return num1 - num2, symbol
	case "mul":
		symbol = " * "
		return num1 * num2, symbol
	case "div":
		if num2 != 0 {
			symbol = " / "
			return num1 / num2, symbol
		}
	case "gcd":
		symbol = fmt.Sprintf("GCD(%d, %d)", num1, num2)
		return gcd(num1, num2), symbol
	case "lcm":
		symbol = fmt.Sprintf("LCM(%d, %d)", num1, num2)
		return lcm(num1, num2), symbol

	}
	return 0, "Error!"
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func Calculator(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	op := r.FormValue("op")
	num1, err1 := strconv.Atoi(r.FormValue("num1"))
	num2, err2 := strconv.Atoi(r.FormValue("num2"))

	if err1 != nil || err2 != nil {
		// 参数格式错误，可以执行相应的操作
		http.ServeFile(w, r, "error.html")
		return
	}

	result, symbol := calculate(op, num1, num2)

	expression := ""
	if op == "gcd" || op == "lcm" {
		expression = symbol
	} else {
		expression = fmt.Sprintf("%d%s%d", num1, symbol, num2)
	}

	data := PageData{
		Operation:  op,
		Num1:       num1,
		Num2:       num2,
		Result:     result,
		Expression: expression,
	}

	if symbol == "Error!" {
		http.ServeFile(w, r, "error.html")
	} else {
		err := template.Must(template.ParseFiles("index.html")).Execute(w, data)

		if err != nil {
			http.ServeFile(w, r, "error.html")
		}
	}

}

func main() {
	http.HandleFunc("/lab04", Calculator)
	http.HandleFunc("/index.html", ServeHTML)
	fmt.Println("Server is running on :8084...")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
