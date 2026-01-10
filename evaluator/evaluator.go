package evaluator

import (
	"context"
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
	ctx               context.Context
}

func NewEvaluator(ctx context.Context, currencyConverter ValueConverter) *Evaluator {
	return &Evaluator{
		ctx:               ctx,
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
	case *ast.CallExpression:
		return e.evalCallExpression(exp.FunctionNameExpression, exp.Args)
	case *ast.OperatorExpression:
		switch exp.Operator.Type {
		case ast.PLUS:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a + b
			})
		case ast.MINUS:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a - b
			})
		case ast.SLASH:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a / b
			})
		case ast.IN:
			return e.evalInExpression(exp.Left, exp.Right)
		case ast.ASTERISK:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return a * b
			})
		case ast.CARET:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, func(a, b float64) float64 {
				return math.Pow(a, b)
			})
		case ast.LOGICAL_AND:
			return e.evalBinaryBooleanLogicalExpression(exp.Left, exp.Right, func(a, b bool) bool {
				return a && b
			})
		case ast.LOGICAL_OR:
			return e.evalBinaryBooleanLogicalExpression(exp.Left, exp.Right, func(a, b bool) bool {
				return a || b
			})
		case ast.EQ, ast.NOT_EQ, ast.LT, ast.GT, ast.LTE, ast.GTE:
			return e.evalBinaryBooleanComparisonExpression(exp.Left, exp.Right, exp.Operator.Type)
		default:
			panic("Invalid operator token")
		}
	case *ast.BooleanExpression:
		return &b.BooleanBox{Value: exp.ActualValue}
	case *ast.PostfixExpression:
		operator := exp.TokenValue.Type
		switch operator {
		case ast.IDENT:
			evaluated := e.evalInExpression(exp.Left, &ast.IdentExpression{
				ActualValue: exp.TokenValue.Literal,
				TokenValue:  exp.TokenValue,
			})
			return evaluated
		case ast.PERCENT:
			left := e.evalExp(exp.Left)
			l, isLeftNumber := left.(*b.NumberBox)
			if !isLeftNumber {
				panic("The left side of percent must be a number")
			}
			return &b.PercentBox{Value: l.Value}
		}
		panic("Postfix not supported")
	case *ast.PrefixExpression:
		operator := exp.TokenValue.Type
		switch right := e.evalExp(exp.Right).(type) {
		case *b.NumberBox:
			if operator == ast.MINUS {
				return &b.NumberBox{Value: -right.Value}
			}
			if operator == ast.PLUS {
				return &b.NumberBox{Value: right.Value}
			}
			panic("Unsupported prefix operation on a number")
		case *b.CurrencyBox:
			if operator == ast.MINUS {
				return &b.CurrencyBox{Value: -right.Value, Unit: right.Unit}
			}
			if operator == ast.PLUS {
				return &b.CurrencyBox{Value: right.Value, Unit: right.Unit}
			}
			panic("Unsupported prefix operation on a number")
		case *b.BooleanBox:
			if operator == ast.BANG {
				return &b.BooleanBox{Value: !right.Value}
			}
			panic("Unsupported prefix operation on a boolean")
		default:
			panic("The right-hand side of this prefix expression is invalid")
		}
	case *ast.IdentExpression:
		found, ok := e.heap[exp.ActualValue]
		if !ok {
			panic(fmt.Sprintf("Identifier %s not found", found.Inspect()))
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
			Value: box.Value,
			Unit:  right.ActualValue,
		}
	case *b.CurrencyBox:
		rightUnit := right.ActualValue
		if rightUnit == box.Unit {
			return &b.CurrencyBox{Value: box.Value, Unit: rightUnit}
		}

		converted, err := e.currencyConverter(box.Value, box.Unit, rightUnit)
		if err != nil {
			panic(err)
		}
		return &b.CurrencyBox{Value: converted, Unit: rightUnit}

	default:
		panic("Invalid left hand side of an in expresison.")
	}
}

func (e *Evaluator) evalBinaryBooleanLogicalExpression(left ast.Expression, right ast.Expression, callable func(a, b bool) bool) b.Box {
	evalLeft, leftIsBool := e.evalExp(left).(*b.BooleanBox)
	evalRight, rightIsBool := e.evalExp(right).(*b.BooleanBox)
	if !leftIsBool || !rightIsBool {
		panic("Both left and right side of a boolean binary operation must be boolean")
	}

	return &b.BooleanBox{Value: callable(evalLeft.Value, evalRight.Value)}
}

