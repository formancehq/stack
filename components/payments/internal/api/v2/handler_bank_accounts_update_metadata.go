package v2

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type bankAccountsUpdateMetadataRequest struct {
	Metadata map[string]string `json:"metadata"`
}

func (u *bankAccountsUpdateMetadataRequest) Validate() error {
	if len(u.Metadata) == 0 {
		return errors.New("metadata must be provided")
	}

	return nil
}

func bankAccountsUpdateMetadata(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_bankAccountsUpdateMetadata")
		defer span.End()

		id, err := uuid.Parse(bankAccountID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var req bankAccountsUpdateMetadataRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		err = req.Validate()
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		err = backend.BankAccountsUpdateMetadata(ctx, id, req.Metadata)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
