package storage

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func Module(connectionOptions bunconnect.ConnectionOptions) fx.Option {
	return fx.Options(
		bunconnect.Module(connectionOptions),
		fx.Provide(func(db *bun.DB) *Storage {
			return NewStorage(db)
		}),
		fx.Invoke(func(lc fx.Lifecycle, repo *Storage) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logging.FromContext(ctx).Debug("Ping database...")

					// TODO: Check migrations state and panic if migrations are not applied

					return nil
				},
			})
		}),
	)
}
