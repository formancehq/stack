package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/service"

	"github.com/formancehq/auth/pkg/api/authorization"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/health"
	"github.com/formancehq/go-libs/httpserver"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.uber.org/fx"
)

func CreateRootRouter(o op.OpenIDProvider, issuer string, debug bool) chi.Router {
	rootRouter := chi.NewRouter()
	rootRouter.Use(service.OTLPMiddleware("auth", debug))
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
	rootRouter.Use(authorization.Middleware(o))
	return rootRouter
}

func addInfoRoute(router chi.Router, serviceInfo api.ServiceInfo) {
	router.Get("/_info", api.InfoHandler(serviceInfo))
}

func Module(addr, issuer string, serviceInfo api.ServiceInfo, debug bool) fx.Option {
	return fx.Options(
		health.Module(),
		fx.Supply(serviceInfo),
		fx.Provide(func(o op.OpenIDProvider) chi.Router {
			return CreateRootRouter(o, issuer, debug)
		}),
		fx.Invoke(
			addInfoRoute,
			addClientRoutes,
			addUserRoutes,
		),
		fx.Invoke(func(lc fx.Lifecycle, r chi.Router, healthController *health.HealthController, o op.OpenIDProvider) {
			finalRouter := chi.NewRouter()
			finalRouter.Get("/_healthcheck", healthController.Check)
			finalRouter.Mount("/", r)

			lc.Append(httpserver.NewHook(finalRouter, httpserver.WithAddress(addr)))
		}),
	)
}
