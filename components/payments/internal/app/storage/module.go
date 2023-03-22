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

func Module(uri, configEncryptionKey string, debug bool, output io.Writer) fx.Option {
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

		fx.Provide(func(client *sql.DB) *Storage {
			db := bun.NewDB(client, pgdialect.New())

			db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(dbName)))

			if debug {
				db.AddQueryHook(bundebug.NewQueryHook(
					bundebug.WithWriter(output),
				))
			}

			return newStorage(db, configEncryptionKey)
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
