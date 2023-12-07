package wallet

import (
	"net/http"

	"go.uber.org/fx"
)

func Module(ledgerURL string, ledgerName, chartPrefix string) fx.Option {
	return fx.Module(
		"wallet",
		fx.Provide(fx.Annotate(func(httpClient *http.Client) *DefaultLedger {
			return NewDefaultLedger(httpClient, ledgerURL)
		}, fx.As(new(Ledger)))),
		fx.Provide(func() *Chart {
			return NewChart(chartPrefix)
		}),
		fx.Provide(func(ledger Ledger, chart *Chart) *Manager {
			return NewManager(ledgerName, ledger, chart)
		}),
	)
}
