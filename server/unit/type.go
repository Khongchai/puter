package unit

type ValueConverter = func(fromValue float64, fromUnit string, toUnit string) (float64, error)

type Converters struct {
	ConvertCurrency  ValueConverter
	ConvertFixedUnit ValueConverter
}
