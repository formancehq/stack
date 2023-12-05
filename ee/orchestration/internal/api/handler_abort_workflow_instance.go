package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func abortWorkflowInstance(backend Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := backend.AbortRun(r.Context(), instanceID(r)); err != nil {
			api.InternalServerError(w, r, err)
			return
		}
		api.NoContent(w)
	}
}
