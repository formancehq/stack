package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
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

func reverseTransferInitiation(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &service.ReverseTransferInitiationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := payload.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		transferID, ok := mux.Vars(r)["transferID"]
		if !ok {
			api.BadRequest(w, ErrInvalidID, errors.New("missing transferID"))
			return
		}

		transferReversal, err := b.GetService().ReverseTransferInitiation(r.Context(), transferID, payload)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &reverseTransferInitiationResponse{
			ID:                  transferID,
			TranferInitiationID: transferID,
			CreatedAt:           transferReversal.CreatedAt,
			UpdatedAt:           transferReversal.UpdatedAt,
			Reference:           transferReversal.ID.Reference,
			Description:         transferReversal.Description,
			Amount:              transferReversal.Amount,
			Asset:               transferReversal.Asset.String(),
			Status:              transferReversal.Status.String(),
			Error:               transferReversal.Error,
			Metadata:            transferReversal.Metadata,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[reverseTransferInitiationResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}
