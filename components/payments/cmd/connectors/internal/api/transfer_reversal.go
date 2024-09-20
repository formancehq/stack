package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/internal/otel"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func reverseTransferInitiationHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "reverseTransferInitiationHandler")
		defer span.End()

		payload := &service.ReverseTransferInitiationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := payload.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		transferID, ok := mux.Vars(r)["transferID"]
		if !ok {
			otel.RecordError(span, errors.New("missing transferID"))
			api.BadRequest(w, ErrInvalidID, errors.New("missing transferID"))
			return
		}

		_, err := b.GetService().ReverseTransferInitiation(ctx, transferID, payload)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
