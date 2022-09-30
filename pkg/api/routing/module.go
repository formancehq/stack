package routing

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/fx"
)

func CreateRootRouter(baseUrl *url.URL, healthController *sharedhealth.HealthController) (*mux.Router, *mux.Router) {
	rootRouter := mux.NewRouter()
	rootRouter.Use(otelmux.Middleware("auth"))
	rootRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	rootRouter.Path("/_healthcheck").HandlerFunc(healthController.Check)

	subRouter := rootRouter
	if baseUrl.Path != "/" {
		subRouter = subRouter.PathPrefix(baseUrl.Path).Subrouter()
		subRouter.Path("/_healthcheck").HandlerFunc(healthController.Check)
	}

	return rootRouter, subRouter
}

func Module(addr string, baseUrl *url.URL) fx.Option {
	return fx.Options(
		sharedhealth.Module(),
		fx.Provide(fx.Annotate(func(healthController *sharedhealth.HealthController) (*mux.Router, *mux.Router) {
			return CreateRootRouter(baseUrl, healthController)
		}, fx.ResultTags(`name:"rootRouter"`, `name:"prefixedRouter"`))),
		fx.Invoke(fx.Annotate(func(lc fx.Lifecycle, router *mux.Router, healthController *sharedhealth.HealthController) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(ctx, addr, router)
				},
			})
		}, fx.ParamTags(``, `name:"rootRouter"`))),
	)
}
