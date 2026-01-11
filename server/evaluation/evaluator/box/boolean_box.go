package box

import (
	"fmt"
)

type BooleanBox struct {
	Value bool
}

func (bb *BooleanBox) Inspect() string {
	text := fmt.Sprintf("%t", bb.Value)
	return text
}

func (bb *BooleanBox) Type() BoxType {
	return BOOLEAN_BOX
}
