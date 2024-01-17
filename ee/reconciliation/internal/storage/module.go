package storage

import (
	"context"
	"database/sql"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func Module(connectionOptions bunconnect.ConnectionOptions) fx.Option {
	return fx.Options(
		fx.Provide(func(client *sql.DB) (*bun.DB, error) {
			return bunconnect.OpenSQLDB(connectionOptions)
		}),
		fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					logging.FromContext(ctx).Infof("Closing database...")
					defer func() {
						logging.FromContext(ctx).Infof("Database closed.")
					}()
					return db.Close()
				},
			})
		}),
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
