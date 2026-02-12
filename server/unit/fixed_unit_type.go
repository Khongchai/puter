// adding new fixed unit type support can be done from this file.

package unit

import (
	"fmt"
	"strings"
)

type FixedUnitType string

type FixedUnitDetail struct {
	UnitFor string

	FullName string

	FullNameSingular string

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	ToBaseUnit func(value float64) float64

	// the equation for transating to base unit
	// what base unit is the smallest unit defined in the group.
	FromBaseUnit func(value float64) float64
}

var FixedUnitTypes = func() map[FixedUnitType]*FixedUnitDetail {
	mapping := map[FixedUnitType]*FixedUnitDetail{
		// Length - Metric (Base: mm)
		"mm": {
			UnitFor:          "length",
			FullName:         "millimeters",
			FullNameSingular: "millimeter",
			ToBaseUnit:       func(value float64) float64 { return value },
			FromBaseUnit:     func(value float64) float64 { return value },
		},
		"cm": {
			UnitFor:          "length",
			FullName:         "centimeters",
			FullNameSingular: "centimeter",
			ToBaseUnit:       func(value float64) float64 { return value * 10 },
			FromBaseUnit:     func(value float64) float64 { return value / 10 },
		},
		"m": {
			UnitFor:          "length",
			FullName:         "meters",
			FullNameSingular: "meter",
			ToBaseUnit:       func(value float64) float64 { return value * 10 * 100 },
			FromBaseUnit:     func(value float64) float64 { return value / 10 / 100 },
		},
		"km": {
			UnitFor:          "length",
			FullName:         "kilometers",
			FullNameSingular: "kilometer",
			ToBaseUnit:       func(value float64) float64 { return value * 10 * 100 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 10 / 100 / 1000 },
		},

		// Length - Imperial (Base: mm)
		"in": {
			UnitFor:          "length",
			FullName:         "inches",
			FullNameSingular: "inch",
			ToBaseUnit:       func(value float64) float64 { return value * 25.4 },
			FromBaseUnit:     func(value float64) float64 { return value / 25.4 },
		},
		"ft": {
			UnitFor:          "length",
			FullName:         "feet",
			FullNameSingular: "foot",
			ToBaseUnit:       func(value float64) float64 { return value * 25.4 * 12 },
			FromBaseUnit:     func(value float64) float64 { return value / 12 / 25.4 },
		},
		"yd": {
			UnitFor:          "length",
			FullName:         "yards",
			FullNameSingular: "yard",
			ToBaseUnit:       func(value float64) float64 { return value * 25.4 * 12 * 3 },
			FromBaseUnit:     func(value float64) float64 { return value / 3 / 12 / 25.4 },
		},
		"mi": {
			UnitFor:          "length",
			FullName:         "miles",
			FullNameSingular: "mile",
			ToBaseUnit:       func(value float64) float64 { return value * 25.4 * 12 * 5280 },
			FromBaseUnit:     func(value float64) float64 { return value / 5280 / 12 / 25.4 },
		},

		// Mass/Weight (Base: mg)
		"mg": {
			UnitFor:          "mass",
			FullName:         "milligrams",
			FullNameSingular: "milligram",
			ToBaseUnit:       func(value float64) float64 { return value },
			FromBaseUnit:     func(value float64) float64 { return value },
		},
		"g": {
			UnitFor:          "mass",
			FullName:         "grams",
			FullNameSingular: "gram",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 },
		},
		"kg": {
			UnitFor:          "mass",
			FullName:         "kilograms",
			FullNameSingular: "kilogram",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 1000 },
		},
		"lbs": {
			UnitFor:          "mass",
			FullName:         "pounds",
			FullNameSingular: "pound",
			ToBaseUnit:       func(value float64) float64 { return value * 453592.37 },
			FromBaseUnit:     func(value float64) float64 { return value / 453592.37 },
		},
		"tonne": {
			UnitFor:          "mass",
			FullName:         "tonnes",
			FullNameSingular: "tonne",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 * 1000 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 1000 / 1000 },
		},

		// Volume (Base: ml)
		"ml": {
			UnitFor:          "volume",
			FullName:         "milliliters",
			FullNameSingular: "milliliter",
			ToBaseUnit:       func(value float64) float64 { return value },
			FromBaseUnit:     func(value float64) float64 { return value },
		},
		"l": {
			UnitFor:          "volume",
			FullName:         "liters",
			FullNameSingular: "liter",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 },
		},

		// Temperature (Base: Celsius)
		"c": {
			UnitFor:          "temperature",
			FullName:         "celsius",
			FullNameSingular: "celsius", // Celsius is generally used for both
			ToBaseUnit:       func(value float64) float64 { return value },
			FromBaseUnit:     func(value float64) float64 { return value },
		},
		"f": {
			UnitFor:          "temperature",
			FullName:         "fahrenheit",
			FullNameSingular: "fahrenheit",
			ToBaseUnit:       func(value float64) float64 { return (value - 32) * 5 / 9 },
			FromBaseUnit:     func(value float64) float64 { return (value * 9 / 5) + 32 },
		},
		"k": {
			UnitFor:          "temperature",
			FullName:         "kelvin",
			FullNameSingular: "kelvin",
			ToBaseUnit:       func(value float64) float64 { return value - 273.15 },
			FromBaseUnit:     func(value float64) float64 { return value + 273.15 },
		},

		// Time (Base: ms)
		"ms": {
			UnitFor:          "time",
			FullName:         "milliseconds",
			FullNameSingular: "millisecond",
			ToBaseUnit:       func(value float64) float64 { return value },
			FromBaseUnit:     func(value float64) float64 { return value },
		},
		"s": {
			UnitFor:          "time",
			FullName:         "seconds",
			FullNameSingular: "second",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 },
		},
		"min": {
			UnitFor:          "time",
			FullName:         "minutes",
			FullNameSingular: "minute",
			ToBaseUnit:       func(value float64) float64 { return value * 60 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 60 },
		},
		"hr": {
			UnitFor:          "time",
			FullName:         "hours",
			FullNameSingular: "hour",
			ToBaseUnit:       func(value float64) float64 { return value * 60 * 60 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 60 / 60 },
		},
		"day": {
			UnitFor:          "time",
			FullName:         "days",
			FullNameSingular: "day",
			ToBaseUnit:       func(value float64) float64 { return value * 24 * 60 * 60 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 60 / 60 / 24 },
		},
		"year": {
			UnitFor:          "time",
			FullName:         "years",
			FullNameSingular: "year",
			ToBaseUnit:       func(value float64) float64 { return value * 365 * 24 * 60 * 60 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 60 / 60 / 24 / 365 },
		},

		// Storage (Base: b)
		"b": {
			UnitFor:          "storage",
			FullName:         "bytes",
			FullNameSingular: "byte",
			ToBaseUnit:       func(value float64) float64 { return value },
			FromBaseUnit:     func(value float64) float64 { return value },
		},
		"kb": {
			UnitFor:          "storage",
			FullName:         "kilobytes",
			FullNameSingular: "kilobyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 },
		},
		"mb": {
			UnitFor:          "storage",
			FullName:         "megabytes",
			FullNameSingular: "megabyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 1000 },
		},
		"gb": {
			UnitFor:          "storage",
			FullName:         "gigabytes",
			FullNameSingular: "gigabyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 * 1000 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 1000 / 1000 },
		},
		"tb": {
			UnitFor:          "storage",
			FullName:         "terabytes",
			FullNameSingular: "terabyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 * 1000 * 1000 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 1000 / 1000 / 1000 },
		},
		"pb": {
			UnitFor:          "storage",
			FullName:         "petabytes",
			FullNameSingular: "petabyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 * 1000 * 1000 * 1000 * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 / 1000 / 1000 / 1000 / 1000 },
		},

		// Binary Storage (Base: b, Power of 2)
		"kib": {
			UnitFor:          "storage",
			FullName:         "kibibytes",
			FullNameSingular: "kibibyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1024 },
			FromBaseUnit:     func(value float64) float64 { return value / 1024 },
		},
		"mib": {
			UnitFor:          "storage",
			FullName:         "mebibytes",
			FullNameSingular: "mebibyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1024 * 1024 },
			FromBaseUnit:     func(value float64) float64 { return value / 1024 / 1024 },
		},
		"gib": {
			UnitFor:          "storage",
			FullName:         "gibibytes",
			FullNameSingular: "gibibyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1024 * 1024 * 1024 },
			FromBaseUnit:     func(value float64) float64 { return value / 1024 / 1024 / 1024 },
		},
		"tib": {
			UnitFor:          "storage",
			FullName:         "tebibytes",
			FullNameSingular: "tebibyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1024 * 1024 * 1024 * 1024 },
			FromBaseUnit:     func(value float64) float64 { return value / 1024 / 1024 / 1024 / 1024 },
		},
		"pib": {
			UnitFor:          "storage",
			FullName:         "pebibytes",
			FullNameSingular: "pebibyte",
			ToBaseUnit:       func(value float64) float64 { return value * 1024 * 1024 * 1024 * 1024 * 1024 },
			FromBaseUnit:     func(value float64) float64 { return value / 1024 / 1024 / 1024 / 1024 / 1024 },
		},

		// Power (Base: Watts - W)
		"w": {
			UnitFor:          "power",
			FullName:         "watts",
			FullNameSingular: "watt",
			ToBaseUnit:       func(value float64) float64 { return value },
			FromBaseUnit:     func(value float64) float64 { return value },
		},
		"kw": {
			UnitFor:          "power",
			FullName:         "kilowatts",
			FullNameSingular: "kilowatt",
			ToBaseUnit:       func(value float64) float64 { return value * 1000 },
			FromBaseUnit:     func(value float64) float64 { return value / 1000 },
		},
		"hp": {
			UnitFor:          "power",
			FullName:         "horsepower",
			FullNameSingular: "horsepower",
			ToBaseUnit:       func(value float64) float64 { return value * 745.7 },
			FromBaseUnit:     func(value float64) float64 { return value / 745.7 },
		},
	}

	for _, value := range mapping {
		mapping[FixedUnitType(value.FullName)] = value
		if value.FullNameSingular == "" {
			panic(fmt.Sprintf("Missing singular for %s", value.FullName))
		}
		mapping[FixedUnitType(value.FullNameSingular)] = value
	}

	return mapping
}()

func IsFixedUnitKeyword(keyword string) (bool, FixedUnitType) {
	lowercased := strings.ToLower(keyword)
	_, is := FixedUnitTypes[FixedUnitType(lowercased)]
	if is {
		return true, FixedUnitType(lowercased)
	}
	return false, ""
}
