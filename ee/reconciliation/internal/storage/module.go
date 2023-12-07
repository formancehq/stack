package storage

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
	"go.uber.org/fx"
)

const dbName = "paymentsDB"

func Module(uri string, debug bool, output io.Writer) fx.Option {
	return fx.Options(
		fx.Provide(func() (*pgx.ConnConfig, error) {
			config, err := pgx.ParseConfig(uri)
			if err != nil {
				return nil, fmt.Errorf("parsing config: %w", err)
			}

			return config, nil
		}),

		fx.Provide(func(config *pgx.ConnConfig) *sql.DB {
			return stdlib.OpenDB(*config)
		}),
		fx.Provide(func(client *sql.DB) *bun.DB {
			db := bun.NewDB(client, pgdialect.New())

			db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(dbName)))

			if debug {
				db.AddQueryHook(bundebug.NewQueryHook(
					bundebug.WithWriter(output),
				))
			}

			return db
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
