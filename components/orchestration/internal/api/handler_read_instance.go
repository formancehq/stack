package api

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func readInstance(m *workflow.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflows, err := m.GetInstance(r.Context(), instanceID(r))
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.Ok(w, workflows)
	}
}
