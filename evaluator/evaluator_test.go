package evaluator

import (
	b "puter/box"
	"testing"
)

func getDefaultCurrencyConverter() ValueConverter {
	return func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
		return 200.0, nil
	}
}

type EvaluationCase struct {
	Line        string
	ExpectPrint string
	ExpectType  b.BoxType
}

func TestEvaluation(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"x = 2",
			"2",
			b.NUMBER_BOX,
		},
		{
			"2 + 4",
			"6",
			b.NUMBER_BOX,
		},
		{
			"2 in usd",
			"2 usd",
			b.CURRENCY_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(getDefaultCurrencyConverter())

		obj := eval.EvalLine(c.Line)

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", "2", obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %s", obj.Type())
		}
	}
}

// try this case too
// func TestEvalWithValueConversion(t *testing.T) {
// 	eval := &Evaluator{}
// 	// The text at this point is expected to have been
// 	// sanitized already. eg. it matches a comment pattern // |, # |, <!-- | -->
// 	result := eval.EvalLine(0, "a = 1 + 2 in usd")
// 	result2 := eval.EvalLine(1, "k = a in thb")
// 	result3 := eval.EvalLine(1, "x = k + 2")
// }
