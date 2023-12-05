package postgres

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/logging"
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
				OnStop: func(ctx context.Context) error {
					logging.FromContext(ctx).Info("Closing database...")
					defer func() {
						logging.FromContext(ctx).Info("Database closed.")
					}()

					return s.Close(ctx)
				},
			})
		}),
	)
}
