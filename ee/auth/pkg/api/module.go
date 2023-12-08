package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/fx"
)

func CreateRootRouter(issuer string) *mux.Router {
	rootRouter := mux.NewRouter()
	rootRouter.Use(otelmux.Middleware("auth"))
	rootRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	rootRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r.WithContext(
				op.ContextWithIssuer(r.Context(), issuer),
			))
		})
	})
	return rootRouter
}

func addInfoRoute(router *mux.Router, serviceInfo api.ServiceInfo) {
	router.Path("/_info").Methods(http.MethodGet).HandlerFunc(api.InfoHandler(serviceInfo))
}

func Module(addr, issuer string, serviceInfo api.ServiceInfo) fx.Option {
	return fx.Options(
		health.Module(),
		fx.Supply(serviceInfo),
		fx.Provide(func() *mux.Router {
			return CreateRootRouter(issuer)
		}),
		fx.Invoke(func(lc fx.Lifecycle, r *mux.Router, healthController *health.HealthController) {
			finalRouter := mux.NewRouter()
			finalRouter.Path("/_healthcheck").HandlerFunc(healthController.Check)
			finalRouter.PathPrefix("/").Handler(r)

			lc.Append(httpserver.NewHook(finalRouter, httpserver.WithAddress(addr)))
		}),
		fx.Invoke(
			addInfoRoute,
			addClientRoutes,
			addUserRoutes,
		),
	)
}
