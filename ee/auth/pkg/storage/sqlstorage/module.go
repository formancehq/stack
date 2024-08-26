package sqlstorage

import (
	"context"
	"crypto/rsa"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/uptrace/bun"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/oidc"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.uber.org/fx"
)

func Module(connectionOptions bunconnect.ConnectionOptions, key *rsa.PrivateKey, debug bool, staticClients ...auth.StaticClient) fx.Option {
	return fx.Options(
		bunconnect.Module(connectionOptions, debug),
		fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logging.FromContext(ctx).Info("Migrate tables")

					return Migrate(ctx, db)
				},
			})
		}),
		fx.Supply(key),
		fx.Supply(staticClients),
		fx.Provide(fx.Annotate(New,
			fx.As(new(oidc.Storage)),
		)),
		health.ProvideHealthCheck(func(storage op.Storage) health.NamedCheck {
			return health.NewNamedCheck("Database", health.CheckFn(storage.Health))
		}),
	)
}
