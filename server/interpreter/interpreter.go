package interpreter

import (
	"context"
	"puter/evaluation/evaluator"
	lsproto "puter/lsp"
	"puter/utils"
)

type Interpreter struct {
	ctx               context.Context
	currencyConverter evaluator.ValueConverter
}

type Interpretation struct {
	LineIndex   int
	Decoration  string
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

	// for each line, check if comment patterns is found at the start of line, not counting
	// whitespace, if so, it's a valid line.
	// If line is not
	lineIndex := 0

	interpretations := []*Interpretation{}

	maybeHandlePipeAndForwardLine := func(pos int, line int) (int, int) {
		forwardedPos := skipEverythingUntilPipeOrNewline(text, pos)
		if text[forwardedPos] == '|' {
			collected, nextPos := collectUntilNewLine(text, forwardedPos+1)
			interpretation := interpreter.evaluateAndInterpretResult(evaluator, collected, lineIndex)
			line++
			nextPos++

			interpretations = append(interpretations, interpretation)

			return nextPos, line
		}
		return forwardedPos, line
	}

	for i := range text {
		c := peek(text, i)
		switch c {
		// line comment in python
		case '#':
			nextPos, nextLine := maybeHandlePipeAndForwardLine(i+2, lineIndex)
			lineIndex = nextLine
			i = nextPos
		case '/':
			peeked := peek(text, i+1)
			// single line comment in c-like languages
			if peeked == '/' {
				nextPos, nextLine := maybeHandlePipeAndForwardLine(i+2, lineIndex)
				lineIndex = nextLine
				i = nextPos
				continue
			}

			// multiline comment in c-like languages
			if peeked == '*' {
				for {
					nextPos, nextLine := maybeHandlePipeAndForwardLine(i+2, lineIndex)
					lineIndex = nextLine
					i = nextPos

					i := skipWhitespace(text, i)
					if peek(text, i) == '*' && peek(text, i+1) == '/' {
						i += 2
						break
					}
				}
			}
		default:
			if isNewLine(c) {
				i++
				lineIndex++
			}
			i++
		}
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

func skipEverythingUntilPipeOrNewline(text string, pos int) int {
	for {
		if pos >= len(text) {
			return pos
		}
		ch := rune(text[pos])
		if ch != '|' && !isNewLine(ch) {
			pos++
			continue
		}
		return pos
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
		Decoration:  decoration,
	}
}

func isWhiteSpace(ch byte) bool {
	return ch == ' '
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
	return collected, pos
}
