package internal

const (
	AdyenConnector         = "adyen"
	AtlarConnector         = "atlar"
	BankingCircleConnector = "bankingcircle"
	CurrencyCloudConnector = "currencycloud"
	ModulrConnector        = "modulr"
	StripeConnector        = "stripe"
	WiseConnector          = "wise"
	MangoPayConnector      = "mangopay"
	MoneycorpConnector     = "moneycorp"
)

var AllConnectors = []string{
	AdyenConnector,
	AtlarConnector,
	BankingCircleConnector,
	CurrencyCloudConnector,
	ModulrConnector,
	StripeConnector,
	WiseConnector,
	MangoPayConnector,
	MoneycorpConnector,
}
