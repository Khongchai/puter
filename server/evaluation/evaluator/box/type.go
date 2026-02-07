package box

import "strings"

type ValueConverter = func(fromValue float64, fromUnit string, toUnit string) (float64, error)

func IsNumberKeyword(keyword string) (bool, NumberType) {
	lowercased := strings.ToLower(keyword)
	is := lowercased == string(Decimal) || keyword == string(Binary) || keyword == string(Hex)
	if is {
		return true, NumberType(lowercased)
	}
	return false, NaN
}
