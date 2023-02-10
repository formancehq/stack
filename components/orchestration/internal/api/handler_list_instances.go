package api

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func listInstances(m *workflow.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runs, err := m.ListInstances(r.Context(), r.URL.Query().Get("workflowID"), r.URL.Query().Get("running") == "true")
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
		api.Ok(w, runs)
	}
}
