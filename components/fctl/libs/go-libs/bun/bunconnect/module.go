package bunconnect

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func Module(connectionOptions ConnectionOptions, debug bool) fx.Option {
	return fx.Options(
		fx.Provide(func(logger logging.Logger) (*bun.DB, error) {
			hooks := make([]bun.QueryHook, 0)
			if debug {
				hooks = append(hooks, bundebug.NewQueryHook())
			}

			return OpenSQLDB(logging.ContextWithLogger(context.Background(), logger), connectionOptions, hooks...)
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
