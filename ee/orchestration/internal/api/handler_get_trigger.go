package api

import (
	"database/sql"
	"net/http"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func getTrigger(backend Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trigger, err := backend.GetTrigger(r.Context(), chi.URLParam(r, "triggerID"))
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				sharedapi.NotFound(w)
				return
			default:
				sharedapi.InternalServerError(w, r, err)
			}
			return
		}

		sharedapi.Ok(w, trigger)
	}
}
