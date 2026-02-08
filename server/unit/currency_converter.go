package unit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetCurrencyConverter() ValueConverter {
	return func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
		if fromUnit == toUnit {
			return fromValue, nil
		}

		fromUnit = strings.ToUpper(fromUnit)
		toUnit = strings.ToUpper(toUnit)
		if _, unitIsCurrency := FiatCurrencies[fromUnit]; !unitIsCurrency {
			return -1, fmt.Errorf("%s is not a valid ISO 4217 currency code.", fromUnit)
		}
		if _, unitIsCurrency := FiatCurrencies[toUnit]; !unitIsCurrency {
			return -1, fmt.Errorf("%s is not a valid ISO 4217 currency code.", toUnit)
		}
		ok := isCurrencyConversionSupported(fromUnit, toUnit)
		if !ok {
			return -1, fmt.Errorf("Conversion between %s and %s not supported", fromUnit, toUnit)
		}

		conversionRate, err := fetchCurrencyConversionRate(fromUnit, toUnit)
		if err != nil {
			return 0.0, err
		}

		converted := conversionRate * fromValue

		return converted, nil
	}
}

func isCurrencyConversionSupported(fromCur Currency, toCur Currency) bool {
	_, ok := FiatCurrencies[fromCur]
	if !ok {
		return false
	}
	_, ok = FiatCurrencies[toCur]
	if !ok {
		return false
	}

	return true
}

type FrankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

var responseCache = map[string]*FrankfurterResponse{}

func fetchCurrencyConversionRate(fromUnit string, toUnit string) (float64, error) {
	data, err := func() (*FrankfurterResponse, error) {
		if result, ok := responseCache[fromUnit]; ok {
			return result, nil
		}

		resp, err := http.Get(fmt.Sprintf("https://api.frankfurter.dev/v1/latest?base=%s", fromUnit))
		if err != nil {
			return nil, errors.New("Request to frankfruter api failed")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("Could not read response body")
		}

		var data FrankfurterResponse
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, errors.New("Could not read response body")
		}

		responseCache[fromUnit] = &data

		return &data, nil
	}()

	if err != nil {
		return -1, err
	}

	rate, ok := data.Rates[toUnit]
	if !ok {
		return -1, fmt.Errorf("Conversion rate not found between %s and %s", fromUnit, toUnit)
	}

	return rate, nil
}
