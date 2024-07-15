package atlar

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	supportedCurrenciesWithDecimal = map[string]int{
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"DKK": currency.ISO4217Currencies["DKK"],
	}
)
