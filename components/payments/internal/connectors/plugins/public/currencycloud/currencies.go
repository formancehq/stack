package currencycloud

import "github.com/formancehq/payments/internal/connectors/plugins/currency"

var (
	// c.f.: https://support.currencycloud.com/hc/en-gb/articles/7840216562972-Currency-Decimal-Places
	supportedCurrenciesWithDecimal = map[string]int{
		"AUD": currency.ISO4217Currencies["AUD"], //  Australian Dollar
		"CAD": currency.ISO4217Currencies["CAD"], //  Canadian Dollar
		"CZK": currency.ISO4217Currencies["CZK"], //  Czech Koruna
		"DKK": currency.ISO4217Currencies["DKK"], //  Danish Krone
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"HKD": currency.ISO4217Currencies["HKD"], //  Hong Kong Dollar
		"INR": currency.ISO4217Currencies["INR"], //  Indian Rupee
		"IDR": currency.ISO4217Currencies["IDR"], //  Indonesia, Rupiah
		"ILS": currency.ISO4217Currencies["ILS"], //  New Israeli Shekel
		"JPY": currency.ISO4217Currencies["JPY"], //  Japan, Yen
		"KES": currency.ISO4217Currencies["KES"], //  Kenyan Shilling
		"MYR": currency.ISO4217Currencies["MYR"], //  Malaysian Ringgit
		"MXN": currency.ISO4217Currencies["MXN"], //  Mexican Peso
		"NZD": currency.ISO4217Currencies["NZD"], //  New Zealand Dollar
		"NOK": currency.ISO4217Currencies["NOK"], //  Norwegian Krone
		"PHP": currency.ISO4217Currencies["PHP"], //  Philippine Peso
		"PLN": currency.ISO4217Currencies["PLN"], //  Poland, Zloty
		"RON": currency.ISO4217Currencies["RON"], //  Romania, New Leu
		"SAR": currency.ISO4217Currencies["SAR"], //  Saudi Riyal
		"SGD": currency.ISO4217Currencies["SGD"], //  Singapore Dollar
		"ZAR": currency.ISO4217Currencies["ZAR"], //  South Africa, Rand
		"SEK": currency.ISO4217Currencies["SEK"], //  Swedish Krona
		"CHF": currency.ISO4217Currencies["CHF"], //  Swiss Franc
		"THB": currency.ISO4217Currencies["THB"], //  Thailand, Baht
		"TRY": currency.ISO4217Currencies["TRY"], //  New Turkish Lira
		"GBP": currency.ISO4217Currencies["GBP"], //  Pound Sterling
		"AED": currency.ISO4217Currencies["AED"], //  UAE Dirham
		"USD": currency.ISO4217Currencies["USD"], //  US Dollar
		"UGX": currency.ISO4217Currencies["UGX"], //  Uganda Shilling
		"QAR": currency.ISO4217Currencies["QAR"], //  Qatari Riyal

		// Unsupported currencies
		// the following currencies are not existing in ISO 4217, so we prefer
		// not to support them for now.
		// "CNH": 2, //  Chinese Yuan

		// The following currencies have a different value in ISO 4217, so we
		// prefer to not support them for now.
		// "HUF": 0, //  Hungarian Forint
		// "KWD": 2, //  Kuwaiti Dinar
		// "OMR": 2, //  Rial Omani
		// "BHD": 2, //  Bahraini Dinar
	}
)
