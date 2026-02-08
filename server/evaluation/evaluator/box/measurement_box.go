// For all fixed-measurement types

package box

import (
	"fmt"
	"puter/evaluation/ast"
	"puter/unit"
)

type MeasurementBox struct {
	Number          *NumberBox
	MeasurementType unit.MeasurementType
}

func NewMeasurementBox(value *NumberBox, measurementType unit.MeasurementType) *MeasurementBox {
	return &MeasurementBox{Number: value, MeasurementType: measurementType}
}

func (mb *MeasurementBox) Inspect() string {
	fullName := unit.MeasurementTypes[mb.MeasurementType].FullName
	return fmt.Sprintf("%g %s", mb.Number.Value, fullName)
}

func (nb *MeasurementBox) Type() BoxType {
	return MEASUREMENT_BOX
}

var _ BinaryNumberOperatable = (*MeasurementBox)(nil)

func (left *MeasurementBox) OperateBinaryNumber(right Box, operator func(a, b float64) float64, converters *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return NewMeasurementBox(NewNumberbox(operator(left.Number.Value, r.Value), r.NumberType), left.MeasurementType), nil
	case *MeasurementBox:
		{
			if left.MeasurementType == r.MeasurementType {
				return NewMeasurementBox(NewNumberbox(operator(left.Number.Value, r.Number.Value), r.Number.NumberType), left.MeasurementType), nil
			}
			leftInRight, err := converters.ConvertMeasurement(left.Number.Value, string(left.MeasurementType), string(r.MeasurementType))
			if err != nil {
				return nil, err
			}
			return NewMeasurementBox(NewNumberbox(operator(leftInRight, r.Number.Value), r.Number.NumberType), r.MeasurementType), nil
		}
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", left.Type(), right.Type())
	}
}

var _ InPrefixOperatable = (*MeasurementBox)(nil)

func (mb *MeasurementBox) OperateIn(keyword string, converters *unit.Converters) (Box, error) {
	if keyword == string(mb.MeasurementType) {
		return NewMeasurementBox(mb.Number, mb.MeasurementType), nil
	}

	inNewUnit, err := converters.ConvertMeasurement(mb.Number.Value, string(mb.MeasurementType), keyword)
	if err != nil {
		return nil, err
	}

	return NewMeasurementBox(NewNumberbox(inNewUnit, mb.Number.NumberType), unit.MeasurementType(keyword)), nil
}

var _ BinaryBooleanOperatable = (*MeasurementBox)(nil)

func (left *MeasurementBox) OperateBinaryBoolean(right Box, operator *ast.Token, converters *unit.Converters) (Box, error) {
	r, is := right.(*MeasurementBox)
	if !is {
		return &BooleanBox{Value: false}, nil
	}

	var leftAsRight, err = converters.ConvertMeasurement(left.Number.Value, string(left.MeasurementType), string(r.MeasurementType))
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
