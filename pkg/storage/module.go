package storage

import (
	"crypto/rsa"

	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
)

func Module(uri string, key *rsa.PrivateKey) fx.Option {
	return fx.Options(
		gormModule(uri),
		fx.Supply(key),
		fx.Provide(fx.Annotate(New, fx.As(new(Storage)), fx.As(new(op.Storage)))),
	)
}
