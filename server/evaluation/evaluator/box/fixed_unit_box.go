package box

import (
	"fmt"
	"puter/evaluation/ast"
	"puter/unit"
)

type FixedUnitBox struct {
	Number        *NumberBox
	FixedUnitType unit.FixedUnitType
}

func NewFixedUnitBox(value *NumberBox, fixedUnitType unit.FixedUnitType) *FixedUnitBox {
	return &FixedUnitBox{Number: value, FixedUnitType: fixedUnitType}
}

func (mb *FixedUnitBox) Inspect() string {
	fullName := unit.FixedUnitTypes[mb.FixedUnitType].FullName
	return fmt.Sprintf("%g %s", mb.Number.Value, fullName)
}

func (nb *FixedUnitBox) Type() BoxType {
	return FIXED_UNIT_BOX
}

var _ BinaryNumberOperatable = (*FixedUnitBox)(nil)

func (left *FixedUnitBox) OperateBinaryNumber(right Box, operator func(a, b float64) float64, converters *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return NewFixedUnitBox(NewNumberbox(operator(left.Number.Value, r.Value), r.NumberType), left.FixedUnitType), nil
	case *FixedUnitBox:
		{
			if left.FixedUnitType == r.FixedUnitType {
				return NewFixedUnitBox(NewNumberbox(operator(left.Number.Value, r.Number.Value), r.Number.NumberType), left.FixedUnitType), nil
			}
			leftInRight, err := converters.ConvertFixedUnit(left.Number.Value, string(left.FixedUnitType), string(r.FixedUnitType))
			if err != nil {
				return nil, err
			}
			return NewFixedUnitBox(NewNumberbox(operator(leftInRight, r.Number.Value), r.Number.NumberType), r.FixedUnitType), nil
		}
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", left.Type(), right.Type())
	}
}

var _ InPrefixOperatable = (*FixedUnitBox)(nil)

func (mb *FixedUnitBox) OperateIn(keyword string, converters *unit.Converters) (Box, error) {
	if keyword == string(mb.FixedUnitType) {
		return NewFixedUnitBox(mb.Number, mb.FixedUnitType), nil
	}

	inNewUnit, err := converters.ConvertFixedUnit(mb.Number.Value, string(mb.FixedUnitType), keyword)
	if err != nil {
		return nil, err
	}

	return NewFixedUnitBox(NewNumberbox(inNewUnit, mb.Number.NumberType), unit.FixedUnitType(keyword)), nil
}

var _ BinaryBooleanOperatable = (*FixedUnitBox)(nil)

func (left *FixedUnitBox) OperateBinaryBoolean(right Box, operator *ast.Token, converters *unit.Converters) (Box, error) {
	r, is := right.(*FixedUnitBox)
	if !is {
		return &BooleanBox{Value: false}, nil
	}

	var leftAsRight, err = converters.ConvertFixedUnit(left.Number.Value, string(left.FixedUnitType), string(r.FixedUnitType))
	if err != nil {
		return nil, err
	}

	result := func() bool {
		switch operator.Type {
		case ast.EQ:
			return leftAsRight == r.Number.Value
		case ast.NOT_EQ:
			return leftAsRight != r.Number.Value
		case ast.LT:
			return leftAsRight < r.Number.Value
		case ast.GT:
			return leftAsRight > r.Number.Value
		case ast.LTE:
			return leftAsRight <= r.Number.Value
		case ast.GTE:
			return leftAsRight >= r.Number.Value
		default:
			return false
		}
	}()

	return &BooleanBox{Value: result}, nil
}

var _ HoldsNumber = (*FixedUnitBox)(nil)

func (m *FixedUnitBox) GetNumber() float64 {
	return m.Number.Value
}
