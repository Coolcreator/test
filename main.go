package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/calculate", calculateHandler)
	http.ListenAndServe(":8080", nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Access") != "superuser" {
		fmt.Println("Доступ запрещен.")
		http.Error(w, "Access denied", http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	expr := strings.TrimSpace(string(body))
	if expr == "" {
		http.Error(w, "Missing expression in request body", http.StatusBadRequest)
		return
	}
	result := evaluateExpression(expr)
	response := map[string]int{"result": result}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func evaluateExpression(s string) int {
	s = strings.Replace(s, " ", "", -1)
	sign := 1
	stack := []int{}
	cursor := 0
	result, operand := 0, 0
	for cursor < len(s) {
		if s[cursor] == '-' {
			result = result + (sign * operand)
			operand = 0
			sign = -1
		} else if s[cursor] == '+' {
			result = result + (sign * operand)
			operand = 0
			sign = 1
		} else if s[cursor] == '(' {
			stack = append(stack, result, sign)
			result, sign = 0, 1
			operand = 0
		} else if s[cursor] == ')' {
			result = result + (sign * operand)
			oldSign := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result = result * oldSign
			oldOperand := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result = result + oldOperand
			operand = 0
		} else {
			operand = 10*operand + int(s[cursor]-'0')
			fmt.Println(operand)
		}
		cursor++
	}
	result = result + (sign * operand)
	return result
}
