package parser

import (
	"fmt"
	"puter/ast"
	"strconv"
	"testing"
)

func TestNumberParsing(t *testing.T) {
	exp := NewParser("1").Parse()
	conv, _ := strconv.ParseFloat(exp.String(), 64)
	if conv != 1 {
		t.Fatalf("Parsing result is not 1, got %f", conv)
	}
}

func TestOperatorExpression(t *testing.T) {
	exp := NewParser("1 + 2").Parse()
	result := exp.String()
	expected := fmt.Sprintf("(%f + %f)", 1.0, 2.0)
	if result != expected {
		t.Fatalf("Parsing result is not %s, got %s", expected, result)
	}
	expression, ok := exp.(*ast.OperatorExpression)
	if !ok {
		t.Fatalf("Operator not binary operator, got: %+v", expression)
	}
}

func TestCallExpressionWithArguments(t *testing.T) {
	exp := NewParser("add(1, 2, 3)").Parse()
	result := exp.String()
	if result != fmt.Sprintf("add(%f, %f, %f)", 1.0, 2.0, 3.0) {
		t.Fatalf("Parsing result is not add(1), got %s", result)
	}
	expression, ok := exp.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Operator not call expression, got: %+v", expression)
	}

	if len(expression.Args) != 3 {
		t.Fatalf("Expression arguments not 3, got %d", len(expression.Args))
	}

	identExpression, ok2 := expression.FunctionNameExpression.(*ast.IdentExpression)
	if !ok2 {
		t.Fatalf("Expected function name expression to be IdentExpression, instead got: %+v", identExpression)
	}
}

func TestCallExpressionWithNoArguments(t *testing.T) {
	exp := NewParser("add()").Parse()
	result := exp.String()
	if result != "add()" {
		t.Fatalf("Parsing result is not add(), got %s", result)
	}

	expression, ok := exp.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Operator call expression, got: %+v", expression)
	}

	if len(expression.Args) != 0 {
		t.Fatalf("Expression arguments not 0, got %d", len(expression.Args))
	}

	identExpression, ok2 := expression.FunctionNameExpression.(*ast.IdentExpression)
	if !ok2 {
		t.Fatalf("Expected function name expression to be IdentExpression, instead got: %+v", identExpression)
	}
}

func TestNumberAssignExpression(t *testing.T) {
	exp := NewParser("a = 2").Parse()
	result := exp.String()
	expected := fmt.Sprintf("a = %f", 2.0)
	if result != expected {
		t.Fatalf("Parsing result is not '%s', got '%s'", expected, result)
	}

	expression, ok := exp.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("Expected assign expression, got: %+v", expression)
	}

	if expression.Name.Token().Type != ast.IDENT {
		t.Fatalf("Token value not IDENT, got: %+v", expression.Name.Token().Type)
	}

	rightExp, ok2 := expression.Right.(*ast.NumberExpression)
	if !ok2 {
		t.Fatalf("Expect right expression to be Number expression, instead got: %+v", rightExp)
	}
}

// func TestFunctionAssignExpression(t *testing.T) {
// 	exp := NewParser("a = 2").Parse()
// 	result := exp.String()
// 	if result != "add()" {
// 		t.Fatalf("Parsing result is not add(), got %s", result)
// 	}
// }

// func TestNameAssignExpression(t *testing.T) {

// }
