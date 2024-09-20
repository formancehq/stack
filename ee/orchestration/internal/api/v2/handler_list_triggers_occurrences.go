package v2

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/orchestration/internal/triggers"

	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/orchestration/internal/api"
)

func listTriggersOccurrences(backend api.Backend) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		query, err := bunpaginate.Extract[triggers.ListTriggersOccurrencesQuery](r, func() (*triggers.ListTriggersOccurrencesQuery, error) {
			pageSize, err := bunpaginate.GetPageSize(r)
			if err != nil {
				return nil, err
			}
			return &triggers.ListTriggersOccurrencesQuery{
				PageSize: pageSize,
				Options: triggers.ListTriggersOccurrencesOptions{
					TriggerID: chi.URLParam(r, "triggerID"),
				},
			}, nil
		})
		if err != nil {
			sharedapi.BadRequest(w, "VALIDATION", err)
			return
		}

		triggersOccurrences, err := backend.ListTriggersOccurrences(r.Context(), *query)
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.RenderCursor(w, *triggersOccurrences)
	}
}
