package postgres

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/pkg/storage"
	"go.uber.org/fx"
)

func NewModule(connectionOptions bunconnect.ConnectionOptions) fx.Option {
	return fx.Options(
		fx.Provide(func() (storage.Store, error) {
			db, err := bunconnect.OpenSQLDB(connectionOptions)
			if err != nil {
				return nil, err
			}
			return NewStore(db)
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
