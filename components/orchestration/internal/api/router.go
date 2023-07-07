package api

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
)

func newRouter(m *workflow.Manager, info ServiceInfo, healthController *health.HealthController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	r.Get("/_healthcheck", healthController.Check)
	r.Get("/_info", getInfo(info))
	r.Group(func(r chi.Router) {
		// Plug middleware to handle traces
		r.Use(otelchi.Middleware("orchestration"))
		r.Route("/workflows", func(r chi.Router) {
			r.Get("/", listWorkflows(m))
			r.Post("/", createWorkflow(m))
			r.Route("/{workflowId}", func(r chi.Router) {
				r.Delete("/", deleteWorkflow(m))
				r.Get("/", readWorkflow(m))
				r.Route("/instances", func(r chi.Router) {
					r.Post("/", runWorkflow(m))
				})
			})
		})
		r.Route("/instances", func(r chi.Router) {
			r.Get("/", listInstances(m))
			r.Route("/{instanceId}", func(r chi.Router) {
				r.Get("/", readInstance(m))
				r.Post("/events", postEventToWorkflowInstance(m))
				r.Put("/abort", abortWorkflowInstance(m))
				r.Get("/history", readInstanceHistory(m))
				r.Route("/stages", func(r chi.Router) {
					r.Route("/{number}", func(r chi.Router) {
						r.Get("/history", readStageHistory(m))
					})
				})
			})
		})
	})

	return r
}

func workflowID(r *http.Request) string {
	return chi.URLParam(r, "workflowId")
}

func instanceID(r *http.Request) string {
	return chi.URLParam(r, "instanceId")
}
