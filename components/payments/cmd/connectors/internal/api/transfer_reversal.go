package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type reverseTransferInitiationResponse struct {
	ID                  string            `json:"id"`
	TranferInitiationID string            `json:"transferInitiationID"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
	Reference           string            `json:"reference"`
	Description         string            `json:"description"`
	Amount              *big.Int          `json:"amount"`
	Asset               string            `json:"asset"`
	Status              string            `json:"status"`
	Error               string            `json:"error"`
	Metadata            map[string]string `json:"metadata"`
}

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
