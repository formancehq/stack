package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/connectors/internal/integration"
	"github.com/formancehq/payments/cmd/connectors/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type paymentHandler func(ctx context.Context, transfer *models.TransferInitiation) error

type transferInitiationResponse struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	ScheduledAt          time.Time `json:"scheduledAt"`
	Description          string    `json:"description"`
	SourceAccountID      string    `json:"sourceAccountID"`
	DestinationAccountID string    `json:"destinationAccountID"`
	ConnectorID          string    `json:"connectorID"`
	Type                 string    `json:"type"`
	Amount               *big.Int  `json:"amount"`
	Asset                string    `json:"asset"`
	Status               string    `json:"status"`
	Error                string    `json:"error"`
}

type createTransferInitiationRequest struct {
	Reference            string    `json:"reference"`
	ScheduledAt          time.Time `json:"scheduledAt"`
	Description          string    `json:"description"`
	SourceAccountID      string    `json:"sourceAccountID"`
	DestinationAccountID string    `json:"destinationAccountID"`
	ConnectorID          string    `json:"connectorID"`
	Provider             string    `json:"provider"`
	Type                 string    `json:"type"`
	Amount               *big.Int  `json:"amount"`
	Asset                string    `json:"asset"`
	Validated            bool      `json:"validated"`
}

func (r *createTransferInitiationRequest) Validate(repo createTransferInitiationRepository) error {
	if r.Reference == "" {
		return errors.New("uniqueRequestId is required")
	}

	if r.Description == "" {
		return errors.New("description is required")
	}

	if r.SourceAccountID != "" {
		_, err := models.AccountIDFromString(r.SourceAccountID)
		if err != nil {
			return err
		}
	}

	_, err := models.AccountIDFromString(r.DestinationAccountID)
	if err != nil {
		return err
	}

	_, err = models.TransferInitiationTypeFromString(r.Type)
	if err != nil {
		return err
	}

	if r.Amount == nil {
		return errors.New("amount is required")
	}

	if r.Asset == "" {
		return errors.New("asset is required")
	}

	return nil
}

type createTransferInitiationRepository interface {
	CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error)
	ListConnectorsByProvider(ctx context.Context, provider models.ConnectorProvider) ([]*models.Connector, error)
	IsInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error)
}

func createTransferInitiationHandler(
	repo createTransferInitiationRepository,
	publisher message.Publisher,
	paymentHandlers map[models.ConnectorProvider]paymentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload := &createTransferInitiationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := payload.Validate(repo); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		status := models.TransferInitiationStatusWaitingForValidation
		if payload.Validated {
			status = models.TransferInitiationStatusValidated
		}

		var connectorID models.ConnectorID
		if payload.ConnectorID == "" {
			provider, err := models.ConnectorProviderFromString(payload.Provider)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}

			connectors, err := repo.ListConnectorsByProvider(r.Context(), provider)
			if err != nil {
				handleStorageErrors(w, r, err)
				return
			}

			if len(connectors) == 0 {
				api.BadRequest(w, ErrValidation, fmt.Errorf("no connector found for provider %s", provider))
				return
			}

			if len(connectors) > 1 {
				api.BadRequest(w, ErrValidation, fmt.Errorf("multiple connectors found for provider %s", provider))
				return
			}

			connectorID = connectors[0].ID
		} else {
			var err error
			connectorID, err = models.ConnectorIDFromString(payload.ConnectorID)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}
		}

		isInstalled, _ := repo.IsInstalledByConnectorID(r.Context(), connectorID)
		if !isInstalled {
			api.BadRequest(w, ErrValidation, fmt.Errorf("connector %s is not installed", payload.ConnectorID))
			return
		}

		if payload.SourceAccountID != "" {
			_, err := repo.GetAccount(r.Context(), payload.SourceAccountID)
			if err != nil {
				handleStorageErrors(w, r, fmt.Errorf("failed to get source account: %w", err))
				return
			}
		}

		_, err := repo.GetAccount(r.Context(), payload.DestinationAccountID)
		if err != nil {
			handleStorageErrors(w, r, fmt.Errorf("failed to get destination account: %w", err))
			return
		}

		if payload.ScheduledAt.IsZero() {
			payload.ScheduledAt = time.Now().UTC()
		}

		createdAt := time.Now()
		tf := &models.TransferInitiation{
			ID: models.TransferInitiationID{
				Reference:   payload.Reference,
				ConnectorID: connectorID,
			},
			CreatedAt:            createdAt,
			UpdatedAt:            createdAt, // When created, should be the same
			ScheduledAt:          payload.ScheduledAt,
			Description:          payload.Description,
			DestinationAccountID: models.MustAccountIDFromString(payload.DestinationAccountID),
			ConnectorID:          connectorID,
			Type:                 models.MustTransferInitiationTypeFromString(payload.Type),
			Amount:               payload.Amount,
			Asset:                models.Asset(payload.Asset),
			Status:               status,
		}

		if payload.SourceAccountID != "" {
			tf.SourceAccountID = models.MustAccountIDFromString(payload.SourceAccountID)
		}

		if err := repo.CreateTransferInitiation(r.Context(), tf); err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if err := publisher.Publish(
			events.TopicPayments,
			publish.NewMessage(
				r.Context(),
				messages.NewEventSavedTransferInitiations(tf),
			),
		); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		if status == models.TransferInitiationStatusValidated {
			connector, err := repo.GetConnector(r.Context(), connectorID)
			if err != nil {
				handleStorageErrors(w, r, err)
				return
			}

			f, ok := paymentHandlers[connector.Provider]
			if !ok {
				api.InternalServerError(w, r, fmt.Errorf("no payment handler for provider %v", payload.ConnectorID))
				return
			}

			err = f(r.Context(), tf)
			if err != nil {
				switch {
				case errors.Is(err, integration.ErrValidation):
					api.BadRequest(w, ErrValidation, err)
				case errors.Is(err, integration.ErrConnectorNotFound):
					api.BadRequest(w, ErrValidation, err)
				default:
					api.InternalServerError(w, r, err)
				}
				return
			}
		}

		data := &transferInitiationResponse{
			ID:                   tf.ID.String(),
			CreatedAt:            tf.CreatedAt,
			UpdatedAt:            tf.UpdatedAt,
			ScheduledAt:          tf.ScheduledAt,
			Description:          tf.Description,
			SourceAccountID:      tf.SourceAccountID.String(),
			DestinationAccountID: tf.DestinationAccountID.String(),
			ConnectorID:          connectorID.String(),
			Type:                 tf.Type.String(),
			Amount:               tf.Amount,
			Asset:                tf.Asset.String(),
			Status:               tf.Status.String(),
			Error:                tf.Error,
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

type udateTransferInitiationStatusRepository interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
}

