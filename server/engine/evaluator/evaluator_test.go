package evaluator

import (
	"errors"
	"fmt"
	"math"
	b "puter/engine/evaluator/box"
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
			"(2 + 3) * 4",
			"20",
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
			"2 usd",
			"2 usd",
			b.CURRENCY_BOX,
		},
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
		{
			"1 in usd + 2 in thb",
			"202 thb",
			b.CURRENCY_BOX,
		},
		{
			"1 usd + 2 in thb",
			"202 thb",
			b.CURRENCY_BOX,
		},
		{
			"1 usd + 2 thb",
			"202 thb",
			b.CURRENCY_BOX,
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
		{
			"1 usd > 2 thb",
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

func TestPercent(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"2%",
			"2%",
			b.PERCENT_BOX,
		},
		{
			"2 + 5%",
			fmt.Sprintf("%g", 2*1.05),
			b.NUMBER_BOX,
		},
		{
			"2 in usd + 10% + 5 in usd",
			fmt.Sprintf("%g usd", 2*1.1+5),
			b.CURRENCY_BOX,
		},
		{
			"2 usd + 10% + 5 thb",
			fmt.Sprintf("%g thb", (2*1.1)*34.4+5),
			b.CURRENCY_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(t.Context(), func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
			if fromUnit == "usd" && toUnit == "thb" {
				return (fromValue * 34.4), nil
			}
			return -1, errors.New("Currency conversion not supported in this test suite")
		})

		obj := eval.EvalLine(c.Line)

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", c.ExpectPrint, obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %+v", obj.Type())
		}
	}
}

func TestBuiltinFunctionEvaluations(t *testing.T) {
	cases := []*EvaluationCase{
		{
			"logE(2)",
			fmt.Sprintf("%g", math.Log(2)),
			b.NUMBER_BOX,
		},
		{
			"log10(2)",
			fmt.Sprintf("%g", math.Log10(2)),
			b.NUMBER_BOX,
		},
		{
			"log10(2 + 5)",
			fmt.Sprintf("%g", math.Log10(2+5)),
			b.NUMBER_BOX,
		},
		{
			"log2(2)",
			fmt.Sprintf("%g", math.Log2(2)),
			b.NUMBER_BOX,
		},
		{
			"sqrt(10)",
			fmt.Sprintf("%g", math.Sqrt(10)),
			b.NUMBER_BOX,
		},
		{
			"lerp(0, 10, 0.5)",
			"5",
			b.NUMBER_BOX,
		},
		{
			"invLerp(10, 20, 15)",
			"0.5",
			b.NUMBER_BOX,
		},
	}
	for _, c := range cases {
		eval := NewEvaluator(t.Context(), getDefaultCurrencyConverter(200))

		obj := eval.EvalLine(c.Line)

		if len(eval.diagnostics) != 0 {
			v := ""
			for _, d := range eval.diagnostics {
				v += fmt.Sprintf("[%s] ", d.Message)
			}
			t.Fatalf("Eval errors for case %s: %s", c.Line, v)
		}

		if obj.Inspect() != c.ExpectPrint {
			t.Fatalf("Expected inspect result to be %s, got %s", c.ExpectPrint, obj.Inspect())
		}
		if obj.Type() != c.ExpectType {
			t.Fatalf("Expected identifier object, got %+v", obj.Type())
		}
	}
}
