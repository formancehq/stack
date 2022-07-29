package storage

import (
	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
)

func Module(uri string) fx.Option {
	return fx.Options(
		gormModule(uri),
		fx.Provide(fx.Annotate(New, fx.As(new(Storage)), fx.As(new(op.Storage)))),
	)
}
