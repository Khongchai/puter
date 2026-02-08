package unit

import (
	"fmt"
)

func GetFixedUnitConverter() ValueConverter {
	return func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
		if fromUnit == toUnit {
			return fromValue, nil
		}
		isFromKeyword, _ := IsFixedUnitKeyword(fromUnit)
		isToKeyword, _ := IsFixedUnitKeyword(toUnit)
		if !isFromKeyword || !isToKeyword {
			return -1, fmt.Errorf("Cannot convert %s to %s", fromUnit, toUnit)
		}
		fromDetail := FixedUnitTypes[FixedUnitType(fromUnit)]
		targetDetail := FixedUnitTypes[FixedUnitType(toUnit)]

		if fromDetail.UnitFor != targetDetail.UnitFor {
			return -1, fmt.Errorf("Cannot convert %s to %s", fromUnit, toUnit)
		}

		normalized := fromDetail.ToBaseUnit(fromValue)
		inNewUnit := targetDetail.FromBaseUnit(normalized)
		return inNewUnit, nil
	}
}
