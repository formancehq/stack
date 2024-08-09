package fawry

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	supportedCurrenciesWithDecimal = map[string]int{
		"EGP": currency.ISO4217Currencies["EGP"],
	}
)
