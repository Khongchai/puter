package extensions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"puter/engine/evaluator"
	"puter/engine/evaluator/box"
)

func GetCurrencyConverter() evaluator.ValueConverter {
	return func(fromValue float64, fromUnit string, toUnit string) (float64, error) {
		if _, unitIsCurrency := FiatCurrencies[fromUnit]; !unitIsCurrency {
			panic(fmt.Sprintf("%s is not a valid ISO 4217 currency code.", fromUnit))
		}
		if _, unitIsCurrency := FiatCurrencies[toUnit]; !unitIsCurrency {
			panic(fmt.Sprintf("%s is not a valid ISO 4217 currency code.", toUnit))
		}
		ok := isCurrencyConversionSupported(fromUnit, toUnit)
		if !ok {
			panic(fmt.Sprintf("Conversion between %s and %s not supported", fromUnit, toUnit))
		}

		conversionRate, err := fetchCurrencyConversionRate(fromValue, fromUnit, toUnit)
		if err != nil {
			return 0.0, err
		}

		converted := conversionRate * fromValue

		return converted, nil
	}
}

func isCurrencyConversionSupported(fromCur box.Currency, toCur box.Currency) bool {
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

func fetchCurrencyConversionRate(fromValue float64, fromUnit string, toUnit string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.frankfurter.dev/v1/latest?base=%s", fromUnit))
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0, errors.New("Could not read response body")
	}

	var data FrankfurterResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0.0, errors.New("Could not read response body")
	}

	rate, ok := data.Rates[toUnit]
	if !ok {
		return 0.0, fmt.Errorf("Conversion rate not found between %s and %s", fromUnit, toUnit)
	}

	return fromValue * rate, nil
}

// TODO include crypto currencies here
var FiatCurrencies = map[box.Currency]struct{}{
	"USD": {}, // US Dollar
	"EUR": {}, // Euro
	"GBP": {}, // British Pound
	"JPY": {}, // Japanese Yen
	"CNY": {}, // Chinese Yuan
	"CHF": {}, // Swiss Franc
	"CAD": {}, // Canadian Dollar
	"AUD": {}, // Australian Dollar
	"NZD": {}, // New Zealand Dollar

	// Asia
	"THB": {}, // Thai Baht
	"SGD": {}, // Singapore Dollar
	"HKD": {}, // Hong Kong Dollar
	"KRW": {}, // South Korean Won
	"INR": {}, // Indian Rupee
	"IDR": {}, // Indonesian Rupiah
	"MYR": {}, // Malaysian Ringgit
	"PHP": {}, // Philippine Peso
	"VND": {}, // Vietnamese Dong
	"PKR": {}, // Pakistani Rupee
	"BDT": {}, // Bangladeshi Taka
	"LKR": {}, // Sri Lankan Rupee
	"NPR": {}, // Nepalese Rupee
	"KHR": {}, // Cambodian Riel
	"MMK": {}, // Myanmar Kyat
	"LAK": {}, // Lao Kip
	"MNT": {}, // Mongolian Tögrög
	"KZT": {}, // Kazakhstani Tenge
	"UZS": {}, // Uzbekistani Som

	// Middle East
	"AED": {}, // UAE Dirham
	"SAR": {}, // Saudi Riyal
	"QAR": {}, // Qatari Riyal
	"KWD": {}, // Kuwaiti Dinar
	"BHD": {}, // Bahraini Dinar
	"OMR": {}, // Omani Rial
	"ILS": {}, // Israeli New Shekel
	"JOD": {}, // Jordanian Dinar
	"IRR": {}, // Iranian Rial
	"IQD": {}, // Iraqi Dinar

	// Europe (non-EUR)
	"SEK": {}, // Swedish Krona
	"NOK": {}, // Norwegian Krone
	"DKK": {}, // Danish Krone
	"PLN": {}, // Polish Złoty
	"CZK": {}, // Czech Koruna
	"HUF": {}, // Hungarian Forint
	"RON": {}, // Romanian Leu
	"BGN": {}, // Bulgarian Lev
	"ISK": {}, // Icelandic Króna
	"UAH": {}, // Ukrainian Hryvnia
	"RSD": {}, // Serbian Dinar
	"ALL": {}, // Albanian Lek
	"MKD": {}, // Macedonian Denar
	"BAM": {}, // Bosnia-Herzegovina Convertible Mark
	"MDL": {}, // Moldovan Leu
	"BYN": {}, // Belarusian Ruble
	"RUB": {}, // Russian Ruble

	// Americas
	"MXN": {}, // Mexican Peso
	"BRL": {}, // Brazilian Real
	"ARS": {}, // Argentine Peso
	"CLP": {}, // Chilean Peso
	"COP": {}, // Colombian Peso
	"PEN": {}, // Peruvian Sol
	"UYU": {}, // Uruguayan Peso
	"BOB": {}, // Bolivian Boliviano
	"PYG": {}, // Paraguayan Guaraní
	"VES": {}, // Venezuelan Bolívar
	"DOP": {}, // Dominican Peso
	"CUP": {}, // Cuban Peso
	"JMD": {}, // Jamaican Dollar
	"TTD": {}, // Trinidad and Tobago Dollar

	// Africa
	"ZAR": {}, // South African Rand
	"NGN": {}, // Nigerian Naira
	"KES": {}, // Kenyan Shilling
	"UGX": {}, // Ugandan Shilling
	"TZS": {}, // Tanzanian Shilling
	"GHS": {}, // Ghanaian Cedi
	"ETB": {}, // Ethiopian Birr
	"MAD": {}, // Moroccan Dirham
	"DZD": {}, // Algerian Dinar
	"TND": {}, // Tunisian Dinar
	"EGP": {}, // Egyptian Pound
	"SDG": {}, // Sudanese Pound
	"ZMW": {}, // Zambian Kwacha
	"BWP": {}, // Botswana Pula
	"MUR": {}, // Mauritian Rupee

	// Special / supranational
	"XAF": {}, // Central African CFA Franc
	"XOF": {}, // West African CFA Franc
	"XPF": {}, // CFP Franc
	"XCD": {}, // East Caribbean Dollar
}
