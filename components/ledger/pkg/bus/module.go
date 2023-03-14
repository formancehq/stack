package bus

import (
	"github.com/formancehq/ledger/pkg/query"
	"go.uber.org/fx"
)

func LedgerMonitorModule() fx.Option {
	return fx.Decorate(fx.Annotate(newLedgerMonitor, fx.As(new(query.Monitor))))
}
