package sqlstorage

import (
	"crypto/rsa"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/oidc"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
)

func Module(uri string, debug bool, key *rsa.PrivateKey, staticClients []auth.StaticClient) fx.Option {
	return fx.Options(
		gormModule(uri, debug),
		fx.Supply(key),
		fx.Supply(staticClients),
		fx.Provide(fx.Annotate(New,
			fx.As(new(oidc.Storage)),
		)),
		sharedhealth.ProvideHealthCheck(func(storage op.Storage) sharedhealth.NamedCheck {
			return sharedhealth.NewNamedCheck("Database", sharedhealth.CheckFn(storage.Health))
		}),
	)
}
