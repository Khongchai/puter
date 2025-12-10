package parser

import (
	"fmt"
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

func TestSimpleNumberExpressionParsing(t *testing.T) {
	exp := NewParser("1 + 2").Parse()
	result := exp.String()
	expected := fmt.Sprintf("(%f + %f)", 1.0, 2.0)
	if result != expected {
		t.Fatalf("Parsing result is not %s, got %s", expected, result)
	}
}
