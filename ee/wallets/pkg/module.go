package wallet

import (
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go/v2"

	"go.uber.org/fx"
)

func Module(ledgerURL string, ledgerName, chartPrefix string) fx.Option {
	return fx.Module(
		"wallet",
		fx.Provide(fx.Annotate(func(httpClient *http.Client) *DefaultLedger {
			sdk := sdk.New(
				sdk.WithClient(httpClient),
				sdk.WithServerURL(ledgerURL),
			)
			return NewDefaultLedger(sdk)
		}, fx.As(new(Ledger)))),
		fx.Provide(func() *Chart {
			return NewChart(chartPrefix)
		}),
		fx.Provide(func(ledger Ledger, chart *Chart) *Manager {
			return NewManager(ledgerName, ledger, chart)
		}),
		fx.Invoke(func(manager *Manager, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: manager.Init,
			})
		}),
	)
}
