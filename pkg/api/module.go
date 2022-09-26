package api

import (
	"context"
	"net/url"

	"github.com/gorilla/mux"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"go.uber.org/fx"
)

func Module(addr string, baseUrl *url.URL) fx.Option {
	return fx.Options(
		sharedhealth.Module(),
		sharedhealth.ProvideHealthCheck(delegatedOIDCServerAvailability),
		fx.Provide(func() *mux.Router {
			return NewRouter(baseUrl)
		}),
		fx.Invoke(
			addClientRoutes,
			addScopeRoutes,
			addUserRoutes,
		),
		fx.Invoke(func(lc fx.Lifecycle, router *mux.Router) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(addr, router)
				},
			})
		}),
	)
}
