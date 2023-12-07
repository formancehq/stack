package api

import (
	"database/sql"
	"net/http"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func deleteTrigger(backend Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := backend.DeleteTrigger(r.Context(), chi.URLParam(r, "triggerID")); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				sharedapi.NotFound(w)
			default:
				sharedapi.InternalServerError(w, r, err)
			}
			return
		}

		sharedapi.NoContent(w)
	}
}
