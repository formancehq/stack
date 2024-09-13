package v3

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func paymentsUpdateMetadata(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_paymentsUpdateMetadata")
		defer span.End()

		id, err := models.PaymentIDFromString(paymentID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var metadata map[string]string
		err = json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		err = backend.PaymentsUpdateMetadata(ctx, id, metadata)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
