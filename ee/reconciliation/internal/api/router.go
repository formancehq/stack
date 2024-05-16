package api

import (
	"net/http"

	v1 "github.com/formancehq/reconciliation/internal/api/v1"
	backendv1 "github.com/formancehq/reconciliation/internal/api/v1/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)

func newRouter(
	backendV1 backendv1.Backend,
	serviceInfo api.ServiceInfo,
	healthController *health.HealthController,
	a auth.Auth,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	r.Get("/_healthcheck", healthController.Check)
	r.Get("/_info", api.InfoHandler(serviceInfo))

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(a))
		r.Use(otelchi.Middleware("reconciliation"))

		v1.NewRouter(backendV1, r)
	})

	return r
}
