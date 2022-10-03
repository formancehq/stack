package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/fx"
)

func CreateRootRouter(healthController *sharedhealth.HealthController) *mux.Router {
	rootRouter := mux.NewRouter()
	rootRouter.Use(otelmux.Middleware("auth"))
	rootRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	rootRouter.Path("/_healthcheck").HandlerFunc(healthController.Check)

	return rootRouter
}

func Module(addr string) fx.Option {
	return fx.Options(
		sharedhealth.ProvideHealthCheck(delegatedOIDCServerAvailability),
		sharedhealth.Module(),
		fx.Provide(func(healthController *sharedhealth.HealthController) *mux.Router {
			return CreateRootRouter(healthController)
		}),
		fx.Invoke(func(lc fx.Lifecycle, router *mux.Router, healthController *sharedhealth.HealthController) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(ctx, addr, router)
				},
			})
		}),
		fx.Invoke(
			addClientRoutes,
			addScopeRoutes,
			addUserRoutes,
		),
	)
}
