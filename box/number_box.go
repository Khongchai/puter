package box

import (
	"fmt"
	"puter/ast"
	"puter/lib"
)

type NumberBox struct {
	Value float64
	Tok   *ast.Token
}

func (nb *NumberBox) Inspect() *lib.Promise[string] {
	formatted := fmt.Sprintf("%g", nb.Value)
	return lib.NewResolvedPromise(formatted)
}

func (nb *NumberBox) Type() BoxType {
	return NUMBER_BOX
}
