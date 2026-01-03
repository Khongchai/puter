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

func TestNumberBinaryOperatorEvaluations(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"2 + 4",
			"6",
			b.NUMBER_BOX,
		},
		{
			"10 - 4",
			"6",
			b.NUMBER_BOX,
		},
		{
			"3 * 5",
			"15",
			b.NUMBER_BOX,
		},
		{
			"10 / 2",
			"5",
			b.NUMBER_BOX,
		},
		{
			"5 / 2",
			"2.5",
			b.NUMBER_BOX,
		},
		{
			"2 + 3 * 4",
			"14",
			b.NUMBER_BOX,
		},
		{
			"(2 + 3) * 4",
			"20",
			b.NUMBER_BOX,
		},
		{
			"100 - 50 + 25",
			"75",
			b.NUMBER_BOX,
		},
		{
			"1.5 + 2.25",
			"3.75",
			b.NUMBER_BOX,
		},
		{
			"2 ^ 4",
			"16",
			b.NUMBER_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(t.Context(), getDefaultCurrencyConverter(200))

		obj := eval.EvalLine(c.Line)

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", c.ExpectPrint, obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %+v", obj.Type())
		}
	}
}

func TestBinaryBooleanOperatorEvaluations(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"true",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"false",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"true && true",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"true && false",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"true || false",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"false || false",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"5 == 5",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"5 != 5",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"10 > 5",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"3 <= 2",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"(5 > 2) && (10 < 20)",
			"true",
			b.BOOLEAN_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(t.Context(), getDefaultCurrencyConverter(200))

		obj := eval.EvalLine(c.Line)

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", c.ExpectPrint, obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %+v", obj.Type())
		}
	}
}

func TestCurrencyEvaluation(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"x = 2",
			"2",
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
		eval := NewEvaluator(t.Context(), getDefaultCurrencyConverter(200))

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
	eval := NewEvaluator(t.Context(), getDefaultCurrencyConverter(100))

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

func TestPrefixEvaluation(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"!true",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"!false",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"-1",
			"-1",
			b.NUMBER_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(t.Context(), getDefaultCurrencyConverter(200))

		obj := eval.EvalLine(c.Line)

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", c.ExpectPrint, obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %+v", obj.Type())
		}
	}
}

func TestComparsionEvaluation(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"1 < 2",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"100 == 1000",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"1 == 1",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"1 != 1.0",
			"false",
			b.BOOLEAN_BOX,
		},
		{
			"1 <= 1",
			"true",
			b.BOOLEAN_BOX,
		},
		{
			"-1 <= 1",
			"true",
			b.BOOLEAN_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(t.Context(), getDefaultCurrencyConverter(200))

		obj := eval.EvalLine(c.Line)

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", c.ExpectPrint, obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %+v", obj.Type())
		}
	}
}
