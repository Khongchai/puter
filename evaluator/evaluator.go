package evaluator

import (
	"fmt"
	"math"
	"puter/ast"
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
		value := e.evalExp(exp.Right)
		ident, ok := exp.Name.(*ast.IdentExpression)
		if !ok {
			panic("Invalid identifier")
		}
		e.heap[ident.ActualValue] = value
		return value
	case *ast.OperatorExpression:
		switch exp.Operator.Type {
		case ast.PLUS:
			return e.evalBinaryArithmeticNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a + b
			})
		case ast.MINUS:
			return e.evalBinaryArithmeticNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a - b
			})
		case ast.SLASH:
			return e.evalBinaryArithmeticNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a / b
			})
		case ast.IN:
			return e.evalInExpression(exp.Left, exp.Right)
		case ast.ASTERISK:
			return e.evalBinaryArithmeticNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a * b
			})
		case ast.GT:
			return nil
		case ast.CARET:
			return e.evalBinaryArithmeticNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return math.Pow(a, b)
			})
		case ast.LT:
			return nil
		case ast.LOGICAL_AND:
			return e.evalBinaryBooleanExpression(exp.Left, exp.Right, func(a, b bool) bool {
				return a && b
			})
		case ast.LOGICAL_OR:
			return e.evalBinaryBooleanExpression(exp.Left, exp.Right, func(a, b bool) bool {
				return a || b
			})
		default:
			panic("Invalid operator token")
		}
	case *ast.BooleanExpression:
		return &b.BooleanBox{Value: exp.ActualValue}
	case *ast.IdentExpression:
		found, ok := e.heap[exp.ActualValue]
		if !ok {
			panic("Identifier not found")
		}
		return found
	case *ast.NumberExpression:
		return &b.NumberBox{Value: exp.ActualValue}
	default:
		x := exp.String()
		panic(fmt.Sprintf("Unhandled case %s", x))
	}

}

// If left is not a number box, but another unit-based expression, convert it first before returning a new value.
func (e *Evaluator) evalInExpression(leftExpr ast.Expression, rightExpr ast.Expression) b.Box {
	right, rightIsIdentifier := rightExpr.(*ast.IdentExpression)
	if !rightIsIdentifier {
		panic("Right side of an in expression must be a unit identifier")
	}

	// if left already a number box, no need for conversion. Just use the unit on the right
	// otherwise try to convert by converting whatever unit left is to the right unit.
	leftBox := e.evalExp(leftExpr)
	switch box := leftBox.(type) {
	case *b.NumberBox:
		return &b.CurrencyBox{
			Number: box,
			Unit:   right.ActualValue,
		}
	case *b.CurrencyBox:
		rightUnit := right.ActualValue
		if rightUnit == box.Unit {
			return &b.CurrencyBox{Number: box.Number, Unit: rightUnit}
		}

		converted, err := e.currencyConverter(box.Number.Value, box.Unit, rightUnit)
		if err != nil {
			panic(err)
		}
		return &b.CurrencyBox{Number: &b.NumberBox{Value: converted}, Unit: rightUnit}

	default:
		panic("Invalid left hand side of an in expresison.")
	}
}

func (e *Evaluator) evalBinaryBooleanExpression(left ast.Expression, right ast.Expression, callable func(a, b bool) bool) b.Box {
	evalLeft, leftIsBool := e.evalExp(left).(*b.BooleanBox)
	evalRight, rightIsBool := e.evalExp(right).(*b.BooleanBox)
	if !leftIsBool || !rightIsBool {
		panic("Both left and right side of a boolean binary operation must be boolean")
	}
	return &b.BooleanBox{Value: callable(evalLeft.Value, evalRight.Value)}
}

func (e *Evaluator) evalBinaryArithmeticNumberExpression(left ast.Expression, right ast.Expression, callable func(a, b float64) float64) b.Box {
	evalLeft := e.evalExp(left)
	evalRight := e.evalExp(right)
	switch l := evalLeft.(type) {
	case *b.NumberBox:
		switch r := evalRight.(type) {
		case *b.NumberBox:
			return &b.NumberBox{Value: callable(l.Value, r.Value)}
		case *b.CurrencyBox:
			return &b.CurrencyBox{Number: &b.NumberBox{Value: callable(l.Value, r.Number.Value)}, Unit: r.Unit}
		default:
			panic("Type not supported. Cannot add.")
		}
	case *b.CurrencyBox:
		switch r := evalRight.(type) {
		case *b.NumberBox:
			return &b.CurrencyBox{Number: &b.NumberBox{Value: callable(l.Number.Value, r.Value)}, Unit: l.Unit}
		case *b.CurrencyBox:
			if r.Unit == l.Unit {
				return &b.CurrencyBox{Number: &b.NumberBox{Value: callable(l.Number.Value, r.Number.Value)}, Unit: l.Unit}
			}

			panic("Can't add numbers of different unit")
		default:
			panic("Type not supported. Cannot add.")
		}

	default:
		panic("Type not supported. Cannot add.")
	}

}
