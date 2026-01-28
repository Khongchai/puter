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

		if firstTwoCharsNotSpace == "/*" {
			hasEndSameLine := strings.Contains(lines[i], "*/")
			if hasEndSameLine {
				pipeIndex := strings.Index(lines[i], "|")
				enderIndex := strings.Index(lines[i], "*/")
				if pipeIndex != -1 {
					middleText := lines[i][pipeIndex+1 : enderIndex]
					interpretation := interpreter.evaluateAndInterpretResult(evaluator, middleText, i)
					interpretations = append(interpretations, interpretation)
				}
			} else {
				pos += 2
				i++
				for i < len(lines) {
					hasEnd := strings.Contains(lines[i], "*/")
					if !hasEnd {
						index := strings.Index(lines[i], "|")
						if index != -1 && len(lines[i]) > index+1 {
							evaluatable := lines[i][index+1:]
							interpretation := interpreter.evaluateAndInterpretResult(evaluator, evaluatable, i)
							interpretations = append(interpretations, interpretation)
						}
					}
					i++
				}
			}
		}

		pos += len(text)
		i++
	}

	return interpretations
}

func peek(text string, i int) rune {
	if i >= len(text) {
		return -1
	}
	return rune(text[i])
}

func skipWhitespace(text string, pos int) int {
	for {
		if pos >= len(text) {
			return pos
		}
		ch := rune(text[pos])
		if ch == ' ' {
			pos++
			continue
		}
		return pos
	}
}

func findPipeIndex(text string) int {
	i := 0
	for {
		if i >= len(text) {
			return -1
		}
		if text[i] != '|' {
			i++
			continue
		}
		return i
	}
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

func isNewLine(ch rune) bool {
	return ch == '\t' || ch == '\n' || ch == '\r'
}

// Return the text and position before new line
func collectUntilNewLine(text string, pos int) (string, int) {
	collected := ""
	for {
		peeked := peek(text, pos)
		if isNewLine(peeked) || peeked == -1 {
			break
		}
		collected += string(text[pos])
		pos++
	}
	return collected, pos - 1
}

func forwardLine(text string, line int, pos int) (int, int) {
	for {
		peeked := peek(text, pos)
		if peeked == -1 {
			break
		}
		if !isNewLine(peeked) {
			break
		}
		pos++
		line++
	}
	return line, pos
}
