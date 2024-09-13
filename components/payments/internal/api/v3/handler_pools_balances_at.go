package v3

import (
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func poolsBalancesAt(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_poolsBalancesAt")
		defer span.End()

		id, err := uuid.Parse(poolID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		atTime := r.URL.Query().Get("at")
		if atTime == "" {
			otel.RecordError(span, errors.New("missing atTime"))
			api.BadRequest(w, ErrValidation, errors.New("missing atTime"))
			return
		}

		at, err := time.Parse(time.RFC3339, atTime)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, errors.Wrap(err, "invalid atTime"))
			return
		}

		balances, err := backend.PoolsBalancesAt(ctx, id, at)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, balances)
	}
}
