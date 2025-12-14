package box

import (
	"fmt"
)

type NumberBox struct {
	Value float64
}

func (nb *NumberBox) Inspect() string {
	return fmt.Sprintf("%g", nb.Value)
}

func (nb *NumberBox) Type() BoxType {
	return NUMBER_BOX
}
