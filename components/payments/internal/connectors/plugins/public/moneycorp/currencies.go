package moneycorp

import "github.com/formancehq/payments/internal/connectors/plugins/currency"

var (
	supportedCurrenciesWithDecimal = map[string]int{
		"AED": currency.ISO4217Currencies["AED"], //  UAE Dirham
		"AUD": currency.ISO4217Currencies["AUD"], //  Australian Dollar
		"BBD": currency.ISO4217Currencies["BBD"], //  Barbados Dollar
		"BGN": currency.ISO4217Currencies["BGN"], //  Bulgarian lev
		"BHD": currency.ISO4217Currencies["BHD"], //  Bahraini dinar
		"BWP": currency.ISO4217Currencies["BWP"], //  Botswana pula
		"CAD": currency.ISO4217Currencies["CAD"], //  Canadian dollar
		"CHF": currency.ISO4217Currencies["CHF"], //  Swiss franc
		"CZK": currency.ISO4217Currencies["CZK"], //  Czech koruna
		"DKK": currency.ISO4217Currencies["DKK"], //  Danish krone
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"GBP": currency.ISO4217Currencies["GBP"], //  Pound sterling
		"GHS": currency.ISO4217Currencies["GHS"], //  Ghanaian cedi
		"HKD": currency.ISO4217Currencies["HKD"], //  Hong Kong dollar
		"ILS": currency.ISO4217Currencies["ILS"], //  Israeli new shekel
		"INR": currency.ISO4217Currencies["INR"], //  Indian rupee
		"JMD": currency.ISO4217Currencies["JMD"], //  Jamaican dollar
		"JPY": currency.ISO4217Currencies["JPY"], //  Japanese yen
		"KES": currency.ISO4217Currencies["KES"], //  Kenyan shilling
		"LKR": currency.ISO4217Currencies["LKR"], //  Sri Lankan rupee
		"MAD": currency.ISO4217Currencies["MAD"], //  Moroccan dirham
		"MUR": currency.ISO4217Currencies["MUR"], //  Mauritian rupee
		"MXN": currency.ISO4217Currencies["MXN"], //  Mexican peso
		"NOK": currency.ISO4217Currencies["NOK"], //  Norwegian krone
		"NPR": currency.ISO4217Currencies["NPR"], //  Nepalese rupee
		"NZD": currency.ISO4217Currencies["NZD"], //  New Zealand dollar
		"OMR": currency.ISO4217Currencies["OMR"], //  Omani rial
		"PHP": currency.ISO4217Currencies["PHP"], //  Philippine peso
		"PKR": currency.ISO4217Currencies["PKR"], //  Pakistani rupee
		"PLN": currency.ISO4217Currencies["PLN"], //  Polish złoty
		"QAR": currency.ISO4217Currencies["QAR"], //  Qatari riyal
		"RON": currency.ISO4217Currencies["RON"], //  Romanian leu
		"RSD": currency.ISO4217Currencies["RSD"], //  Serbian dinar
		"SAR": currency.ISO4217Currencies["SAR"], //  Saudi riyal
		"SEK": currency.ISO4217Currencies["SEK"], //  Swedish krona/kronor
		"SGD": currency.ISO4217Currencies["SGD"], //  Singapore dollar
		"THB": currency.ISO4217Currencies["THB"], //  Thai baht
		"TRY": currency.ISO4217Currencies["TRY"], //  Turkish lira
		"TTD": currency.ISO4217Currencies["TTD"], //  Trinidad and Tobago dollar
		"UGX": currency.ISO4217Currencies["UGX"], //  Ugandan shilling
		"USD": currency.ISO4217Currencies["USD"], //  United States dollar
		"XCD": currency.ISO4217Currencies["XCD"], //  East Caribbean dollar
		"ZAR": currency.ISO4217Currencies["ZAR"], //  South African rand
		"ZMW": currency.ISO4217Currencies["ZMW"], //  Zambian kwacha

		// All following currencies have not the same decimals given by
		// Moneycorp compared to the ISO 4217 standard.
		// Let's not handle them for now.
		// "CNH": 2, //  Chinese Yuan
		// "IDR": 2, //  Indonesian rupiah
		// "ISK": 0, //  Icelandic króna
		// "HUF": 2, //  Hungarian forint
		// "JOD": 3, //  Jordanian dinar
		// "KWD": 3, //  Kuwaiti dinar
		// "TND": 3, //  Tunisian dinar
	}
)
