package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/triggers"
)

func listTriggersOccurrences(backend api.Backend) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		triggersOccurrences, err := backend.ListTriggersOccurrences(r.Context(), triggers.ListTriggersOccurrencesQuery{
			Options: triggers.ListTriggersOccurrencesOptions{
				TriggerID: chi.URLParam(r, "triggerID"),
			},
		})
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.Ok(w, triggersOccurrences.Data)
	}
}
