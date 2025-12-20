package evaluator

import (
	"fmt"
	ast "puter/ast"
	b "puter/box"
	"puter/lib"
	p "puter/parser"
)

type ValueConverter = func(fromValue float64, toValue float64, fromUnit string, toUnit string) (*lib.Promise[float64], bool)

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
		e.heap[identifier.TokenValue().Literal] = value
		return value
	case *ast.OperatorExpression:
		switch exp.Operator.Type {
		case ast.PLUS:
			return nil
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
	// todo do we need to promisify these guys below?
	case *ast.BooleanExpression:
		return &b.BooleanBox{Value: lib.NewResolvedPromise(exp.ActualValue)}
	case *ast.IdentExpression:
		return &b.IdentBox{Value: exp.ActualValue}
	case *ast.NumberExpression:
		return &b.NumberBox{Value: exp.ActualValue}
	default:
		x := exp.String()
		panic(fmt.Sprintf("Unhandled case %s", x))
	}

}

func (e *Evaluator) evalInExpression(left ast.Expression, right ast.Expression) b.Box {
	leftBox := e.evalExp(left)
	rightBox := e.evalExp(right)
	leftBoxCasted, leftIsNumberBox := leftBox.(*b.NumberBox)
	if !leftIsNumberBox {
		panic("Left side of an in expression must be a number")
	}
	if rightBox.Type() != b.IDENTIFIER_BOX {
		panic("Right side of an in expression must be a unit identifier")
	}

	unitIdentifier := rightBox.Inspect().Await()

	// TODO check against other possible units. For now just currency.
	if _, unitIsCurrency := b.ValidCurrencies[unitIdentifier]; !unitIsCurrency {
		panic(fmt.Sprintf("%s is not a valid ISO 4217 currency code.", unitIdentifier))
	}

	return &b.CurrencyBox{
		Number: leftBoxCasted,
		Unit:   unitIdentifier,
	}
}
