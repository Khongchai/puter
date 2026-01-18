package interpreter

import "puter/evaluation/evaluator"

type EvaluatorState struct {
	StartLineInclusive int
	EndLineInclusive   int
	Evaluator          *evaluator.Evaluator
}

type DocumentUri = string

type Interpreter struct {
	evaluatorStates map[DocumentUri][]*EvaluatorState
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
func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (interpreter *Interpreter) Interpret(uri string, line int, lineText string) {
}
