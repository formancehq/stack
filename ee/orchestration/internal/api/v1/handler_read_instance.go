package v1

import (
	"net/http"

	api2 "github.com/formancehq/orchestration/internal/api"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func readInstance(backend api2.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflows, err := backend.GetInstance(r.Context(), instanceID(r))
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.Ok(w, workflows)
	}
}
