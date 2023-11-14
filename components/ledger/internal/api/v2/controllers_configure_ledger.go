package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/ledger/internal/api/backend"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func configureLedger(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		configuration := backend.Configuration{}
		if err := json.NewDecoder(r.Body).Decode(&configuration); err != nil {
			sharedapi.BadRequest(w, ErrValidation, err)
			return
		}
		if err := b.ConfigureLedger(r.Context(), chi.URLParam(r, "ledger"), configuration); err != nil {
			switch {
			case errors.Is(err, backend.ErrAlreadyConfigured):
				sharedapi.BadRequest(w, ErrValidation, err)
			default:
				sharedapi.InternalServerError(w, r, err)
			}
			return
		}
		sharedapi.NoContent(w)
	}
}
