package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type opFuncType func(int, int) (int, error)

func add(i, j int) (int, error) { return i + j, nil }
func sub(i, j int) (int, error) { return i - j, nil }
func mul(i, j int) (int, error) { return i * j, nil }
func div(i, j int) (int, error) {
	if j == 0 {
		return 0, errors.New("division by 0")
	}
	return i / j, nil
}

var opMap = map[string]opFuncType{
	"+": add,
	"-": sub,
	"*": mul,
	"/": div,
}

func evalExpr(expr []string) (int, error) {
	if len(expr) != 3 {
		return 0, fmt.Errorf("invalid expression: %v", expr)

	}
	p1, err := strconv.Atoi(expr[0])
	if err != nil {
		return 0, fmt.Errorf("%q: %w", expr[0], err)

	}
	op := expr[1]
	opFunc, ok := opMap[op]
	if !ok {
		return 0, fmt.Errorf("Unsupported operator: %q", op)

	}
	p2, err := strconv.Atoi(expr[2])
	if err != nil {
		return 0, fmt.Errorf("%q: %w", expr[2], err)

	}
	return opFunc(p1, p2)
}

func main() {
	expressions := [][]string{
		{"2", "+", "3"},
		{"2", "-", "3"},
		{"2", "*", "3"},
		{"2", "/", "3"},
		{"2", "%", "3"},
		{"two", "+", "three"},
		{"5"},
	}
	for _, expression := range expressions {
		res, err := evalExpr(expression)
		if err != nil {
			fmt.Printf("%-12s -> error: %v\n", strings.Join(expression, " "), err)
			continue
		}
		fmt.Printf("%-12s -> %d\n", strings.Join(expression, " "), res)
	}
}
