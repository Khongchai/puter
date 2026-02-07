package evaluator

import (
	"context"
	"fmt"
	"log"
	"math"
	"puter/evaluation/ast"
	b "puter/evaluation/evaluator/box"
	p "puter/evaluation/parser"
)

type Evaluator struct {
	parser p.Parser
	// A map of identifier to puter object
	heap              map[string]b.Box
	currencyConverter b.ValueConverter
	ctx               context.Context
	// Eval stage holds multiple diagnostics error.
	// For simplicity, parser and tokenizer always return one errors but this stage returns
	// multiple since it's the most complex.
	//
	// For example, it might be useful to let the user know that x and y in sum(1) + sum(x, y, z) are not numbers. Instead
	// of just x or y.
	//
	// With the parser stage, if the syntax is broken -- it's broken. Just emit whatever error first encountered and return.
	//
	// There is also the fact that evaluations are line-by-line here so error here does not mean the entire program halts.
	diagnostics []*ast.Diagnostic
}

func NewEvaluator(ctx context.Context, currencyConverter b.ValueConverter) *Evaluator {
	return &Evaluator{
		ctx:               ctx,
		parser:            *p.NewParser(),
		heap:              make(map[string]b.Box),
		currencyConverter: currencyConverter,
	}
}

// Evaluate the content of a line. Line separation is assumed
// to have been done by some earlier stage.
//
// The returned b.Box is nullable if an error is encountered during evaluation
func (e *Evaluator) EvalLine(text string) b.Box {
	e.diagnostics = []*ast.Diagnostic{}
	expression, err := e.parser.Parse(text)
	if err != nil {
		e.diagnostics = append(e.diagnostics, err)
		return nil
	}
	result := e.evalExp(expression)
	return result
}

func (e *Evaluator) evalExp(expression ast.Expression) b.Box {
	switch exp := expression.(type) {
	case *ast.AssignExpression:
		value := e.evalExp(exp.Right)
		ident, ok := exp.Name.(*ast.IdentExpression)
		if ok {
			e.heap[ident.ActualValue] = value
			return value

		}
		e.diagnostics = append(
			e.diagnostics,
			ast.NewDiagnosticAtToken("Expected an identifier", ident.Token()),
		)
		return value
	case *ast.CallExpression:
		return e.evalCallExpression(exp.FunctionNameExpression, exp.Args)
	case *ast.OperatorExpression:
		switch exp.Operator.Type {
		case ast.PLUS:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, exp.Operator, func(a, b float64) float64 {
				return a + b
			})
		case ast.MINUS:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, exp.Operator, func(a, b float64) float64 {
				return a - b
			})
		case ast.SLASH:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, exp.Operator, func(a, b float64) float64 {
				return a / b
			})
		case ast.IN:
			return e.evalInExpression(exp.Left, exp.Right)
		case ast.ASTERISK:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, exp.Operator, func(a, b float64) float64 {
				return a * b
			})
		case ast.CARET:
			return e.evalBinaryNumberExpression(exp.Left, exp.Right, exp.Operator, func(a, b float64) float64 {
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
			return e.evalBinaryBooleanComparisonExpression(exp.Left, exp.Right, exp.Operator)
		default:
			return nil
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
				e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
					"The left-hand side expression of a percent symbol must be a number type",
					exp.Left.Token(),
				))
				return nil
			}
			return &b.PercentBox{Value: l.Value}
		}
		e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
			"Unrecognized postfix operator",
			exp.Token(),
		))
		return nil
	case *ast.PrefixExpression:
		operator := exp.TokenValue.Type
		switch right := e.evalExp(exp.Right).(type) {
		case *b.NumberBox:
			if operator == ast.MINUS {
				return b.NewNumberbox(-right.Value, right.NumberType)
			}
			if operator == ast.PLUS {
				return b.NewNumberbox(right.Value, right.NumberType)
			}
			e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
				"Unsupported prefix operation on a number",
				exp.Token(),
			))
			return nil
		case *b.CurrencyBox:
			if operator == ast.MINUS {
				return &b.CurrencyBox{Number: b.NewNumberbox(-right.Number.Value, right.Number.NumberType), Unit: right.Unit}
			}
			if operator == ast.PLUS {
				return &b.CurrencyBox{Number: b.NewNumberbox(right.Number.Value, right.Number.NumberType), Unit: right.Unit}
			}
			e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
				"Unsupported prefix operation on a currency",
				exp.Token(),
			))
			return nil
		case *b.BooleanBox:
			if operator == ast.BANG {
				return &b.BooleanBox{Value: !right.Value}
			}
			e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
				"Unsupported prefix operation on a boolean",
				exp.Token(),
			))
			return nil
		default:
			e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
				"The right-hand expression of this prefix operation is not supported.",
				exp.Token(),
			))
			return nil
		}
	case *ast.IdentExpression:
		found, ok := e.heap[exp.ActualValue]
		if !ok {
			e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
				fmt.Sprintf("Identifier %s not found", exp.ActualValue),
				exp.Token(),
			))
			return nil
		}
		return found
	case *ast.NumberExpression:
		return b.NewNumberbox(exp.ActualValue, b.Decimal)
	default:
		x := exp.String()
		log.Fatalf("Evaluator error: unhandled case %s", x)
		return nil
	}
}

