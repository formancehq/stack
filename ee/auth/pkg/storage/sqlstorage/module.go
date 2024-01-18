package sqlstorage

import (
	"crypto/rsa"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/oidc"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.uber.org/fx"
)

func Module(connectionOptions bunconnect.ConnectionOptions, key *rsa.PrivateKey, staticClients ...auth.StaticClient) fx.Option {
	return fx.Options(
		bunModule(connectionOptions),
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
