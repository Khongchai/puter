package box

import (
	"fmt"
)

type NumberBox struct {
	Value float64
}

func (nb *NumberBox) Inspect() string {
	formatted := fmt.Sprintf("%g", nb.Value)
	return formatted
}

func (nb *NumberBox) Type() BoxType {
	return NUMBER_BOX
}
