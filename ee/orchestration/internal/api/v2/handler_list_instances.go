package v2

import (
	"net/http"

	api2 "github.com/formancehq/orchestration/internal/api"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func listInstances(backend api2.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runs, err := backend.ListInstances(r.Context(), r.URL.Query().Get("workflowID"), r.URL.Query().Get("running") == "true")
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
		api.Ok(w, runs)
	}
}
