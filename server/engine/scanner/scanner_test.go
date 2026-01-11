package scanner

import (
	"puter/engine/ast"
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
	scanner := NewScanner("foo = 2")

	expectations := []string{
		"foo",
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
	cases := []string{
		"1",
		"100.111",
	}

	for _, c := range cases {
		scanner := NewScanner(c)
		result := scanner.Next()
		if result.Type != ast.NUMBER || result.Literal != c {
			t.Fatalf("Expected %s, got %s", c, result.Literal)
		}
	}

}

func TestKeywords(t *testing.T) {
	scanner := NewScanner("oops in true lol false 4")
	expectations := []ast.TokenType{
		ast.IDENT,
		ast.IN,
		ast.TRUE,
		ast.IDENT,
		ast.FALSE,
		ast.NUMBER,
		ast.EOF,
	}

	for _, e := range expectations {
		r := scanner.Next()
		if r.Type != e {
			t.Fatalf("Expected %s, got %s", e, r.Type)
		}
	}
}

func TestNameFollowedByNumber(t *testing.T) {
	scanner := NewScanner("log10")
	result := scanner.Next()
	if result.Type != ast.IDENT {
		t.Fatalf("Expected %s, got %s", ast.IDENT, result.Type)
	}
	if result.Literal != "log10" {
		t.Fatalf("Expected log10, got %s", result.Literal)
	}
}

func TestLogicalOperators(t *testing.T) {
	scanner := NewScanner("&& &! ||")
	expectations := []ast.TokenType{
		ast.LOGICAL_AND,
		ast.ILLEGAL,
		ast.BANG,
		ast.LOGICAL_OR,
		ast.EOF,
	}

	for _, e := range expectations {
		r := scanner.Next()
		if r.Type != e {
			t.Fatalf("Expected %s, got %s", e, r.Type)
		}
	}
}
