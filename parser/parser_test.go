package parser

import (
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
