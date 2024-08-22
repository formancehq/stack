package v2

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/service"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/go-chi/chi/v5"
)

func newRouter(backend api.Backend, authenticator auth.Authenticator, debug bool) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		// Plug middleware to handle traces
		r.Use(auth.Middleware(authenticator))
		r.Use(service.OTLPMiddleware("orchestration", debug))
		r.Route("/triggers", func(r chi.Router) {
			r.Get("/", listTriggers(backend))
			r.Post("/", createTrigger(backend))
			r.Route("/{triggerID}", func(r chi.Router) {
				r.Get("/", getTrigger(backend))
				r.Delete("/", deleteTrigger(backend))
				r.Get("/occurrences", listTriggersOccurrences(backend))
				r.Post("/test", testTrigger(backend))
			})
		})
		r.Route("/workflows", func(r chi.Router) {
			r.Get("/", listWorkflows(backend))
			r.Post("/", createWorkflow(backend))
			r.Route("/{workflowId}", func(r chi.Router) {
				r.Delete("/", deleteWorkflow(backend))
				r.Get("/", readWorkflow(backend))
				r.Route("/instances", func(r chi.Router) {
					r.Post("/", runWorkflow(backend))
				})
			})
		})
		r.Route("/instances", func(r chi.Router) {
			r.Get("/", listInstances(backend))
			r.Route("/{instanceId}", func(r chi.Router) {
				r.Get("/", readInstance(backend))
				r.Post("/events", postEventToWorkflowInstance(backend))
				r.Put("/abort", abortWorkflowInstance(backend))
				r.Get("/history", readInstanceHistory(backend))
				r.Route("/stages", func(r chi.Router) {
					r.Route("/{number}", func(r chi.Router) {
						r.Get("/history", readStageHistory(backend))
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
