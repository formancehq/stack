package api

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func deleteWorkflow(m *workflow.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := m.DeleteWorkflow(r.Context(), workflowID(r))
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
