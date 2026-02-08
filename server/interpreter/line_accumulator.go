package interpreter

import (
	"fmt"
	"puter/evaluation/evaluator/box"
	"puter/unit"
)

var startingValues = map[string]float64{
	"sum": 0, "difference": 0, "product": 1, "quotient": 1,
}

var commandToOperation = map[string]func(a, b float64) float64{
	"sum": add, "product": multiply, "quotient": quotient, "difference": difference,
}

type LineAccumulator struct {
	line       int
	command    string
	acc        box.Box
	operation  func(a, b float64) float64
	converters *unit.Converters
}

func NewLineAccumulator(command string, line int, converters *unit.Converters) *LineAccumulator {
	if !IsAccumulationCommand(command) {
		panic(fmt.Sprintf("Invalid line command. Got %s", command))
	}
	operation := commandToOperation[command]
	got := &LineAccumulator{
		line,
		command,
		nil,
		operation,
		converters,
	}
	return got
}

func (l *LineAccumulator) Print() string {
	return l.acc.Inspect()
}

func (l *LineAccumulator) GetLine() int {
	return l.line
}

// If result is not a valid type, this method does nothing.
func (l *LineAccumulator) Accept(result box.Box) {
	if l.acc == nil {
		l.setStartingAcc(result)
	}

	operatable, ok := l.acc.(box.BinaryNumberOperatable)
	if !ok {
		return
	}
	newAcc, err := operatable.OperateBinaryNumber(result, l.operation, l.converters)
	if err != nil {
		return
	}
	l.acc = newAcc
}

func (l *LineAccumulator) setStartingAcc(result box.Box) {
	num := startingValues[l.command]
	switch v := result.(type) {
	case *box.NumberBox:
		l.acc = box.NewNumberbox(num, v.NumberType)
	case *box.CurrencyBox:
		l.acc = &box.CurrencyBox{Number: box.NewNumberbox(num, v.Number.NumberType), Unit: v.Unit}
	}
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
