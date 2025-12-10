package scanner

import (
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
