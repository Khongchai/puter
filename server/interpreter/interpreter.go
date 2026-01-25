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

	// get all valid lines as task and store the line index, char start (inclusive), and char end (exclusive)
	// then loop through each valid lines, if there is an interpretation or diagnostics to be added, then add them.

	type EvalTask struct {
		text string
		line int
		// inclusive
		charStart int
		// inclusive
		charEnd int
	}

	tasks := []*EvalTask{}

	for i := range text {
		c := peek(text, i)
		switch c {
		// line comment in python
		case '#':
			collected, stoppedPos := collectUntilNewLine(text, i)
			tasks = append(tasks, &EvalTask{
				text:      collected,
				line:      lineIndex,
				charStart: i,
				charEnd:   stoppedPos,
			})
			newLine, newPos := forwardLine(text, lineIndex, stoppedPos+1)
			lineIndex = newLine
			i = newPos
		case '/':
			peeked := peek(text, i+1)
			// single line comment in c-like languages
			if peeked == '/' {
				collected, stoppedPos := collectUntilNewLine(text, i)
				tasks = append(tasks, &EvalTask{
					text:      collected,
					line:      lineIndex,
					charStart: i,
					charEnd:   stoppedPos,
				})
				newLine, newPos := forwardLine(text, lineIndex, stoppedPos+1)
				lineIndex = newLine
				i = newPos
			}

			// multiline comment in c-like languages
			if peeked == '*' {
				// multiline can also be single line, eg /* something */ so for simplicity,
				// just collect everything until line termination and handle evaluation later
				collected := ""
				j := i
				for {
					peeked := peek(text, j)
					if peeked == -1 {
						break
					}
					if peeked == '*' && peek(text, j+1) == '/' {
						j += 2
						collected += "*/"
						break
					}
					collected += string(peeked)
					j++
				}

				// at this point, collected is /* {...} */ where inside can include newline char too
				lineIndexInner := 0
				j = 0
				for {
					j = skipWhitespace(collected, j)
					collectedInner, stoppedPos := collectUntilNewLine(collected, j)
					tasks = append(tasks, &EvalTask{
						text:      collectedInner,
						line:      lineIndex + lineIndexInner,
						charStart: i,
						charEnd:   stoppedPos,
					})
					newLine, newPos := forwardLine(collected, lineIndexInner, stoppedPos+1)
					lineIndexInner = newLine
					j = newPos
					if j >= len(collected) {
						break
					}
				}
				i += len(collected)
				lineIndex = lineIndexInner
			}
		default:
			if isNewLine(c) {
				i++
				lineIndex++
			}
			i++
		}
	}

	for _, task := range tasks {
		pos := findPipeIndex(task.text)
		peeked := peek(task.text, pos)
		if peeked != '|' {
			continue
		}
		pos++ // skip pipe
		trimmed := task.text[pos:]
		interpretation := interpreter.evaluateAndInterpretResult(evaluator, trimmed, task.line)
		interpretations = append(interpretations, interpretation)
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
			return i
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
		Decoration:  decoration,
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
