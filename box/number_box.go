package box

import (
	"fmt"
	"puter/ast"
)

type NumberBox struct {
	Value float64
	Tok   *ast.Token
}

func (nb *NumberBox) Inspect() string {
	formatted := fmt.Sprintf("%g", nb.Value)
	return formatted
}

func (nb *NumberBox) Type() BoxType {
	return NUMBER_BOX
}
