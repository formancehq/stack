package bankingcircle

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	// All supported BankingCircle currencies and decimal are on par with
	// ISO4217Currencies.
	supportedCurrenciesWithDecimal = map[string]int{
		"AED": currency.ISO4217Currencies["AED"], // UAE Dirham
		"AUD": currency.ISO4217Currencies["AUD"], // Australian Dollar
		"CAD": currency.ISO4217Currencies["CAD"], // Canadian Dollar
		"CHF": currency.ISO4217Currencies["CHF"], // Swiss Franc
		"CNY": currency.ISO4217Currencies["CNY"], //  China Yuan Renminbi
		"CZK": currency.ISO4217Currencies["CZK"], //  Czech Koruna
		"DKK": currency.ISO4217Currencies["DKK"], //  Danish Krone
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"GBP": currency.ISO4217Currencies["GBP"], //  Pound Sterling
		"HKD": currency.ISO4217Currencies["HKD"], //  Hong Kong Dollar
		"ILS": currency.ISO4217Currencies["ILS"], //  Israeli Shekel
		"JPY": currency.ISO4217Currencies["JPY"], //  Japanese Yen
		"MXN": currency.ISO4217Currencies["MXN"], //  Mexican Peso
		"NOK": currency.ISO4217Currencies["NOK"], //  Norwegian Krone
		"NZD": currency.ISO4217Currencies["NZD"], //  New Zealand Dollar
		"PLN": currency.ISO4217Currencies["PLN"], //  Polish Zloty
		"RON": currency.ISO4217Currencies["RON"], //  Romanian Leu
		"SAR": currency.ISO4217Currencies["SAR"], //  Saudi Riyal
		"SEK": currency.ISO4217Currencies["SEK"], //  Swedish Krona
		"SGD": currency.ISO4217Currencies["SGD"], //  Singapore Dollar
		"TRY": currency.ISO4217Currencies["TRY"], //  Turkish Lira
		"USD": currency.ISO4217Currencies["USD"], //  US Dollar
		"ZAR": currency.ISO4217Currencies["ZAR"], //  South African Rand

		// Unsupported currencies
		// Since we're not sure about decimals for these currencies, we prefer
		// to not support them for now.
		// "HUF": 2, //  Hungarian Forint
	}
)
