package unit

import (
	"fmt"
)

func GetMeasurementConverter() ValueConverter {
	return func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
		if fromUnit == toUnit {
			return fromValue, nil
		}
		isFromKeyword, _ := IsMeasurementKeyword(fromUnit)
		isToKeyword, _ := IsMeasurementKeyword(toUnit)
		if !isFromKeyword || !isToKeyword {
			return -1, fmt.Errorf("Cannot convert %s to %s", fromUnit, toUnit)
		}
		fromDetail := MeasurementTypes[MeasurementType(fromUnit)]
		targetDetail := MeasurementTypes[MeasurementType(toUnit)]

		if fromDetail.Measures != targetDetail.Measures {
			return -1, fmt.Errorf("Cannot convert %s to %s", fromUnit, toUnit)
		}

		normalized := fromDetail.ToBaseUnit(fromValue)
		inNewUnit := targetDetail.FromBaseUnit(normalized)
		return inNewUnit, nil
	}
}
