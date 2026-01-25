package interpreter

import (
	"puter/evaluation/evaluator"
	"testing"
)

func getDefaultCurrencyConverter(defaultValue float64) evaluator.ValueConverter {
	return func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
		return defaultValue, nil
	}
}

func joinLines(lines ...string) string {
	collected := ""
	for _, line := range lines {
		collected += line
		collected += "\n"
	}
	return collected
}

func TestInterpretEmptyFile(t *testing.T) {
	interpreter := NewInterpreter(t.Context(), getDefaultCurrencyConverter(200))
	result := interpreter.Interpret("")
	if len(result) != 0 {
		t.Fatalf("Expected empty result, got instead %d results", len(result))
	}
}

func TestInterpretingValidSingleLineResult(t *testing.T) {
	interpreter := NewInterpreter(t.Context(), getDefaultCurrencyConverter(200))
	validCases := []string{
		"//|1+2",
		"// | 1+2",
		"#|1+2",
		"# | 1+2",
		joinLines(
			"/*",
			"* | 1 + 2",
			"*/",
		),
		joinLines(
			"/*",
			"| 1 + 2",
			"*/",
		),
		// not handling this, kind of useless pattern, just do // bro
		// "/* | 1 + 2*/",
	}

	for _, validCase := range validCases {
		interpretations := interpreter.Interpret(validCase)
		if len(interpretations) != 1 {
			t.Fatalf("Expected case %s interpretations length to be 1, got %d", validCase, len(interpretations))
		}
		if interpretations[0].Decoration != "3" {
			t.Fatalf("Decoration of %s is not 3, got %s", validCase, interpretations[0].Decoration)
		}
		if len(interpretations[0].Diagnostics) != 0 {
			t.Fatalf("Diagnostics length of %s should be 0, got %d", validCase, len(interpretations[0].Diagnostics))
		}
	}
}

// test that we won't ever be stuck in an infinite loop...
// nor do we encounter nil error
func FuzzInterpretation(f *testing.F) {
	f.Add("//| 1+1")
	f.Add("/* | 2+2 */")
	f.Add("# | 3+3")
	f.Add("/ / | 1")
	f.Add("/*\n| 1\n*/")
	f.Add("\r\n#|1\r\n")
	f.Add("🔥🔥🔥 #| 1+1")

	f.Fuzz(func(t *testing.T, input string) {
		interpreter := NewInterpreter(t.Context(), getDefaultCurrencyConverter(200))

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Recovered from panic with input %q: %v", input, r)
			}
		}()

		interpretations := interpreter.Interpret(input)

		for _, interp := range interpretations {
			if interp == nil {
				t.Error("Returned a nil interpretation")
			}
			if interp.LineIndex < 0 {
				t.Errorf("Negative line index: %d", interp.LineIndex)
			}
		}
	})
}
