package storage

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/logging"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func Module(connectionOptions bunconnect.ConnectionOptions, configEncryptionKey string, debug bool) fx.Option {
	return fx.Options(
		fx.Supply(&connectionOptions),
		bunconnect.Module(connectionOptions, debug),
		fx.Provide(func(db *bun.DB) *Storage {
			return NewStorage(db, configEncryptionKey)
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
