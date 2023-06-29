package storage

import (
	"context"
	"database/sql"
	"io"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/fx"
)

func LoadDB(dsn string, debug bool, output io.Writer) *bun.DB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	if debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(debug),
			bundebug.FromEnv("BUNDEBUG"),
			bundebug.WithWriter(output),
		))
	}
	return db
}

func NewModule(dsn string, debug bool, output io.Writer) fx.Option {
	return fx.Options(
		fx.Provide(func() *bun.DB {
			return LoadDB(dsn, debug, output)
		}),
		fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return db.Close()
				},
			})
		}),
	)
}
