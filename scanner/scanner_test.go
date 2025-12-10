package scanner

import (
	"puter/ast"
	"testing"
)

func TestScan(t *testing.T) {
	scanner := NewScanner("=!02938#*Hello")

	expectations := []string{
		"=",
		"!",
		"02938",
		"#",
		"*",
		"Hello",

		// EOF below
		"",
		"",
		"",
		"",
	}

	for _, e := range expectations {
		r := scanner.Next()
		if r.Literal != e {
			t.Fatalf("Expected =, got %s", r.Literal)
		}
	}
}

func TestAssignmentScan(t *testing.T) {
	scanner := NewScanner("a = 2")

	expectations := []string{
		"a",
		"=",
		"2",
		"",
	}

	for _, e := range expectations {
		r := scanner.Next()
		if r.Literal != e {
			t.Fatalf("Expected =, got %s", r.Literal)
		}
	}
}

func TestNumberScan(t *testing.T) {
	scanner := NewScanner("1")
	result := scanner.Next()
	if result.Type != ast.NUMBER && result.Literal != "1" {
		t.Fatalf("Expected 1, got %s", result.Literal)
	}
}
