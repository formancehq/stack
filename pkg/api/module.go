package api

import (
	"context"
	"net/http"

	sharedhealth "github.com/formancehq/go-libs/sharedhealth/pkg"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/fx"
)

func CreateRootRouter() *mux.Router {
	rootRouter := mux.NewRouter()
	rootRouter.Use(otelmux.Middleware("auth"))
	rootRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	return rootRouter
}

func Module(addr string) fx.Option {
	return fx.Options(
		sharedhealth.ProvideHealthCheck(delegatedOIDCServerAvailability),
		sharedhealth.Module(),
		fx.Provide(CreateRootRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *mux.Router, healthController *sharedhealth.HealthController) {
			finalRouter := mux.NewRouter()
			finalRouter.Path("/_healthcheck").HandlerFunc(healthController.Check)
			finalRouter.PathPrefix("/").Handler(r)
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(ctx, addr, finalRouter)
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
