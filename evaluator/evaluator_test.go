package evaluator

import (
	b "puter/box"
	"testing"
)

func getDefaultCurrencyConverter(defaultValue float64) ValueConverter {
	return func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
		return defaultValue, nil
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
		{
			"2 in usd in thb",
			"200 thb",
			b.CURRENCY_BOX,
		},
		{
			"x = 2 in usd in thb",
			"200 thb",
			b.CURRENCY_BOX,
		},
		{
			"x = 1 + 2 in usd",
			"3 usd",
			b.CURRENCY_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(getDefaultCurrencyConverter(200))

		obj := eval.EvalLine(c.Line)

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", "2", obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %+v", obj.Type())
		}
	}
}

func TestCurrencyConversionMultiline(t *testing.T) {
	eval := NewEvaluator(getDefaultCurrencyConverter(100))

	eval.EvalLine("a = 1 + 2 in usd") // 3 usd
	eval.EvalLine("k = a in thb")     // 100 thb
	eval.EvalLine("x = k + 2")        // 102 thb
	result := eval.EvalLine("x")      // 102 thb
	if result.Inspect() != "102 thb" {
		t.Fatalf("Expected 102 thb, got %s", result.Inspect())
	}
	if result.Type() != b.CURRENCY_BOX {
		t.Fatalf("Expected currency, got %+v", result.Type())
	}
}
