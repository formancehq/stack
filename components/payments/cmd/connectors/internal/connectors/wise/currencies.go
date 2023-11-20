package wise

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	// c.f. https://wise.com/help/articles/2897238/which-currencies-can-i-add-keep-and-receive-in-my-wise-account
	supportedCurrenciesWithDecimal = map[string]int{
		"AUD": currency.ISO4217Currencies["AUD"], //  Australian dollar
		"BGN": currency.ISO4217Currencies["BGN"], //  Bulgarian lev
		"BRL": currency.ISO4217Currencies["BRL"], //  Brazilian real
		"CAD": currency.ISO4217Currencies["CAD"], //  Canadian dollar
		"CNY": currency.ISO4217Currencies["CNY"], //  Chinese yuan
		"CHF": currency.ISO4217Currencies["CHF"], //  Swiss franc
		"CZK": currency.ISO4217Currencies["CZK"], //  Czech koruna
		"DKK": currency.ISO4217Currencies["DKK"], //  Danish krone
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"GBP": currency.ISO4217Currencies["GBP"], //  Pound sterling
		"IDR": currency.ISO4217Currencies["IDR"], //  Indonesian rupiah
		"JPY": currency.ISO4217Currencies["JPY"], //  Japanese yen
		"MYR": currency.ISO4217Currencies["MYR"], //  Malaysian ringgit
		"NOK": currency.ISO4217Currencies["NOK"], //  Norwegian krone
		"NZD": currency.ISO4217Currencies["NZD"], //  New Zealand dollar
		"PLN": currency.ISO4217Currencies["PLN"], //  Polish złoty
		"RON": currency.ISO4217Currencies["RON"], //  Romanian leu
		"SEK": currency.ISO4217Currencies["SEK"], //  Swedish krona/kronor
		"SGD": currency.ISO4217Currencies["SGD"], //  Singapore dollar
		"TRY": currency.ISO4217Currencies["TRY"], //  Turkish lira
		"USD": currency.ISO4217Currencies["USD"], //  United States dollar

		// Unsupported currencies
		// "HUF": currency.ISO4217Currencies["HUF"], //  Hungarian forint
	}
)
