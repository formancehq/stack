package v1

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/orchestration/internal/api"
	"github.com/pkg/errors"
)

func deleteTrigger(backend api.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := backend.DeleteTrigger(r.Context(), chi.URLParam(r, "triggerID")); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				sharedapi.NotFound(w, err)
			default:
				sharedapi.InternalServerError(w, r, err)
			}
			return
		}

		sharedapi.NoContent(w)
	}
}
