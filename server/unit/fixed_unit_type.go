// adding new fixed unit type support can be done from this file.

package unit

import "strings"

type FixedUnitType string

type FixedUnitDetail struct {
	UnitFor string

	FullName string

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	ToBaseUnit func(value float64) float64

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	FromBaseUnit func(value float64) float64
}

var FixedUnitTypes = map[FixedUnitType]*FixedUnitDetail{
	// Length - Metric (Base: mm)
	"mm": {
		UnitFor:      "length",
		FullName:     "millimeters",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"cm": {
		UnitFor:      "length",
		FullName:     "centimeters",
		ToBaseUnit:   func(value float64) float64 { return value * 10 },
		FromBaseUnit: func(value float64) float64 { return value / 10 },
	},
	"m": {
		UnitFor:      "length",
		FullName:     "meters",
		ToBaseUnit:   func(value float64) float64 { return value * 10 * 100 },
		FromBaseUnit: func(value float64) float64 { return value / 10 / 100 },
	},
	"km": {
		UnitFor:      "length",
		FullName:     "kilometers",
		ToBaseUnit:   func(value float64) float64 { return value * 10 * 100 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 10 / 100 / 1000 },
	},

	// Length - Imperial (Base: mm)
	"in": {
		UnitFor:      "length",
		FullName:     "inches",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 },
		FromBaseUnit: func(value float64) float64 { return value / 25.4 },
	},
	"ft": {
		UnitFor:      "length",
		FullName:     "feet",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 * 12 },
		FromBaseUnit: func(value float64) float64 { return value / 12 / 25.4 },
	},
	"yd": {
		UnitFor:      "length",
		FullName:     "yards",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 * 12 * 3 },
		FromBaseUnit: func(value float64) float64 { return value / 3 / 12 / 25.4 },
	},
	"mi": {
		UnitFor:      "length",
		FullName:     "miles",
		ToBaseUnit:   func(value float64) float64 { return value * 25.4 * 12 * 5280 },
		FromBaseUnit: func(value float64) float64 { return value / 5280 / 12 / 25.4 },
	},

	// Mass/Weight (Base: mg)
	"mg": {
		UnitFor:      "mass",
		FullName:     "milligrams",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"g": {
		UnitFor:      "mass",
		FullName:     "grams",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 },
	},
	"kg": {
		UnitFor:      "mass",
		FullName:     "kilograms",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 },
	},
	// https://en.wikipedia.org/wiki/Avoirdupois#:~:text=This%20term%20originally%20referred%20to,used%20to%20weigh%20such%20merchandise.
	// 1 lb = 453,592.37 mg
	"lbs": {
		UnitFor:      "mass",
		FullName:     "pounds",
		ToBaseUnit:   func(value float64) float64 { return value * 453592.37 },
		FromBaseUnit: func(value float64) float64 { return value / 453592.37 },
	},
	"ton": {
		UnitFor:      "mass",
		FullName:     "tons", // Metric Tonne
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 / 1000 },
	},

	// Volume (Base: ml)
	"ml": {
		UnitFor:      "volume",
		FullName:     "milliliters",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"l": {
		UnitFor:      "volume",
		FullName:     "liters",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 },
	},

	// Temperature (Base: Celsius)
	"c": {
		UnitFor:      "temperature",
		FullName:     "celsius",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"f": {
		UnitFor:      "temperature",
		FullName:     "fahrenheit",
		ToBaseUnit:   func(value float64) float64 { return (value - 32) * 5 / 9 },
		FromBaseUnit: func(value float64) float64 { return (value * 9 / 5) + 32 },
	},
	"k": {
		UnitFor:      "temperature",
		FullName:     "kelvin",
		ToBaseUnit:   func(value float64) float64 { return value - 273.15 },
		FromBaseUnit: func(value float64) float64 { return value + 273.15 },
	},

	// Time (Base: ms)
	"ms": {
		UnitFor:      "time",
		FullName:     "milliseconds",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"s": {
		UnitFor:      "time",
		FullName:     "seconds",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 },
	},
	"min": {
		UnitFor:      "time",
		FullName:     "minutes",
		ToBaseUnit:   func(value float64) float64 { return value * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 },
	},
	"hr": {
		UnitFor:      "time",
		FullName:     "hours",
		ToBaseUnit:   func(value float64) float64 { return value * 60 * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 / 60 },
	},
	"day": {
		UnitFor:      "time",
		FullName:     "days",
		ToBaseUnit:   func(value float64) float64 { return value * 24 * 60 * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 / 60 / 24 },
	},
	"year": {
		UnitFor:      "time",
		FullName:     "years",
		ToBaseUnit:   func(value float64) float64 { return value * 365 * 24 * 60 * 60 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 60 / 60 / 24 / 365 },
	},

	// storage
	"b": {
		UnitFor:      "storage",
		FullName:     "bytes",
		ToBaseUnit:   func(value float64) float64 { return value },
		FromBaseUnit: func(value float64) float64 { return value },
	},
	"kb": {
		UnitFor:      "storage",
		FullName:     "kilobytes",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 },
	},
	"mb": {
		UnitFor:      "storage",
		FullName:     "megabytes",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 },
	},
	"gb": {
		UnitFor:      "storage",
		FullName:     "gigabytes",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 / 1000 },
	},
	"tb": {
		UnitFor:      "storage",
		FullName:     "terabytes",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 / 1000 / 1000 },
	},
	"pb": {
		UnitFor:      "storage",
		FullName:     "petabytes",
		ToBaseUnit:   func(value float64) float64 { return value * 1000 * 1000 * 1000 * 1000 * 1000 },
		FromBaseUnit: func(value float64) float64 { return value / 1000 / 1000 / 1000 / 1000 / 1000 },
	},
}

func IsFixedUnitKeyword(keyword string) (bool, FixedUnitType) {
	lowercased := strings.ToLower(keyword)
	_, is := FixedUnitTypes[FixedUnitType(lowercased)]
	if is {
		return true, FixedUnitType(lowercased)
	}
	return false, ""
}
