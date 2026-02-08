package box

type ValueConverter = func(fromValue float64, fromUnit string, toUnit string) (float64, error)
