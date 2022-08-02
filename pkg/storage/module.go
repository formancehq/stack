package storage

import (
	"crypto/rsa"

	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
)

func Module(uri string, key *rsa.PrivateKey) fx.Option {
	return fx.Options(
		gormModule(uri),
		fx.Supply(key),
		fx.Provide(fx.Annotate(New, fx.As(new(Storage)), fx.As(new(op.Storage)))),
		sharedhealth.ProvideHealthCheck(func(storage op.Storage) sharedhealth.NamedCheck {
			return sharedhealth.NewNamedCheck("Database", sharedhealth.CheckFn(storage.Health))
		}),
	)
}
