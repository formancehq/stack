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

type transferInitiationResponse struct {
	ID                   string            `json:"id"`
	CreatedAt            time.Time         `json:"createdAt"`
	UpdatedAt            time.Time         `json:"updatedAt"`
	ScheduledAt          time.Time         `json:"scheduledAt"`
	Description          string            `json:"description"`
	SourceAccountID      string            `json:"sourceAccountID"`
	DestinationAccountID string            `json:"destinationAccountID"`
	ConnectorID          string            `json:"connectorID"`
	Type                 string            `json:"type"`
	Amount               *big.Int          `json:"amount"`
	Asset                string            `json:"asset"`
	Status               string            `json:"status"`
	Error                string            `json:"error"`
	Metadata             map[string]string `json:"metadata"`
}

func createTransferInitiationHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload := &service.CreateTransferInitiationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := payload.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		tf, err := b.GetService().CreateTransferInitiation(r.Context(), payload)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &transferInitiationResponse{
			ID:                   tf.ID.String(),
			CreatedAt:            tf.CreatedAt,
			UpdatedAt:            tf.UpdatedAt,
			ScheduledAt:          tf.ScheduledAt,
			Description:          tf.Description,
			SourceAccountID:      tf.SourceAccountID.String(),
			DestinationAccountID: tf.DestinationAccountID.String(),
			ConnectorID:          tf.ConnectorID.String(),
			Type:                 tf.Type.String(),
			Amount:               tf.Amount,
			Asset:                tf.Asset.String(),
			Status:               tf.Status.String(),
			Error:                tf.Error,
			Metadata:             tf.Metadata,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[transferInitiationResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func updateTransferInitiationStatusHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &service.UpdateTransferInitiationStatusRequest{}
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

		if err := b.GetService().UpdateTransferInitiationStatus(r.Context(), transferID, payload); err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func retryTransferInitiationHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferID, ok := mux.Vars(r)["transferID"]
		if !ok {
			api.BadRequest(w, ErrInvalidID, errors.New("missing transferID"))
			return
		}

		if err := b.GetService().RetryTransferInitiation(r.Context(), transferID); err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteTransferInitiationHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferID, ok := mux.Vars(r)["transferID"]
		if !ok {
			api.BadRequest(w, ErrInvalidID, errors.New("missing transferID"))
			return
		}

		if err := b.GetService().DeleteTransferInitiation(r.Context(), transferID); err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
