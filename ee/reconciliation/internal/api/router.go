package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/service"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/health"
	"github.com/formancehq/reconciliation/internal/api/backend"
)

func newRouter(
	b backend.Backend,
	serviceInfo api.ServiceInfo,
	authenticator auth.Authenticator,
	healthController *health.HealthController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	r.Get("/_healthcheck", healthController.Check)
	r.Get("/_info", api.InfoHandler(serviceInfo))

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(authenticator))
		r.Use(service.OTLPMiddleware("reconciliation", serviceInfo.Debug))

		r.Get("/reconciliations/{reconciliationID}", getReconciliationHandler(b))
		r.Get("/reconciliations", listReconciliationsHandler(b))

		r.Post("/policies", createPolicyHandler(b))
		r.Get("/policies", listPoliciesHandler(b))
		r.Delete("/policies/{policyID}", deletePolicyHandler(b))
		r.Get("/policies/{policyID}", getPolicyHandler(b))
		r.Post("/policies/{policyID}/reconciliation", reconciliationHandler(b))
	})

	return r
}
