package atlar

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	supportedCurrenciesWithDecimal = map[string]int{
		// For now we're only working with EUR. In the end this will be a bit more complex than
		// it is with other connectors, because in the end not only atlar, but also the banks
		// connected to Atlar will have to support the currencies.
		"EUR": currency.ISO4217Currencies["EUR"],
	}
)
