package v2

import (
	"encoding/json"
	"github.com/formancehq/ledger/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"net/http"
)

func updateLedgerMetadata(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m := metadata.Metadata{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			api.BadRequest(w, "VALIDATION", errors.New("invalid format"))
			return
		}

		if err := b.UpdateLedgerMetadata(r.Context(), chi.URLParam(r, "ledger"), m); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
