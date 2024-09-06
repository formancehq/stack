package storage

import (
	"context"

	"github.com/formancehq/stack/ee/ingester/internal/drivers"

	"github.com/formancehq/stack/ee/ingester/internal/controller"
	"github.com/formancehq/stack/ee/ingester/internal/runner"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func Module(debug bool, connectionOptions bunconnect.ConnectionOptions) fx.Option {
	return fx.Options(
		bunconnect.Module(connectionOptions, debug),
		fx.Provide(fx.Annotate(NewPostgresStore,
			fx.As(new(runner.Store)),
			fx.As(new(controller.Store)),
			fx.As(new(drivers.Store)),
		)),
		fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return errors.Wrap(NewMigrator().Up(ctx, db), "migrating database")
				},
			})
		}),
	)
}
