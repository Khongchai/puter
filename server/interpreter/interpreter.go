package interpreter

import (
	"context"
	"puter/evaluation/evaluator"
	lsproto "puter/lsp"
	"puter/utils"
	"slices"
	"strings"
)

type Interpreter struct {
	ctx               context.Context
	currencyConverter evaluator.ValueConverter
}

type Interpretation struct {
	LineIndex   int
	EvalResult  string
	Diagnostics []*lsproto.Diagnostic
}

// Interpreter takes in a text file, finds out if there is a line in that text file
// that starts with `//|` or `#|` (space between ignored), then start an evaluator for
// that line
// Example
//
// ```js
// // | a = log(2) < this starts evaluator 1
// // | b = a
//
// const something = 2;
//
// // | 2 + 3 < this starts evaluator 2
//
// for (let i = 0; i < 3; i++) {
//
// }
// ```
func NewInterpreter(ctx context.Context, currencyConverter evaluator.ValueConverter) *Interpreter {
	return &Interpreter{
		ctx,
		currencyConverter,
	}
}

// We do not yet need to care about the uri since we're doing full parsing
func (interpreter *Interpreter) Interpret(text string) []*Interpretation {
	evaluator := evaluator.NewEvaluator(interpreter.ctx, interpreter.currencyConverter)

	interpretations := []*Interpretation{}

	pos := 0
	i := 0 // line index
	lines := slices.Collect(strings.SplitSeq(text, "\n"))
	for i < len(lines) {
		if len(lines[i]) < 2 { // 2 is double slash, use this + "|" as minimum line length
			pos += len(lines[i])
			i++
			continue
		}

		firstTwoCharsNotSpace := func() string {
			s := ""
			for _, c := range lines[i] {
				if len(s) == 2 {
					return s
				}
				if c != ' ' {
					s += string(c)
				}
			}
			return s
		}()

		// python or c-like normal comment
		if firstTwoCharsNotSpace[0] == '#' || firstTwoCharsNotSpace == "//" {
			index := strings.Index(lines[i], "|")
			if index != -1 && len(lines[i]) > index+1 {
				evaluatable := lines[i][index+1:]
				interpretation := interpreter.evaluateAndInterpretResult(evaluator, evaluatable, i)
				interpretations = append(interpretations, interpretation)
			}
		}

		pos += len(text)
		i++
	}

	return interpretations
}

func (interpreter *Interpreter) evaluateAndInterpretResult(
	evaluator *evaluator.Evaluator,
	collected string,
	lineIndex int,
) *Interpretation {
	box := evaluator.EvalLine(collected)
	evalDiag := evaluator.GetDiagnostics()
	lsDiag := []*lsproto.Diagnostic{}
	if len(evalDiag) > 0 {
		for _, e := range evalDiag {
			lsDiag = append(lsDiag, &lsproto.Diagnostic{
				Severity: utils.PointerTo(lsproto.DiagnosticSeverityError),
				Range: lsproto.Range{
					Start: lsproto.Position{
						Line:      uint32(lineIndex),
						Character: uint32(e.StartPos),
					},
					End: lsproto.Position{
						Line:      uint32(lineIndex),
						Character: uint32(e.EndPos),
					},
				},
				Message: e.Message,
			})
		}
	}

	decoration := ""
	if box != nil {
		decoration = box.Inspect()
	}
	return &Interpretation{
		LineIndex:   lineIndex,
		Diagnostics: lsDiag,
		EvalResult:  decoration,
	}
}
