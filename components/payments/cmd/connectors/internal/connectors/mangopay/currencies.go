package mangopay

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	// c.f. https://mangopay.com/docs/api-basics/data-formats
	supportedCurrenciesWithDecimal = map[string]int{
		"AED": currency.ISO4217Currencies["AED"], //  UAE Dirham
		"AUD": currency.ISO4217Currencies["AUD"], //  Australian Dollar
		"CAD": currency.ISO4217Currencies["CAD"], //  Canadian Dollar
		"CHF": currency.ISO4217Currencies["CHF"], //  Swiss Franc
		"CZK": currency.ISO4217Currencies["CZK"], //  Czech Koruna
		"DKK": currency.ISO4217Currencies["DKK"], //  Danish Krone
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"GBP": currency.ISO4217Currencies["GBP"], //  Pound Sterling
		"HKD": currency.ISO4217Currencies["HKD"], //  Hong Kong Dollar
		"JPY": currency.ISO4217Currencies["JPY"], //  Japan, Yen
		"NOK": currency.ISO4217Currencies["NOK"], //  Norwegian Krone
		"PLN": currency.ISO4217Currencies["PLN"], //  Poland, Zloty
		"SEK": currency.ISO4217Currencies["SEK"], //  Swedish Krona
		"USD": currency.ISO4217Currencies["USD"], //  US Dollar
		"ZAR": currency.ISO4217Currencies["ZAR"], //  South Africa, Rand
	}
)
