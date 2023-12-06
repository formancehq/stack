package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func readInstance(backend Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflows, err := backend.GetInstance(r.Context(), instanceID(r))
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.Ok(w, workflows)
	}
}
