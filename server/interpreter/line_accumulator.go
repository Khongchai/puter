package interpreter

import (
	"fmt"
	"puter/evaluation/evaluator"
	"puter/evaluation/evaluator/box"
)

var startingValues = map[string]float64{
	"sum": 0, "difference": 0, "product": 1, "quotient": 1,
}

var commandToOperation = map[string]func(a, b float64) float64{
	"sum": add, "product": multiply, "quotient": quotient, "difference": difference,
}

type LineAccumulator struct {
	line              int
	command           string
	acc               box.Box
	operation         func(a, b float64) float64
	currencyConverter evaluator.ValueConverter
}

func NewLineAccumulator(command string, line int, currencyConverter evaluator.ValueConverter) *LineAccumulator {
	if !IsAccumulationCommand(command) {
		panic(fmt.Sprintf("Invalid line command. Got %s", command))
	}
	operation := commandToOperation[command]
	got := &LineAccumulator{
		line,
		command,
		nil,
		operation,
		currencyConverter,
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

	switch acc := l.acc.(type) {
	case *box.NumberBox:
		{
			switch r := result.(type) {
			case *box.NumberBox:
				acc.Value = l.operation(r.Value, acc.Value)
			case *box.CurrencyBox:
				num := l.operation(r.Number.Value, acc.Value)
				l.acc = &box.CurrencyBox{Number: box.NewNumberbox(num, r.Number.NumberType), Unit: r.Unit}
			}
		}
	case *box.CurrencyBox:
		{
			switch r := result.(type) {
			case *box.NumberBox:
				acc.Number.Value = l.operation(r.Value, acc.Number.Value)
			case *box.CurrencyBox:
				// convert acc unit to r unit
				if acc.Unit != r.Unit {
					converted, err := l.currencyConverter(acc.Number.Value, acc.Unit, r.Unit)
					if err != nil {
						return
					}
					acc.Unit = r.Unit
					acc.Number.Value = converted
				}
				acc.Number.Value = l.operation(r.Number.Value, acc.Number.Value)
			}
		}
	}
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
