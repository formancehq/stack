package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func readWorkflow(backend Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflow, err := backend.ReadWorkflow(r.Context(), workflowID(r))
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.Ok(w, workflow)
	}
}
