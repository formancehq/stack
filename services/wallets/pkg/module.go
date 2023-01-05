package pkg

import (
	"github.com/formancehq/wallets/pkg/api"
	"github.com/formancehq/wallets/pkg/wallet"
	"go.uber.org/fx"
)

func Module(ledgerName, chartPrefix string) fx.Option {
	return fx.Module(
		"wallets-core",
		wallet.Module(ledgerName, chartPrefix),
		api.Module(),
	)
}
