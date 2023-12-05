package api

import (
	"net/http"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
)

func listTriggersOccurrences(backend Backend) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		triggersOccurrences, err := backend.ListTriggersOccurrences(r.Context(), chi.URLParam(r, "triggerID"))
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.Ok(w, triggersOccurrences)
	}
}
