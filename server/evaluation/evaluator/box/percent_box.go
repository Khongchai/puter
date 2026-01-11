package box

import (
	"fmt"
)

type PercentBox struct {
	Value float64
}

func (nb *PercentBox) Inspect() string {
	formatted := fmt.Sprintf("%g%%", nb.Value)
	return formatted
}

func (nb *PercentBox) Type() BoxType {
	return PERCENT_BOX
}
