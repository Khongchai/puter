package scanner

import (
	"iter"
	"testing"
)

func TestScan(t *testing.T) {
	next, stop := iter.Pull(Scan("=!02938#*Hello", 0))
	defer stop()

	expectations := []string{
		"=",
		"!",
		"02938",
		"#",
		"*",
		"Hello",
	}

	for _, e := range expectations {
		r, ok := next()
		if !ok {
			t.Fatalf("Scan iterator ends too early")
		}
		if r.Literal != e {
			t.Fatalf("Expected =, got %s", r.Literal)
		}
	}
}
