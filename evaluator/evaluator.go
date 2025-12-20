package evaluator

import (
	"fmt"
	ast "puter/ast"
	b "puter/box"
	p "puter/parser"
)

type ValueConverter = func(fromValue float64, fromUnit string, toUnit string) (float64, error)

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
	case *ast.OperatorExpression:
		switch exp.Operator.Type {
		case ast.PLUS:
			return e.evalPlusExpression(exp.Left, exp.Right)
		case ast.MINUS:
			return nil
		case ast.IN:
			return e.evalInExpression(exp.Left, exp.Right)
		case ast.ASTERISK:
			return nil
		case ast.LOGICAL_AND:
			return nil
		case ast.LOGICAL_OR:
			return nil
		case ast.LT:
			return nil
		case ast.GT:
			return nil
		case ast.CARET:
			return nil
		default:
			panic("Invalid operator token")
		}
	case *ast.BooleanExpression:
		return &b.BooleanBox{Value: exp.ActualValue}
	case *ast.IdentExpression:
		return &b.IdentBox{Value: exp.ActualValue}
	case *ast.NumberExpression:
		return &b.NumberBox{Value: exp.ActualValue, Tok: exp.Token()}
	default:
		x := exp.String()
		panic(fmt.Sprintf("Unhandled case %s", x))
	}

}

// If left is not a number box, but another unit-based expression, convert it first before returning a new value.
func (e *Evaluator) evalInExpression(leftExpr ast.Expression, rightExpr ast.Expression) b.Box {
	rightBox := e.evalExp(rightExpr)
	if rightBox.Type() != b.IDENTIFIER_BOX {
		panic("Right side of an in expression must be a unit identifier")
	}

	leftBox := func() b.Box {
		evaluated := e.evalExp(leftExpr)

		_, leftBoxIsIdent := evaluated.(*b.IdentBox)
		if !leftBoxIsIdent {
			return evaluated
		}
		got, ok := e.heap[evaluated.Inspect()]
		if !ok {
			panic(fmt.Sprintf("Identifier %s not set"))
		}
		return got
	}()

	// if left already a number box, no need for conversion. Just use the unit on the right
	// otherwise try to convert by converting whatever unit left is to the right unit.
	switch box := leftBox.(type) {
	case *b.NumberBox:
		return &b.CurrencyBox{
			Number: box,
			Unit:   rightBox.Inspect(),
		}
	case *b.CurrencyBox:
		rightUnit := rightBox.Inspect()
		if rightUnit == box.Unit {
			return &b.CurrencyBox{Number: box.Number, Unit: rightUnit}
		}

		converted, err := e.currencyConverter(box.Number.Value, box.Unit, rightUnit)
		if err != nil {
			panic(err)
		}
		return &b.CurrencyBox{Number: &b.NumberBox{Value: converted, Tok: box.Number.Tok}, Unit: rightUnit}

	default:
		panic("Invalid left hand side of an in expresison.")
	}

}

func (e *Evaluator) evalPlusExpression(left ast.Expression, right ast.Expression) {
	// if left and right is number, just add them, otherwise
	// if they are of different unit, convert right to left unit if possible then add them.
	// otherwise if either one has unit but the other one is just a number, then use that unit

}
