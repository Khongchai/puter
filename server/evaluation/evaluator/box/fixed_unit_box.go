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

func (fub *FixedUnitBox) Inspect() string {
	fullName := unit.FixedUnitTypes[fub.FixedUnitType].FullName
	return fmt.Sprintf("%g %s", fub.Number.Value, fullName)
}

func (fub *FixedUnitBox) Type() BoxType {
	return FIXED_UNIT_BOX
}

var _ BinaryNumberOperatable = (*FixedUnitBox)(nil)

func (fub *FixedUnitBox) OperateBinaryNumber(right Box, operator func(a, b float64) float64, converters *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return NewFixedUnitBox(NewNumberbox(operator(fub.Number.Value, r.Value), r.NumberType), fub.FixedUnitType), nil
	case *FixedUnitBox:
		{
			if fub.FixedUnitType == r.FixedUnitType {
				return NewFixedUnitBox(NewNumberbox(operator(fub.Number.Value, r.Number.Value), r.Number.NumberType), fub.FixedUnitType), nil
			}
			leftInRight, err := converters.ConvertFixedUnit(fub.Number.Value, string(fub.FixedUnitType), string(r.FixedUnitType))
			if err != nil {
				return nil, err
			}
			return NewFixedUnitBox(NewNumberbox(operator(leftInRight, r.Number.Value), r.Number.NumberType), r.FixedUnitType), nil
		}
	default:
		return nil, fmt.Errorf("Cannot perform this operation on these unit types")
	}
}

var _ InPrefixOperatable = (*FixedUnitBox)(nil)

func (fub *FixedUnitBox) OperateIn(keyword string, converters *unit.Converters) (Box, error) {
	if keyword == string(fub.FixedUnitType) {
		return NewFixedUnitBox(fub.Number, fub.FixedUnitType), nil
	}

	inNewUnit, err := converters.ConvertFixedUnit(fub.Number.Value, string(fub.FixedUnitType), keyword)
	if err != nil {
		return nil, err
	}

	return NewFixedUnitBox(NewNumberbox(inNewUnit, fub.Number.NumberType), unit.FixedUnitType(keyword)), nil
}

var _ BinaryBooleanOperatable = (*FixedUnitBox)(nil)

func (fub *FixedUnitBox) OperateBinaryBoolean(right Box, operator *ast.Token, converters *unit.Converters) (Box, error) {
	r, is := right.(*FixedUnitBox)
	if !is {
		return &BooleanBox{Value: false}, nil
	}

	var leftAsRight, err = converters.ConvertFixedUnit(fub.Number.Value, string(fub.FixedUnitType), string(r.FixedUnitType))
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

var _ NumericType = (*FixedUnitBox)(nil)

func (fub *FixedUnitBox) GetNumber() float64 {
	return fub.Number.Value
}

func (fub *FixedUnitBox) SetNumber(v float64) {
	fub.Number.Value = v
}

func (fub *FixedUnitBox) Clone() Box {
	return &FixedUnitBox{Number: NewNumberbox(fub.Number.Value, fub.Number.NumberType), FixedUnitType: fub.FixedUnitType}
}