func (e *Evaluator) evalBinaryBooleanComparisonExpression(left ast.Expression, right ast.Expression, operatorType ast.TokenType) b.Box {
	evaluatedRight := e.evalExp(right)
	evaluatedLeft := e.evalExp(left)
	if operatorType == ast.EQ {
		return &b.BooleanBox{Value: evaluatedRight.Inspect() == evaluatedLeft.Inspect()}
	}
	if operatorType == ast.NOT_EQ {
		return &b.BooleanBox{Value: evaluatedRight.Inspect() != evaluatedLeft.Inspect()}
	}

	// TODO let's support this later, for example 3 usd > 2 thb. But for now 3 usd > (2 thb in usd) is fine.
	if evaluatedRight.Type() != evaluatedLeft.Type() {
		panic("Comparsion of different type or unit")
	}

	comp := func(comp func(a, b float64) bool) *b.BooleanBox {
		switch l := (evaluatedLeft).(type) {
		case *b.NumberBox:
			r, _ := (evaluatedRight).(*b.NumberBox)
			return &b.BooleanBox{Value: comp(l.Value, r.Value)}
		case *b.CurrencyBox:
			r, _ := (evaluatedRight).(*b.CurrencyBox)
			return &b.BooleanBox{Value: comp(l.Value, r.Value)}
		}
		return nil
	}

	result := func() *b.BooleanBox {
		switch operatorType {
		case ast.GT:
			return comp(func(a, b float64) bool { return a > b })
		case ast.LT:
			return comp(func(a, b float64) bool { return a < b })
		case ast.GTE:
			return comp(func(a, b float64) bool { return a >= b })
		case ast.LTE:
			return comp(func(a, b float64) bool { return a <= b })
		}
		return nil
	}()

	if result != nil {
		return result
	}

	panic("Unsupported boolean comparsion expression")
}

// For now we only support builtin functions and all builtin functions are simple math functions.
func (e *Evaluator) evalCallExpression(functionName ast.Expression, arguments []ast.Expression) b.Box {
	readArgs := func(expectedCount int, unlimited bool) []float64 {
		if !unlimited && len(arguments) != expectedCount {
			panic(fmt.Sprintf("Expected %d arguments, got %d len(arguments)", expectedCount, len(arguments)))
		}

		var results []float64
		for _, arg := range arguments {
			evaluated, ok := e.evalExp(arg).(*b.NumberBox)
			if !ok {
				panic(fmt.Sprintf("Method expect number, but got %s instead", evaluated.Type()))
			}
			results = append(results, evaluated.Value)
		}

		return results
	}

	result := func() float64 {
		switch functionName.String() {
		case "log10":
			return math.Log10(readArgs(1, false)[0])
		case "logE":
			return math.Log(readArgs(1, false)[0])
		case "log2":
			return math.Log2(readArgs(1, false)[0])
		case "round":
			return math.Round(readArgs(1, false)[0])
		case "floor":
			return math.Floor(readArgs(1, false)[0])
		case "ceil":
			return math.Ceil(readArgs(1, false)[0])
		case "abs":
			return math.Abs(readArgs(1, false)[0])
		case "sin":
			return math.Sin(readArgs(1, false)[0])
		case "cos":
			return math.Cos(readArgs(1, false)[0])
		case "tan":
			return math.Tan(readArgs(1, false)[0])
		case "sqrt":
			return math.Sqrt(readArgs(1, false)[0])
		case "lerp":
			read := readArgs(3, false)
			return (1-read[2])*read[0] + read[2]*read[1]
		case "invLerp":
			read := readArgs(3, false)
			return (read[2] - read[0]) / (read[1] - read[0])
		case "sum":
			read := readArgs(-1, true)
			s := 0.0
			for _, n := range read {
				s += n
			}
			return s
		case "product":
			read := readArgs(-1, true)
			s := 0.0
			for _, n := range read {
				s *= n
			}
			return s
		default:
			panic("Unrecognized function name")
		}
	}()

	return &b.NumberBox{Value: result}
}

func (e *Evaluator) evalBinaryNumberExpression(left ast.Expression, right ast.Expression, callable func(a, b float64) float64) b.Box {
	evalLeft := e.evalExp(left)
	evalRight := e.evalExp(right)
	switch l := evalLeft.(type) {
	case *b.NumberBox:
		switch r := evalRight.(type) {
		case *b.NumberBox:
			return &b.NumberBox{Value: callable(l.Value, r.Value)}
		case *b.CurrencyBox:
			return &b.CurrencyBox{Value: callable(l.Value, r.Value), Unit: r.Unit}
		default:
			panic("Type not supported. Cannot add.")
		}
	case *b.CurrencyBox:
		switch r := evalRight.(type) {
		case *b.NumberBox:
			return &b.CurrencyBox{Value: callable(l.Value, r.Value), Unit: l.Unit}
		case *b.CurrencyBox:
			if r.Unit == l.Unit {
				return &b.CurrencyBox{Value: callable(l.Value, r.Value), Unit: l.Unit}
			}

			// convert left to right
			leftConverted, err := e.currencyConverter(l.Value, l.Unit, r.Unit)
			if err != nil {
				panic(err)
			}

			return &b.CurrencyBox{Value: callable(leftConverted, r.Value), Unit: r.Unit}
		default:
			panic("Type not supported. Cannot perform binary expression.")
		}

	default:
		panic("Type not supported. Cannot perform binary expression.")
	}

}
