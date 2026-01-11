package interpreter

import "puter/evaluation/evaluator"

type FileState struct {
	evaluators []*evaluator.Evaluator
}

type Interpreter struct {
	evaluator map[int]*evaluator.Evaluator
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func Interpret(lint int, lineText string) {

}
