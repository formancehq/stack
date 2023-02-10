package api

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func listWorkflows(m *workflow.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflows, err := m.ListWorkflows(r.Context())
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.Ok(w, workflows)
	}
}
