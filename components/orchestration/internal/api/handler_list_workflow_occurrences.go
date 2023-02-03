package api

import (
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/orchestration/internal/workflow"
)

func listInstances(m *workflow.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runs, err := m.ListInstances(r.Context(), r.URL.Query().Get("workflowID"))
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
		api.Ok(w, runs)
	}
}
