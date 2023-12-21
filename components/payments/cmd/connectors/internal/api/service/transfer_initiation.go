package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/pkg/errors"
)

type CreateTransferInitiationRequest struct {
	Reference            string            `json:"reference"`
	ScheduledAt          time.Time         `json:"scheduledAt"`
	Description          string            `json:"description"`
	SourceAccountID      string            `json:"sourceAccountID"`
	DestinationAccountID string            `json:"destinationAccountID"`
	ConnectorID          string            `json:"connectorID"`
	Provider             string            `json:"provider"`
	Type                 string            `json:"type"`
	Amount               *big.Int          `json:"amount"`
	Asset                string            `json:"asset"`
	Validated            bool              `json:"validated"`
	Metadata             map[string]string `json:"metadata"`
}

func (r *CreateTransferInitiationRequest) Validate() error {
	if r.Reference == "" {
		return errors.New("reference is required")
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

func (s *Service) CreateTransferInitiation(ctx context.Context, req *CreateTransferInitiationRequest) (*models.TransferInitiation, error) {
	status := models.TransferInitiationStatusWaitingForValidation
	if req.Validated {
		status = models.TransferInitiationStatusValidated
	}

	var connectorID models.ConnectorID
	if req.ConnectorID == "" {
		provider, err := models.ConnectorProviderFromString(req.Provider)
		if err != nil {
			return nil, errors.Wrap(ErrValidation, err.Error())
		}

		connectors, err := s.store.ListConnectorsByProvider(ctx, provider)
		if err != nil {
			return nil, newStorageError(err, "listing connectors")
		}

		if len(connectors) == 0 {
			return nil, errors.Wrap(ErrValidation, fmt.Sprintf("no connector found for provider %s", provider))
		}

		if len(connectors) > 1 {
			return nil, errors.Wrap(ErrValidation, fmt.Sprintf("multiple connectors found for provider %s", provider))
		}

		connectorID = connectors[0].ID
	} else {
		var err error
		connectorID, err = models.ConnectorIDFromString(req.ConnectorID)
		if err != nil {
			return nil, errors.Wrap(ErrValidation, err.Error())
		}
	}

	isInstalled, _ := s.store.IsInstalledByConnectorID(ctx, connectorID)
	if !isInstalled {
		return nil, errors.Wrap(ErrValidation, fmt.Sprintf("connector %s is not installed", req.ConnectorID))
	}

	if req.SourceAccountID != "" {
		_, err := s.store.GetAccount(ctx, req.SourceAccountID)
		if err != nil {
			return nil, newStorageError(err, "getting source account")
		}
	}

	_, err := s.store.GetAccount(ctx, req.DestinationAccountID)
	if err != nil {
		return nil, newStorageError(err, "getting destination account")
	}

	createdAt := time.Now()
	tf := &models.TransferInitiation{
		ID: models.TransferInitiationID{
			Reference:   req.Reference,
			ConnectorID: connectorID,
		},
		CreatedAt:            createdAt,
		UpdatedAt:            createdAt, // When created, should be the same
		ScheduledAt:          req.ScheduledAt,
		Description:          req.Description,
		DestinationAccountID: models.MustAccountIDFromString(req.DestinationAccountID),
		ConnectorID:          connectorID,
		Provider:             connectorID.Provider,
		Type:                 models.MustTransferInitiationTypeFromString(req.Type),
		Amount:               req.Amount,
		Asset:                models.Asset(req.Asset),
		Status:               status,
		Metadata:             req.Metadata,
	}

	if req.SourceAccountID != "" {
		sID := models.MustAccountIDFromString(req.SourceAccountID)
		tf.SourceAccountID = &sID
	}

	if err := s.store.CreateTransferInitiation(ctx, tf); err != nil {
		return nil, newStorageError(err, "creating transfer initiation")
	}

	if err := s.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			s.messages.NewEventSavedTransferInitiations(tf),
		),
	); err != nil {
		return nil, errors.Wrap(ErrPublish, err.Error())
	}

	if status == models.TransferInitiationStatusValidated {
		connector, err := s.store.GetConnector(ctx, connectorID)
		if err != nil {
			return nil, newStorageError(err, "getting connector")
		}

		handlers, ok := s.connectorHandlers[connector.Provider]
		if !ok {
			return nil, errors.Wrap(ErrValidation, fmt.Sprintf("no payment handler for provider %v", connector.Provider))
		}

		err = handlers.PaymentHandler(ctx, tf)
		if err != nil {
			switch {
			case errors.Is(err, manager.ErrValidation):
				return nil, errors.Wrap(ErrValidation, err.Error())
			case errors.Is(err, manager.ErrConnectorNotFound):
				return nil, errors.Wrap(ErrValidation, err.Error())
			default:
				return nil, err
			}
		}
	}

	return tf, nil
}

type UpdateTransferInitiationStatusRequest struct {
	Status string `json:"status"`
}

func (r *UpdateTransferInitiationStatusRequest) Validate() error {
	if r.Status == "" {
		return errors.New("status is required")
	}

	return nil
}

