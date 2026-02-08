// For all fixed-measurement types

package box

import (
	"fmt"
	"puter/unit"
)

type MeasurementBox struct {
	Value           *NumberBox
	MeasurementType unit.MeasurementType
}

func NewMeasurementBox(value *NumberBox, measurementType unit.MeasurementType) *MeasurementBox {
	return &MeasurementBox{Value: value, MeasurementType: measurementType}
}

func (mb *MeasurementBox) Inspect() string {
	fullName := unit.MeasurementTypes[mb.MeasurementType].FullName
	return fmt.Sprintf("%g %s", mb.Value.Value, fullName)
}

func (nb *MeasurementBox) Type() BoxType {
	return MEASUREMENT_BOX
}

var _ BinaryNumberOperatable = (*MeasurementBox)(nil)

func (left *MeasurementBox) OperateBinary(right Box, operator BinaryOperation[float64], converters *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return NewMeasurementBox(NewNumberbox(operator(left.Value.Value, r.Value), r.NumberType), left.MeasurementType), nil
	case *MeasurementBox:
		{
			if left.MeasurementType == r.MeasurementType {
				return NewMeasurementBox(NewNumberbox(operator(left.Value.Value, r.Value.Value), r.Value.NumberType), left.MeasurementType), nil
			}
			leftInRight, err := converters.ConvertMeasurement(left.Value.Value, string(left.MeasurementType), string(r.MeasurementType))
			if err != nil {
				return nil, err
			}
			return NewMeasurementBox(NewNumberbox(operator(leftInRight, r.Value.Value), r.Value.NumberType), r.MeasurementType), nil
		}
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", left.Type(), right.Type())
	}
}

var _ InPrefixOperatable = (*MeasurementBox)(nil)

func (mb *MeasurementBox) OperateIn(keyword string, converters *unit.Converters) (Box, error) {
	if keyword == string(mb.MeasurementType) {
		return NewMeasurementBox(mb.Value, mb.MeasurementType), nil
	}

	inNewUnit, err := converters.ConvertMeasurement(mb.Value.Value, string(mb.MeasurementType), keyword)
	if err != nil {
		return nil, err
	}

	return NewMeasurementBox(NewNumberbox(inNewUnit, mb.Value.NumberType), unit.MeasurementType(keyword)), nil
}
