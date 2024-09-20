package v2

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/orchestration/internal/api"
	"github.com/pkg/errors"
)

func getTrigger(backend api.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trigger, err := backend.GetTrigger(r.Context(), chi.URLParam(r, "triggerID"))
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				sharedapi.NotFound(w, err)
				return
			default:
				sharedapi.InternalServerError(w, r, err)
			}
			return
		}

		sharedapi.Ok(w, trigger)
	}
}
