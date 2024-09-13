package v3

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type bankAccountsForwardToConnectorRequest struct {
	ConnectorID string `json:"connectorID"`
}

func (f *bankAccountsForwardToConnectorRequest) Validate() error {
	if f.ConnectorID == "" {
		return errors.New("connectorID must be provided")
	}

	return nil
}

func bankAccountsForwardToConnector(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_bankAccountsForwardToConnector")
		defer span.End()

		id, err := uuid.Parse(bankAccountID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var req bankAccountsForwardToConnectorRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		err = req.Validate()
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		connectorID, err := models.ConnectorIDFromString(req.ConnectorID)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		bankAccount, err := backend.BankAccountsForwardToConnector(ctx, id, connectorID)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, bankAccount)
	}
}
