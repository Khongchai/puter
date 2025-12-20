package box

import (
	"fmt"
)

type Currency = string

type CurrencyBox struct {
	Number *NumberBox
	Unit   Currency
}

func (bb *CurrencyBox) Inspect() string {
	text := fmt.Sprintf("%g %s", bb.Number.Value, bb.Unit)
	return text
}

func (bb *CurrencyBox) Type() BoxType {
	return CURRENCY_BOX
}

// TODO generate this programmatically
var ValidCurrencies = map[Currency]struct{}{
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