type updateTransferInitiationStatusRequest struct {
	Status string `json:"status"`
}

func updateTransferInitiationStatusHandler(
	repo udateTransferInitiationStatusRepository,
	publisher message.Publisher,
	paymentHandlers map[models.ConnectorProvider]paymentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &updateTransferInitiationStatusRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		status, err := models.TransferInitiationStatusFromString(payload.Status)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		switch status {
		case models.TransferInitiationStatusWaitingForValidation:
			api.BadRequest(w, ErrValidation, errors.New("cannot set back transfer initiation status to waiting for validation"))
			return
		case models.TransferInitiationStatusFailed,
			models.TransferInitiationStatusProcessed,
			models.TransferInitiationStatusProcessing:
			api.BadRequest(w, ErrValidation, errors.New("either VALIDATED or REJECTED status can be set"))
			return
		default:
		}

		transferID, err := models.TransferInitiationIDFromString((mux.Vars(r)["transferID"]))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		previousTransferInitiation, err := repo.ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if previousTransferInitiation.Status != models.TransferInitiationStatusWaitingForValidation {
			api.BadRequest(w, ErrValidation, errors.New("only waiting for validation transfer initiation can be updated"))
			return
		}
		previousTransferInitiation.Status = status
		previousTransferInitiation.Attempts++

		err = repo.UpdateTransferInitiationPaymentsStatus(r.Context(), transferID, nil, status, "", previousTransferInitiation.Attempts, time.Now())
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if err := publisher.Publish(
			events.TopicPayments,
			publish.NewMessage(
				r.Context(),
				messages.NewEventSavedTransferInitiations(previousTransferInitiation),
			),
		); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		if status == models.TransferInitiationStatusValidated {
			f, ok := paymentHandlers[previousTransferInitiation.Provider]
			if !ok {
				api.InternalServerError(w, r, errors.New("no payment handler for provider "+previousTransferInitiation.Provider.String()))
				return
			}

			err = f(r.Context(), previousTransferInitiation)
			if err != nil {
				switch {
				case errors.Is(err, integration.ErrValidation):
					api.BadRequest(w, ErrValidation, err)
				case errors.Is(err, integration.ErrConnectorNotFound):
					api.BadRequest(w, ErrValidation, err)
				default:
					api.InternalServerError(w, r, err)
				}
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type retryTransferInitiationRepository interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
}

func retryTransferInitiationHandler(
	repo retryTransferInitiationRepository,
	publisher message.Publisher,
	paymentHandlers map[models.ConnectorProvider]paymentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferID, err := models.TransferInitiationIDFromString((mux.Vars(r)["transferID"]))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		previousTransferInitiation, err := repo.ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if previousTransferInitiation.Status != models.TransferInitiationStatusFailed {
			api.BadRequest(w, ErrValidation, errors.New("only failed transfer initiation can be updated"))
			return
		}
		previousTransferInitiation.Status = models.TransferInitiationStatusProcessing
		previousTransferInitiation.Attempts++

		err = repo.UpdateTransferInitiationPaymentsStatus(r.Context(), transferID, nil, models.TransferInitiationStatusProcessing, "", previousTransferInitiation.Attempts, time.Now())
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if err := publisher.Publish(
			events.TopicPayments,
			publish.NewMessage(
				r.Context(),
				messages.NewEventSavedTransferInitiations(previousTransferInitiation),
			),
		); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		f, ok := paymentHandlers[previousTransferInitiation.Provider]
		if !ok {
			api.InternalServerError(w, r, errors.New("no payment handler for provider "+previousTransferInitiation.Provider.String()))
			return
		}

		err = f(r.Context(), previousTransferInitiation)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type deleteTransferInitiationRepository interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	DeleteTransferInitiation(ctx context.Context, id models.TransferInitiationID) error
}

func deleteTransferInitiationHandler(
	repo deleteTransferInitiationRepository,
	publisher message.Publisher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferID, err := models.TransferInitiationIDFromString(mux.Vars(r)["transferID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		tf, err := repo.ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if tf.Status != models.TransferInitiationStatusWaitingForValidation {
			api.BadRequest(w, ErrValidation, errors.New("cannot delete transfer initiation not waiting for validation"))
			return
		}

		err = repo.DeleteTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if err := publisher.Publish(
			events.TopicPayments,
			publish.NewMessage(
				r.Context(),
				messages.NewEventDeleteTransferInitiation(tf.ID),
			),
		); err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
