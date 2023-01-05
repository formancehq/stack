package wallet

import (
	"github.com/formancehq/wallets/pkg/core"
	"go.uber.org/fx"
)

func Module(ledgerName, chartPrefix string) fx.Option {
	return fx.Module(
		"wallet",
		fx.Provide(fx.Annotate(NewDefaultLedger, fx.As(new(Ledger)))),
		fx.Provide(func() *core.Chart {
			return core.NewChart(chartPrefix)
		}),
		fx.Provide(func(ledger Ledger, chart *core.Chart) *FundingService {
			return NewFundingService(ledgerName, ledger, chart)
		}),
		fx.Provide(func(ledger Ledger, chart *core.Chart) *Repository {
			return NewRepository(ledgerName, ledger, chart)
		}),
	)
}
