package box

import (
	"fmt"
	"puter/unit"
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

var _ BinaryNumberOperatable = (*PercentBox)(nil)

func (pb *PercentBox) OperateBinary(right Box, operator BinaryOperation[float64], _ *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return &NumberBox{Value: operator(r.Value, (pb.Value/100)*r.Value)}, nil
	case *CurrencyBox:
		return &CurrencyBox{Number: NewNumberbox(operator(r.Number.Value, (pb.Value/100)*r.Number.Value), r.Number.NumberType), Unit: r.Unit}, nil
	case *PercentBox:
		return &PercentBox{Value: operator(pb.Value, r.Value)}, nil
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", pb.Type(), right.Type())
	}
}
