package interpreter

import (
	"fmt"
	"puter/evaluation/evaluator/box"
)

type LineAccumulator struct {
	line      int
	command   string
	acc       float64
	operation func(a, b float64) float64
}

func NewLineAccumulator(command string, line int) (*LineAccumulator, error) {
	if !IsAccumulationCommand(command) {
		return nil, fmt.Errorf("Invalid line command. Got %s", command)
	}
	var acc float64
	var operation func(a, b float64) float64
	switch command {
	case "sum":
		operation = add
		acc = 0
	case "product":
		operation = multiply
		acc = 1
	case "difference":
		operation = difference
		acc = 0
	case "quotient":
		operation = quotient
		acc = 1
	}
	return &LineAccumulator{
		line,
		command,
		acc,
		operation,
	}, nil
}

func (l *LineAccumulator) Print() string {
	return fmt.Sprintf("%g", l.acc)
}

func (l *LineAccumulator) GetLine() int {
	return l.line
}

func (l *LineAccumulator) Accept(result box.Box) {
	resultNumber, resultIsNumber := result.(*box.NumberBox)
	if resultIsNumber {
		l.acc = l.operation(l.acc, resultNumber.Value)
		return
	}

	// resultCurrency, resultIsCurrency := result.(*box.CurrencyBox)
	// if resultIsCurrency {

	// }
}

func IsAccumulationCommand(text string) bool {
	return text == "sum" || text == "product" || text == "quotient" || text == "difference"
}

func multiply(a, b float64) float64 {
	return a * b
}

func add(a, b float64) float64 {
	return a + b
}

func difference(a, b float64) float64 {
	return a - b
}

func quotient(a, b float64) float64 {
	return a / b
}
