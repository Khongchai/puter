package unit

import "strings"

type MeasurementType string

type MeasurementDetail struct {
	// what this Measures, for example length, length-imperial, etc
	Measures string

	FullName string

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	ToBaseUnit func(value float64) float64

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	FromBaseUnit func(value float64) float64
}

var MeasurementTypes = map[MeasurementType]*MeasurementDetail{
	// Length - Metric (Base: mm)
	"mm": {
		Measures:     "length",
		FullName:     "millimeters",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"cm": {
		Measures:     "length",
		FullName:     "centimeters",
		ToBaseUnit:   func(value float64) float64 { return value * 10 },
		FromBaseUnit: func(value float64) float64 { return value / 10 },
	},
	"m": {
		Measures:     "length",
		FullName:     "meters",
		ToBaseUnit:   func(value float64) float64 { return value * 10 * 100 },
		FromBaseUnit: func(value float64) float64 { return value / 10 / 100 },
	},
	"km": {
		Measures:     "length",
		FullName:     "kilometers",
		ToBaseUnit:   func(value float64) float64 { return value * 10 * 100 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 10 / 100 / 1000 },
	},

	// Length - Imperial (Base: mm)
	"in": {
		Measures:     "length",
		FullName:     "inches",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 },
		FromBaseUnit: func(value float64) float64 { return value / 25.4 },
	},
	"ft": {
		Measures:     "length",
		FullName:     "feet",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 * 12 },
		FromBaseUnit: func(value float64) float64 { return value / 12 / 25.4 },
	},
	"yd": {
		Measures:     "length",
		FullName:     "yards",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 * 12 * 3 },
		FromBaseUnit: func(value float64) float64 { return value / 3 / 12 / 25.4 },
	},
	"mi": {
		Measures:     "length",
		FullName:     "miles",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 * 12 * 5280 },
		FromBaseUnit: func(value float64) float64 { return value / 5280 / 12 / 25.4 },
	},

	// Mass/Weight (Base: mg)
	"mg": {
		Measures:     "mass",
		FullName:     "milligrams",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"g": {
		Measures:     "mass",
		FullName:     "grams",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 },
	},
	"kg": {
		Measures:     "mass",
		FullName:     "kilograms",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 },
	},
	// https://en.wikipedia.org/wiki/Avoirdupois#:~:text=This%20term%20originally%20referred%20to,used%20to%20weigh%20such%20merchandise.
	// 1 lb = 453,592.37 mg
	"lbs": {
		Measures:     "mass",
		FullName:     "pounds",
		ToBaseUnit:   func(value float64) float64 { return value * 453592.37 },
		FromBaseUnit: func(value float64) float64 { return value / 453592.37 },
	},
	"ton": {
		Measures:     "mass",
		FullName:     "tons", // Metric Tonne
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 / 1000 },
	},

	// Volume (Base: ml)
	"ml": {
		Measures:     "volume",
		FullName:     "milliliters",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"l": {
		Measures:     "volume",
		FullName:     "liters",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 },
	},

	// Temperature (Base: Celsius)
	"c": {
		Measures:     "temperature",
		FullName:     "celsius",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"f": {
		Measures:     "temperature",
		FullName:     "fahrenheit",
		ToBaseUnit:   func(value float64) float64 { return (value - 32) * 5 / 9 },
		FromBaseUnit: func(value float64) float64 { return (value * 9 / 5) + 32 },
	},
	"k": {
		Measures:     "temperature",
		FullName:     "kelvin",
		ToBaseUnit:   func(value float64) float64 { return value - 273.15 },
		FromBaseUnit: func(value float64) float64 { return value + 273.15 },
	},

	// Time (Base: ms)
	"ms": {
		Measures:     "time",
		FullName:     "milliseconds",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"s": {
		Measures:     "time",
		FullName:     "seconds",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 },
	},
	"min": {
		Measures:     "time",
		FullName:     "minutes",
		ToBaseUnit:   func(value float64) float64 { return value * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 },
	},
	"hr": {
		Measures:     "time",
		FullName:     "hours",
		ToBaseUnit:   func(value float64) float64 { return value * 60 * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 / 60 },
	},
	"day": {
		Measures:     "time",
		FullName:     "days",
		ToBaseUnit:   func(value float64) float64 { return value * 24 * 60 * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 / 60 / 24 },
	},
	"year": {
		Measures:     "time",
		FullName:     "years",
		ToBaseUnit:   func(value float64) float64 { return value * 365 * 24 * 60 * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 / 60 / 24 / 365 },
	},
}

func IsMeasurementKeyword(keyword string) (bool, MeasurementType) {
	lowercased := strings.ToLower(keyword)
	_, is := MeasurementTypes[MeasurementType(lowercased)]
	if is {
		return true, MeasurementType(lowercased)
	}
	return false, ""
}
