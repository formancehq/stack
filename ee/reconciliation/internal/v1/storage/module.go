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
		fx.Provide(func() *bunconnect.ConnectionOptions {
			return &connectionOptions
		}),
		bunconnect.Module(connectionOptions),
		fx.Provide(func(db *bun.DB) *Storage {
			return NewStorage(db)
		}),
		fx.Invoke(func(lc fx.Lifecycle, repo *Storage) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logging.FromContext(ctx).Debug("Ping database...")

					if err := repo.Ping(); err != nil {
						return err
					}

					return nil
				},
			})
		}),
	)
}
