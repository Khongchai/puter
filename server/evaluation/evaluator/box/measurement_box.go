// For all fixed-measurement types

package box

import (
	"fmt"
	"strings"
)

type MeasurementType string

var measurementTypes = map[MeasurementType]string{
	// Length - Metric
	"mm": "millimeters",
	"cm": "centimeters",
	"m":  "meters",
	"km": "kilometers",

	// Length - Imperial
	"in": "inches",
	"ft": "feet",
	"yd": "yards",
	"mi": "miles",

	// Mass/Weight
	"mg":  "milligrams",
	"g":   "grams",
	"kg":  "kilograms",
	"lbs": "pounds",
	"ton": "tons",

	// Volume
	"ml": "milliliters",
	"l":  "liters",

	// Temperature
	"c": "celsius",
	"f": "fahrenheit",
	"k": "kelvin",

	// Time
	"ms":   "milliseconds",
	"s":    "seconds",
	"min":  "minutes",
	"hr":   "hours",
	"day":  "days",
	"year": "years",
}

type MeasurementBox struct {
	Value           *NumberBox
	MeasurementType MeasurementType
}

func IsMeasurementKeyword(keyword string) (bool, MeasurementType) {
	lowercased := strings.ToLower(keyword)
	_, is := measurementTypes[MeasurementType(lowercased)]
	if is {
		return true, MeasurementType(lowercased)
	}
	return false, ""
}

func NewMeasurementBox(value *NumberBox, measurementType MeasurementType) *MeasurementBox {
	return &MeasurementBox{Value: value, MeasurementType: measurementType}
}

func (mb *MeasurementBox) Inspect() string {
	return fmt.Sprintf("%g %s", mb.Value, mb.MeasurementType)
}

func (nb *MeasurementBox) Type() BoxType {
	return MEASUREMENT_BOX
}

var _ BinaryNumberOperatable = (*MeasurementBox)(nil)

func (left *MeasurementBox) OperateBinary(right Box, operator BinaryOperation[float64], _ ValueConverter) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return NewMeasurementBox(NewNumberbox(operator(left.Value.Value, r.Value), r.NumberType), left.MeasurementType), nil
	case *MeasurementBox:
		return NewMeasurementBox(NewNumberbox(operator(left.Value.Value, r.Value.Value), r.Value.NumberType), left.MeasurementType), nil
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", left.Type(), right.Type())
	}
}

var _ InPrefixOperatable = (*MeasurementBox)(nil)

func (mb *MeasurementBox) OperateIn(keyword string, _ ValueConverter) (Box, error) {
	if keyword == string(mb.MeasurementType) {
		return NewMeasurementBox(mb.Value, mb.MeasurementType), nil
	}

	isMeasurementKeyword, measurementType := IsMeasurementKeyword(keyword)
	if isMeasurementKeyword {
		// TODO convert measurement
		return nil, nil
	}

	return nil, fmt.Errorf("Cannot convert %s to %s", mb.MeasurementType, keyword)
}
