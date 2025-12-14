package box

import "fmt"

type BooleanBox struct {
	Value bool
}

func (bb *BooleanBox) Inspect() string {
	return fmt.Sprintf("%t", bb.Value)
}

func (bb *BooleanBox) Type() BoxType {
	return BOOLEAN_BOX
}