func (e *Evaluator) evalInExpression(leftExpr ast.Expression, rightExpr ast.Expression) b.Box {
	right, rightIsIdentifier := rightExpr.(*ast.IdentExpression)
	if !rightIsIdentifier {
		e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken("Right side of an in expression must be a unit identifier", right.Token()))
		return nil
	}

	leftBox := e.evalExp(leftExpr)
	if operatable, ok := leftBox.(b.InPrefixOperatable); !ok {
		e.diagnostics = append(e.diagnostics, ast.NewDiagnostic(
			"Left hand side of this expression is not evaluable by this operator",
			leftExpr.Token().StartPos(),
			rightExpr.Token().EndPos(),
		))
		return nil
	} else {
		res, err := operatable.OperateIn(right.ActualValue, e.currencyConverter)
		if err != nil {
			e.diagnostics = append(e.diagnostics, ast.NewDiagnostic(
				err.Error(),
				leftExpr.Token().StartPos(),
				right.Token().EndPos(),
			))
		}
		return res
	}
}

func (e *Evaluator) evalBinaryBooleanLogicalExpression(left ast.Expression, right ast.Expression, callable func(a, b bool) bool) b.Box {
	evalLeft, leftIsBool := e.evalExp(left).(*b.BooleanBox)
	evalRight, rightIsBool := e.evalExp(right).(*b.BooleanBox)
	if !leftIsBool || !rightIsBool {
		e.diagnostics = append(e.diagnostics, ast.NewDiagnostic(
			"Both must be boolean",
			left.Token().StartPos(),
			right.Token().EndPos(),
		))
		return nil
	}

	return &b.BooleanBox{Value: callable(evalLeft.Value, evalRight.Value)}
}

func (e *Evaluator) evalBinaryBooleanComparisonExpression(left ast.Expression, right ast.Expression, operator *ast.Token) b.Box {

	evaluatedRight := e.evalExp(right)
	evaluatedLeft := e.evalExp(left)
	if operator.Type == ast.EQ {
		return &b.BooleanBox{Value: evaluatedRight.Inspect() == evaluatedLeft.Inspect()}
	}
	if operator.Type == ast.NOT_EQ {
		return &b.BooleanBox{Value: evaluatedRight.Inspect() != evaluatedLeft.Inspect()}
	}

	if evaluatedRight.Type() != evaluatedLeft.Type() {
		e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(
			fmt.Sprintf("Can't compare %s and %s", evaluatedLeft.Type(), evaluatedRight.Type()),
			operator,
		))
		return nil
	}

	comp := func(comp func(a, b float64) bool) *b.BooleanBox {
		switch l := (evaluatedLeft).(type) {
		case *b.NumberBox:
			r, _ := (evaluatedRight).(*b.NumberBox)
			return &b.BooleanBox{Value: comp(l.Value, r.Value)}
		case *b.CurrencyBox:
			r, _ := (evaluatedRight).(*b.CurrencyBox)
			if l.Unit == r.Unit {
				return &b.BooleanBox{Value: comp(l.Number.Value, r.Number.Value)}
			}

			converted, err := e.currencyConverter(l.Number.Value, l.Unit, r.Unit)
			if err != nil {
				e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken(err.Error(), operator))
				return nil
			}
			return &b.BooleanBox{Value: comp(converted, r.Number.Value)}
		}
		e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken("Relational operators only applicable to currency or number types", operator))
		return nil
	}
	switch operator.Type {
	case ast.GT:
		return comp(func(a, b float64) bool { return a > b })
	case ast.LT:
		return comp(func(a, b float64) bool { return a < b })
	case ast.GTE:
		return comp(func(a, b float64) bool { return a >= b })
	case ast.LTE:
		return comp(func(a, b float64) bool { return a <= b })
	default:
		e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken("Unrecognized operator", operator))
		return nil
	}
}

