package modulr

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	// c.f. https://modulr.readme.io/docs/international-payments
	supportedCurrenciesWithDecimal = map[string]int{
		"GBP": currency.ISO4217Currencies["GBP"], //  Pound Sterling
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"CZK": currency.ISO4217Currencies["CZK"], //  Czech Koruna
		"DKK": currency.ISO4217Currencies["DKK"], //  Danish Krone
		"NOK": currency.ISO4217Currencies["NOK"], //  Norwegian Krone
		"PLN": currency.ISO4217Currencies["PLN"], //  Poland, Zloty
		"SEK": currency.ISO4217Currencies["SEK"], //  Swedish Krona
		"CHF": currency.ISO4217Currencies["CHF"], //  Swiss Franc
		"USD": currency.ISO4217Currencies["USD"], //  US Dollar
		"HKD": currency.ISO4217Currencies["HKD"], //  Hong Kong Dollar
		"JPY": currency.ISO4217Currencies["JPY"], //  Japan, Yen
	}
)
