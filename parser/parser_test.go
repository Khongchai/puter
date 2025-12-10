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
	exp := NewParser("add(1)").Parse()
	result := exp.String()
	if result != fmt.Sprintf("add(%f)", 1.0) {
		t.Fatalf("Parsing result is not add(1), got %s", result)
	}
	expression, ok := exp.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Operator not call expression, got: %+v", expression)
	}

	if len(expression.Args) != 1 {
		t.Fatalf("Expression arguments not 1, got %d", len(expression.Args))
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
