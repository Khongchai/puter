package interpreter

import (
	"puter/evaluation/evaluator"
	"strconv"
	"strings"
	"testing"
	"time"
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
		if interpretations[0].EvalResult != "3" {
			t.Fatalf("Decoration of %s is not 3, got %s", validCase, interpretations[0].EvalResult)
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

func TestIntegration(t *testing.T) {
	type TestCase struct {
		ExpectPrint []string
		ExpectLine  []int
		InputText   string
	}
	cases := []*TestCase{
		{
			ExpectPrint: []string{"3", "5"},
			ExpectLine:  []int{5, 6},
			InputText: joinLines(
				"func (e *Engine) handleTextDocumentDidChange(ctx context.Context, params *lsproto.DidChangeTextDocumentParams) error {",
				"// uri := params.TextDocument.Uri",
				"return nil",
				"}",
				"",
				"// | a = 1 + 2",
				"// | a + 2",
			),
		},
	}

	for _, testCase := range cases {
		interpreter := NewInterpreter(t.Context(), getDefaultCurrencyConverter(200))
		interpretations := interpreter.Interpret(testCase.InputText)
		if len(testCase.ExpectPrint) != len(testCase.ExpectLine) {
			t.Fatalf("Invalid test case")
		}
		for i := range interpretations {
			if testCase.ExpectPrint[i] != interpretations[i].EvalResult {
				t.Fatalf("Expected %s, instead got %s", testCase.ExpectPrint[i], interpretations[i].EvalResult)
			}
			if testCase.ExpectLine[i] != interpretations[i].LineIndex {
				t.Fatalf("Expected line of result %s to be %d, not %d", interpretations[i].EvalResult, testCase.ExpectLine[i], interpretations[i].LineIndex)
			}
		}
	}

}

func TestPerformanceLargeFile(t *testing.T) {
	interpreter := NewInterpreter(t.Context(), getDefaultCurrencyConverter(200))

	// 1. generate a "huge" source file (10,000 lines)
	lineCount := 10000
	lines := make([]string, lineCount)
	for i := 0; i < lineCount; i++ {
		// every 10th line is an evaluation trigger, others are "JS"
		if i%10 == 0 {
			lines[i] = "// | 1 + " + strconv.Itoa(i)
		} else {
			lines[i] = "const x" + strconv.Itoa(i) + " = () => { console.log('hello'); };"
		}
	}
	hugeText := strings.Join(lines, "\n")

	start := time.Now()
	results := interpreter.Interpret(hugeText)
	elapsed := time.Since(start)

	t.Logf("\n--- Performance Result ---")
	t.Logf("File Size:    %.2f MB", float64(len(hugeText))/(1024*1024))
	t.Logf("Total Lines:  %d", lineCount)
	t.Logf("Evaluations:  %d", len(results))
	t.Logf("Time Taken:   %f", elapsed.Seconds())
	t.Logf("Avg per line: %f", (elapsed / time.Duration(lineCount)).Seconds())
	t.Logf("--------------------------")
}
