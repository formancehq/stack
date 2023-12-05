package routes

import (
	"net/http"

	"github.com/formancehq/stack/components/stargate/internal/server/http/controllers"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
)

func NewRouter(
	logger logging.Logger,
	healthController *health.HealthController,
	stargateController *controllers.StargateController,
) chi.Router {
	router := chi.NewMux()

	router.Use(
		cors.New(cors.Options{
			AllowOriginFunc: func(r *http.Request, origin string) bool {
				return true
			},
			AllowCredentials: true,
		}).Handler,
		func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if logger != nil {
					r = r.WithContext(
						logging.ContextWithLogger(r.Context(), logger),
					)
				}
				handler.ServeHTTP(w, r)
			})
		},
		middleware.Recoverer,
	)

	router.Get("/_healthcheck", healthController.Check)

	router.Group(func(router chi.Router) {
		router.Use(otelchi.Middleware("stargate"))
		router.Get("/_info", stargateController.GetInfo)

		router.HandleFunc("/*", stargateController.HandleCalls)
	})

	return router
}
