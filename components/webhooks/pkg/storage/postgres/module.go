package postgres

import (
	"github.com/formancehq/webhooks/pkg/storage"
	"go.uber.org/fx"
)

func NewModule(dsn string) fx.Option {
	return fx.Options(
		fx.Provide(func() (storage.Store, error) {
			return NewStore(dsn)
		}),
		fx.Invoke(func(lc fx.Lifecycle, s storage.Store) {
			lc.Append(fx.Hook{
				OnStop: s.Close,
			})
		}),
	)
}
