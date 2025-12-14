package evaluator

import (
	"fmt"
	ast "puter/ast"
	b "puter/box"
	p "puter/parser"
)

type ValueConverter = func(fromValue float64, toValue float64, fromUnit string, toUnit string) (float64, bool)

type Evaluator struct {
	parser p.Parser
	// A map of identifier to puter object
	heap              map[string]b.Box
	currencyConverter ValueConverter
}

func NewEvaluator(currencyConverter ValueConverter) *Evaluator {
	return &Evaluator{
		parser:            *p.NewParser(),
		heap:              make(map[string]b.Box),
		currencyConverter: currencyConverter,
	}
}

// Evaluate the content of a line. Line separation is assumed
// to have been done by some earlier stage.
func (e *Evaluator) EvalLine(text string) b.Box {
	expression := e.parser.Parse(text)
	result := e.evalExp(expression)
	return result
}

func (e *Evaluator) evalExp(expression ast.Expression) b.Box {
	switch exp := expression.(type) {
	case *ast.AssignExpression:
		identifier := e.evalExp(exp.Name)
		value := e.evalExp(exp.Right)
		e.heap[identifier.Inspect()] = value
		return value
	case *ast.BooleanExpression:
		return &b.BooleanBox{Value: exp.ActualValue}
	case *ast.IdentExpression:
		return &b.IdentBox{Value: exp.ActualValue}
	case *ast.NumberExpression:
		return &b.NumberBox{Value: exp.ActualValue}
	default:
		x := exp.String()
		panic(fmt.Sprintf("Unhandled case %s", x))
	}

}
