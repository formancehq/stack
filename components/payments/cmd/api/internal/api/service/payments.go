package service

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/pkg/errors"
)

type CreatePaymentRequest struct {
	Reference            string    `json:"reference"`
	ConnectorID          string    `json:"connectorID"`
	CreatedAt            time.Time `json:"createdAt"`
	Amount               *big.Int  `json:"amount"`
	Type                 string    `json:"type"`
	Status               string    `json:"status"`
	Scheme               string    `json:"scheme"`
	Asset                string    `json:"asset"`
	SourceAccountID      string    `json:"sourceAccountID"`
	DestinationAccountID string    `json:"destinationAccountID"`
}

func (r *CreatePaymentRequest) Validate() error {
	if r.Reference == "" {
		return errors.New("reference is required")
	}

	if r.ConnectorID == "" {
		return errors.New("connectorID is required")
	}

	if r.CreatedAt.IsZero() || r.CreatedAt.After(time.Now()) {
		return errors.New("createdAt is empty or in the future")
	}

	if r.Amount == nil {
		return errors.New("amount is required")
	}

	if r.Type == "" {
		return errors.New("type is required")
	}

	if _, err := models.PaymentTypeFromString(r.Type); err != nil {
		return errors.Wrap(err, "invalid type")
	}

	if r.Status == "" {
		return errors.New("status is required")
	}

	if _, err := models.PaymentStatusFromString(r.Status); err != nil {
		return errors.Wrap(err, "invalid status")
	}

	if r.Scheme == "" {
		return errors.New("scheme is required")
	}

	if _, err := models.PaymentSchemeFromString(r.Scheme); err != nil {
		return errors.Wrap(err, "invalid scheme")
	}

	if r.Asset == "" {
		return errors.New("asset is required")
	}

	_, _, err := models.GetCurrencyAndPrecisionFromAsset(models.Asset(r.Asset))
	if err != nil {
		return errors.Wrap(err, "invalid asset")
	}

	return nil
}

func (s *Service) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*models.Payment, error) {
	connectorID, err := models.ConnectorIDFromString(req.ConnectorID)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	isInstalled, err := s.store.IsConnectorInstalledByConnectorID(ctx, connectorID)
	if err != nil {
		return nil, newStorageError(err, "checking if connector is installed")
	}

	if !isInstalled {
		return nil, errors.Wrap(ErrValidation, "connector is not installed")
	}

	var sourceAccountID *models.AccountID
	if req.SourceAccountID != "" {
		sourceAccountID, err = models.AccountIDFromString(req.SourceAccountID)
		if err != nil {
			return nil, errors.Wrap(ErrValidation, err.Error())
		}
	}

	var destinationAccountID *models.AccountID
	if req.DestinationAccountID != "" {
		destinationAccountID, err = models.AccountIDFromString(req.DestinationAccountID)
		if err != nil {
			return nil, errors.Wrap(ErrValidation, err.Error())
		}
	}

	raw, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: req.Reference,
				Type:      models.PaymentType(req.Type),
			},
			ConnectorID: connectorID,
		},
		ConnectorID:          connectorID,
		CreatedAt:            req.CreatedAt,
		Reference:            req.Reference,
		Amount:               req.Amount,
		InitialAmount:        req.Amount,
		Type:                 models.PaymentType(req.Type),
		Status:               models.PaymentStatus(req.Status),
		Scheme:               models.PaymentScheme(req.Scheme),
		Asset:                models.Asset(req.Asset),
		SourceAccountID:      sourceAccountID,
		DestinationAccountID: destinationAccountID,
		RawData:              raw,
	}

	err = s.store.UpsertPayments(ctx, []*models.Payment{payment})
	if err != nil {
		return nil, newStorageError(err, "creating payment")
	}

	err = s.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, s.messages.NewEventSavedPayments(connectorID.Provider, payment)))
	if err != nil {
		return nil, errors.Wrap(err, "publishing message")
	}

	return payment, nil
}

func (s *Service) ListPayments(ctx context.Context, q storage.ListPaymentsQuery) (*api.Cursor[models.Payment], error) {
	cursor, err := s.store.ListPayments(ctx, q)
	return cursor, newStorageError(err, "listing payments")
}

func (s *Service) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	_, err := models.PaymentIDFromString(id)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	payment, err := s.store.GetPayment(ctx, id)
	return payment, newStorageError(err, "getting payment")
}

type UpdateMetadataRequest map[string]string

func (s *Service) UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error {
	err := s.store.UpdatePaymentMetadata(ctx, paymentID, metadata)
	return newStorageError(err, "updating payment metadata")
}
