// For all fixed-measurement types

package box

import (
	"fmt"
	"strings"
)

type MeasurementType string

type MeasurementDetail struct {
	// what this measures, for example length, length-imperial, etc
	measures string

	fullName string

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	toBaseUnit func(value float64) float64

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	fromBaseUnit func(value float64) float64
}

var measurementTypes = map[MeasurementType]*MeasurementDetail{
	// Length - Metric
	"mm": {
		measures:     "length",
		fullName:     "millimeters",
		toBaseUnit:   func(value float64) float64 { return value },
		fromBaseUnit: func(value float64) float64 { return value },
	},
	"cm": {
		measures:     "length",
		fullName:     "centimeters",
		toBaseUnit:   func(value float64) float64 { return value * 10 },
		fromBaseUnit: func(value float64) float64 { return value / 10 },
	},
	"m": {
		measures:     "length",
		fullName:     "meters",
		toBaseUnit:   func(value float64) float64 { return value * 10 * 100 },
		fromBaseUnit: func(value float64) float64 { return value / 10 / 100 },
	},
	"km": {
		measures:     "length",
		fullName:     "kilometers",
		toBaseUnit:   func(value float64) float64 { return value * 10 * 100 * 1000 },
		fromBaseUnit: func(value float64) float64 { return value / 10 / 100 / 1000 },
	},

	// Length - Imperial
	// "in": "inches",
	// "ft": "feet",
	// "yd": "yards",
	// "mi": "miles",

	// // Mass/Weight
	// "mg":  "milligrams",
	// "g":   "grams",
	// "kg":  "kilograms",
	// "lbs": "pounds",
	// "ton": "tons",

	// // Volume
	// "ml": "milliliters",
	// "l":  "liters",

	// // Temperature
	// "c": "celsius",
	// "f": "fahrenheit",
	// "k": "kelvin",

	// // Time
	// "ms":   "milliseconds",
	// "s":    "seconds",
	// "min":  "minutes",
	// "hr":   "hours",
	// "day":  "days",
	// "year": "years",
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
	fullName := measurementTypes[mb.MeasurementType].fullName
	return fmt.Sprintf("%g %s", mb.Value.Value, fullName)
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
	if !isMeasurementKeyword {
		return nil, fmt.Errorf("Cannot convert %s to %s", mb.MeasurementType, keyword)
	}

	fromDetail := measurementTypes[mb.MeasurementType]
	targetDetail := measurementTypes[measurementType]

	normalized := fromDetail.toBaseUnit(mb.Value.Value)
	inNewUnit := targetDetail.fromBaseUnit(normalized)

	return NewMeasurementBox(NewNumberbox(inNewUnit, mb.Value.NumberType), mb.MeasurementType), nil

}
