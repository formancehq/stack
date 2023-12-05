package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func listInstances(backend Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runs, err := backend.ListInstances(r.Context(), r.URL.Query().Get("workflowID"), r.URL.Query().Get("running") == "true")
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
		api.Ok(w, runs)
	}
}