func (s *Service) UpdateTransferInitiationStatus(ctx context.Context, id string, req *UpdateTransferInitiationStatusRequest) error {
	status, err := models.TransferInitiationStatusFromString(req.Status)
	if err != nil {
		return errors.Wrap(ErrValidation, err.Error())
	}

	switch status {
	case models.TransferInitiationStatusWaitingForValidation:
		return errors.Wrap(ErrValidation, "cannot set back transfer initiation status to waiting for validation")
	case models.TransferInitiationStatusFailed,
		models.TransferInitiationStatusProcessed,
		models.TransferInitiationStatusProcessing:
		return errors.Wrap(ErrValidation, "either VALIDATED or REJECTED status can be set")
	default:
	}

	transferID, err := models.TransferInitiationIDFromString(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	previousTransferInitiation, err := s.store.ReadTransferInitiation(ctx, transferID)
	if err != nil {
		return newStorageError(err, "reading transfer initiation")
	}

	if previousTransferInitiation.Status != models.TransferInitiationStatusWaitingForValidation {
		return errors.Wrap(ErrValidation, "only waiting for validation transfer initiation can be updated")
	}
	previousTransferInitiation.Status = status
	previousTransferInitiation.Attempts++

	err = s.store.UpdateTransferInitiationPaymentsStatus(ctx, transferID, nil, status, "", previousTransferInitiation.Attempts, time.Now())
	if err != nil {
		return newStorageError(err, "updating transfer initiation payments status")
	}

	if err := s.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			s.messages.NewEventSavedTransferInitiations(previousTransferInitiation),
		),
	); err != nil {
		return errors.Wrap(ErrPublish, err.Error())
	}

	if status == models.TransferInitiationStatusValidated {
		handlers, ok := s.connectorHandlers[previousTransferInitiation.Provider]
		if !ok {
			return errors.Wrap(ErrValidation, fmt.Sprintf("no payment handler for provider %v", previousTransferInitiation.Provider))
		}

		err = handlers.PaymentHandler(ctx, previousTransferInitiation)
		if err != nil {
			switch {
			case errors.Is(err, manager.ErrValidation):
				return errors.Wrap(ErrValidation, err.Error())
			case errors.Is(err, manager.ErrConnectorNotFound):
				return errors.Wrap(ErrValidation, err.Error())
			default:
				return err
			}
		}
	}

	return nil
}

func (s *Service) RetryTransferInitiation(ctx context.Context, id string) error {
	transferID, err := models.TransferInitiationIDFromString(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	previousTransferInitiation, err := s.store.ReadTransferInitiation(ctx, transferID)
	if err != nil {
		return newStorageError(err, "reading transfer initiation")
	}

	if previousTransferInitiation.Status != models.TransferInitiationStatusFailed {
		return errors.Wrap(ErrValidation, "only failed transfer initiation can be retried")
	}
	previousTransferInitiation.Status = models.TransferInitiationStatusProcessing
	previousTransferInitiation.Attempts++

	err = s.store.UpdateTransferInitiationPaymentsStatus(ctx, transferID, nil, models.TransferInitiationStatusProcessing, "", previousTransferInitiation.Attempts, time.Now())
	if err != nil {
		return newStorageError(err, "updating transfer initiation payments status")
	}

	if err := s.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			s.messages.NewEventSavedTransferInitiations(previousTransferInitiation),
		),
	); err != nil {
		return errors.Wrap(ErrPublish, err.Error())
	}

	handlers, ok := s.connectorHandlers[previousTransferInitiation.Provider]
	if !ok {
		return errors.Wrap(ErrValidation, fmt.Sprintf("no payment handler for provider %v", previousTransferInitiation.Provider))
	}

	err = handlers.PaymentHandler(ctx, previousTransferInitiation)
	if err != nil {
		switch {
		case errors.Is(err, manager.ErrValidation):
			return errors.Wrap(ErrValidation, err.Error())
		case errors.Is(err, manager.ErrConnectorNotFound):
			return errors.Wrap(ErrValidation, err.Error())
		default:
			return err
		}
	}

	return nil
}

func (s *Service) DeleteTransferInitiation(ctx context.Context, id string) error {
	transferID, err := models.TransferInitiationIDFromString(id)
	if err != nil {
		return errors.Wrap(ErrInvalidID, err.Error())
	}

	tf, err := s.store.ReadTransferInitiation(ctx, transferID)
	if err != nil {
		return newStorageError(err, "reading transfer initiation")
	}

	if tf.Status != models.TransferInitiationStatusWaitingForValidation {
		return errors.Wrap(ErrValidation, "only waiting for validation transfer initiation can be deleted")
	}

	err = s.store.DeleteTransferInitiation(ctx, transferID)
	if err != nil {
		return newStorageError(err, "deleting transfer initiation")
	}

	if err := s.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			s.messages.NewEventDeleteTransferInitiation(tf.ID),
		),
	); err != nil {
		return errors.Wrap(ErrPublish, err.Error())
	}

	return nil
}
