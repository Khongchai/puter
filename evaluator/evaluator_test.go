package evaluator

import (
	"testing"
)

func TestNumberAssignment(t *testing.T) {
	eval := NewEvaluator()

	obj := eval.EvalLine("x = 2")

	if obj.Inspect() != "2" {
		t.Fatalf("Expected inspect result to be %s, got %s", "2", obj.Inspect())
	}
	if obj.Type() != NUMBER_OBJ {
		t.Fatalf("Expected identifier object, got %s", obj.Type())
	}
}

// func TestEvalWithValueConversion(t *testing.T) {
// 	eval := &Evaluator{}
// 	// The text at this point is expected to have been
// 	// sanitized already. eg. it matches a comment pattern // |, # |, <!-- | -->
// 	result := eval.EvalLine(0, "a = 1 + 2 in usd")
// 	result2 := eval.EvalLine(1, "k = a in thb")
// 	result3 := eval.EvalLine(1, "x = k + 2")
// }
