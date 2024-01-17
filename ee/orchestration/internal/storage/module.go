package storage

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func NewModule(options bunconnect.ConnectionOptions) fx.Option {
	return fx.Options(
		fx.Provide(func() (*bun.DB, error) {
			return bunconnect.OpenSQLDB(options)
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
