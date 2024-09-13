package v3

import (
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

func poolsRemoveAccount(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_poolsRemoveAccount")
		defer span.End()

		id, err := uuid.Parse(poolID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		accountID, err := models.AccountIDFromString(accountID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		err = backend.PoolsRemoveAccount(ctx, id, accountID)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
