package delegatedauth

import (
	"go.uber.org/fx"
)

func Module(issuer string) fx.Option {
	return fx.Options(
		fx.Provide(ProvideDelegatedOIDCProvider),
		fx.Supply(Issuer(issuer)),
	)
}
