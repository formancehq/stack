package v2

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
)

func listTriggersOccurrences(backend api.Backend) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		triggersOccurrences, err := backend.ListTriggersOccurrences(r.Context(), chi.URLParam(r, "triggerID"))
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.Ok(w, triggersOccurrences)
	}
}
