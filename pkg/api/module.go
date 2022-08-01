package api

import (
	"context"
	"net/http"

	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"go.uber.org/fx"
)

func Module(issuer, addr string) fx.Option {
	return fx.Options(
		fx.Supply(Issuer(issuer)),
		fx.Provide(fx.Annotate(NewRouter, fx.As(new(http.Handler)))),
		fx.Provide(NewOpenIDProvider),
		sharedhealth.Module(),
		sharedhealth.ProvideHealthCheck(delegatedOIDCServerAvailability),
		fx.Invoke(func(lc fx.Lifecycle, handler http.Handler) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(addr, handler)
				},
			})
		}),
	)
}