// For now we only support builtin functions and all builtin functions are simple math functions.
func (e *Evaluator) evalCallExpression(functionName ast.Expression, arguments []ast.Expression) b.Box {
	var args []float64
	var diagnostics []*ast.Diagnostic
	for _, arg := range arguments {
		evaluated, ok := e.evalExp(arg).(*b.NumberBox)
		if !ok {
			diagnostics = append(diagnostics, ast.NewDiagnosticAtToken(
				fmt.Sprintf("Method expect number, but got %s instead", evaluated.Type()),
				arg.Token(),
			))
		} else {
			args = append(args, evaluated.Value)
		}
	}

	if len(diagnostics) > 0 {
		return nil
	}

	matchArgsAndReturn := func(expectedCount int, fn func([]float64) float64) b.Box {
		if len(arguments) != expectedCount {
			text := fmt.Sprintf("Expected %d arguments, got %d len(arguments)", expectedCount, len(arguments))
			e.diagnostics = append(
				e.diagnostics,
				ast.NewDiagnosticAtToken(text, functionName.Token()),
			)
			return nil
		}
		v := fn(args)
		return &b.NumberBox{Value: v}
	}
	switch functionName.String() {
	case "log10":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Log10(v[0]) })
	case "logE":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Log(v[0]) })
	case "log2":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Log2(v[0]) })
	case "round":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Round(v[0]) })
	case "floor":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Floor(v[0]) })
	case "ceil":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Ceil(v[0]) })
	case "abs":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Abs(v[0]) })
	case "sin":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Sin(v[0]) })
	case "cos":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Cos(v[0]) })
	case "tan":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Tan(v[0]) })
	case "sqrt":
		return matchArgsAndReturn(1, func(v []float64) float64 { return math.Sqrt(v[0]) })
	case "lerp":
		return matchArgsAndReturn(3, func(v []float64) float64 {
			return (1-args[2])*args[0] + args[2]*args[1]
		})
	case "invLerp":
		return matchArgsAndReturn(3, func(v []float64) float64 {
			return (args[2] - args[0]) / (args[1] - args[0])
		})
	case "sum":
		s := 0.0
		for _, n := range args {
			s += n
		}
		return &b.NumberBox{Value: s}
	case "product":
		s := 0.0
		for _, n := range args {
			s *= n
		}
		return &b.NumberBox{Value: s}
	default:
		e.diagnostics = append(e.diagnostics, ast.NewDiagnosticAtToken("Unknown function name", functionName.Token()))
		return nil
	}
}

func (e *Evaluator) evalBinaryNumberExpression(left ast.Expression, right ast.Expression, operator *ast.Token, operation func(a, b float64) float64) b.Box {
	var boxLeft b.Box = e.evalExp(left)
	var boxRight b.Box = e.evalExp(right)
	if operatable, ok := boxLeft.(b.BinaryNumberOperatable); !ok {
		e.diagnostics = append(e.diagnostics, ast.NewDiagnostic(
			"Left hand side of this expression is not evaluable by this operator",
			left.Token().StartPos(),
			operator.EndPos(),
		))
		return nil
	} else {
		res, err := operatable.OperateBinary(boxRight, operation, e.currencyConverter)
		if err != nil {
			e.diagnostics = append(e.diagnostics, ast.NewDiagnostic(
				err.Error(),
				left.Token().StartPos(),
				right.Token().EndPos(),
			))
		}
		return res
	}
}

func (e *Evaluator) GetDiagnostics() []*ast.Diagnostic {
	return e.diagnostics
}
