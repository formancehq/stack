package ledger

import (
	"github.com/formancehq/ledger/pkg/ledger/lock"
	"github.com/formancehq/ledger/pkg/storage"
	"go.uber.org/fx"
)

func Module(allowPastTimestamp bool) fx.Option {
	return fx.Options(
		lock.Module(),
		fx.Provide(func(
			storageDriver storage.Driver,
			locker lock.Locker,
		) *Resolver {
			return NewResolver(storageDriver, locker, allowPastTimestamp)
		}),
	)
}
