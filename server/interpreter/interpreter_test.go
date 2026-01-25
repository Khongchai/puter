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
	}

	for _, validCase := range validCases {
		interpretations := interpreter.Interpret(validCase)
		if len(interpretations) != 1 {
			t.Fatalf("Expected case %s interpretations length to be 1, got %d", validCase, len(interpretations))
		}
	}
}

// test that we won't ever be stuck in an infinite loop...
// nor do we encounter nil error
func FuzzInterpretation(f *testing.F) {

}
