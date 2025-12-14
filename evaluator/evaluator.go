package evaluator

import (
	"fmt"
	ast "puter/ast"
	p "puter/parser"
)

type Evaluator struct {
	parser p.Parser
	// A map of identifier to puter object
	heap map[string]Box
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		parser: *p.NewParser(),
		heap:   make(map[string]Box),
	}
}

// Evaluate the content of a line. Line separation is assumed
// to have been done by some earlier stage.
func (e *Evaluator) EvalLine(text string) Box {
	expression := e.parser.Parse(text)
	result := e.evalExp(expression)
	return result
}

func (e *Evaluator) evalExp(expression ast.Expression) Box {
	switch exp := expression.(type) {
	case *ast.AssignExpression:
		identifier := e.evalExp(exp.Name)
		value := e.evalExp(exp.Right)
		e.heap[identifier.Inspect()] = value
		return value
	case *ast.BooleanExpression:
		return &BooleanBox{exp.ActualValue}
	case *ast.IdentExpression:
		return &IdentBox{exp.ActualValue}
	case *ast.NumberExpression:
		return &NumberBox{exp.ActualValue}
	default:
		x := exp.String()
		panic(fmt.Sprintf("Unhandled case %s", x))
	}

}
