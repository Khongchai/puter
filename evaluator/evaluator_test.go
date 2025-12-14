package evaluator

import (
	b "puter/box"
	"testing"
)

func defaultCurrencyConverter(fromValue float64, toValue float64, fromUnit string, toUnit string) (float64, bool) {
	return 200, true
}

func TestNumberAssignment(t *testing.T) {
	eval := NewEvaluator(defaultCurrencyConverter)

	obj := eval.EvalLine("x = 2")

	if obj.Inspect() != "2" {
		t.Fatalf("Expected inspect result to be %s, got %s", "2", obj.Inspect())
	}
	if obj.Type() != b.NUMBER_BOX {
		t.Fatalf("Expected identifier object, got %s", obj.Type())
	}
}

func TestCurrencyConversion(t *testing.T) {
	eval := NewEvaluator(defaultCurrencyConverter)

	eval.EvalLine("x = 2 in usd")
	obj2 := eval.EvalLine("a = x in thb")

	if obj2.Inspect() != "200" {
		t.Fatalf("Expected inspect result to be %s, got %s", "200", obj2.Inspect())
	}
	if obj2.Type() != b.NUMBER_BOX {
		t.Fatalf("Expected identifier object, got %s", obj2.Type())
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
