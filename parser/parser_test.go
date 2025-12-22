package parser

import (
	"fmt"
	"puter/ast"
	"strconv"
	"testing"
)

func TestNumberParsing(t *testing.T) {
	exp := NewParser().Parse("1")
	conv, _ := strconv.ParseFloat(exp.String(), 64)
	if conv != 1 {
		t.Fatalf("Parsing result is not 1, got %f", conv)
	}
}

func TestOperatorExpression(t *testing.T) {
	exp := NewParser().Parse("1 + 2")
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
	exp := NewParser().Parse("add(1, 2, 3)")
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
	exp := NewParser().Parse("add()")
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
	exp := NewParser().Parse("a = 2")
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

func TestNameAssignExpression(t *testing.T) {
	exp := NewParser().Parse("a = b")
	result := exp.String()
	expected := "a = b"
	if result != expected {
		t.Fatalf("Parsing result is not '%s', got '%s'", expected, result)
	}

	expression, ok := exp.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("Expected assign expression, got: %+v", expression)
	}

	nameExp, ok2 := expression.Name.(*ast.IdentExpression)
	if !ok2 {
		t.Fatalf("Expect name expression to be ident expression, instead got: %+v", nameExp)
	}

	rightExp, ok3 := expression.Right.(*ast.IdentExpression)
	if !ok3 {
		t.Fatalf("Expect right expression to be Number expression, instead got: %+v", rightExp)
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a",
			"(-a)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"4 > 4",
			"(4 > 4)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"2 != 5",
			"(2 != 5)",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"x in usd",
			"(x in usd)",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a && b",
			"(a && b)",
		},
		{
			"a && b || c",
			"((a && b) || c)",
		},
		{
			"add(a, sub(-b))",
			"add(a, sub((-b)))",
		},
		{
			"2 usd",
			"(2 usd)",
		},
		{
			"x = 2 usd in thb",
			"x = ((2 usd) in thb)",
		},
		{
			"3 in usd in thb in btc in xx",
			"((((3 in usd) in thb) in btc) in xx)",
		},
	}

	for _, tt := range tests {
		p := NewParser()
		expression := p.Parse(tt.input)

		actual := expression.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
