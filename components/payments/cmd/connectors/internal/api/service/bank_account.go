package service

import (
	"context"
	"time"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateBankAccountRequest struct {
	AccountNumber string            `json:"accountNumber"`
	IBAN          string            `json:"iban"`
	SwiftBicCode  string            `json:"swiftBicCode"`
	Country       string            `json:"country"`
	ConnectorID   string            `json:"connectorID"`
	Name          string            `json:"name"`
	Metadata      map[string]string `json:"metadata"`
}

func (c *CreateBankAccountRequest) Validate() error {
	if c.AccountNumber == "" && c.IBAN == "" {
		return errors.New("either accountNumber or iban must be provided")
	}

	if c.Name == "" {
		return errors.New("name must be provided")
	}

	if c.Country == "" {
		return errors.New("country must be provided")
	}

	return nil
}

func (s *Service) CreateBankAccount(ctx context.Context, req *CreateBankAccountRequest) (*models.BankAccount, error) {
	var handlers *ConnectorHandlers
	var connectorID models.ConnectorID
	if req.ConnectorID != "" {
		var err error
		connectorID, err = models.ConnectorIDFromString(req.ConnectorID)
		if err != nil {
			return nil, errors.Wrap(ErrValidation, err.Error())
		}

		connector, err := s.store.GetConnector(ctx, connectorID)
		if err != nil && errors.Is(err, storage.ErrNotFound) {
			return nil, errors.Wrap(ErrValidation, "connector not installed")
		} else if err != nil {
			return nil, newStorageError(err, "getting connector")
		}

		var ok bool
		handlers, ok = s.connectorHandlers[connector.Provider]
		if !ok || handlers.BankAccountHandler == nil {
			return nil, errors.Wrap(ErrValidation, "no bank account handler for connector")
		}
	}

	bankAccount := &models.BankAccount{
		CreatedAt:     time.Now().UTC(),
		AccountNumber: req.AccountNumber,
		IBAN:          req.IBAN,
		SwiftBicCode:  req.SwiftBicCode,
		Country:       req.Country,
		Name:          req.Name,
		Metadata:      req.Metadata,
	}
	err := s.store.CreateBankAccount(ctx, bankAccount)
	if err != nil {
		return nil, newStorageError(err, "creating bank account")
	}

	if handlers != nil {
		if err := handlers.BankAccountHandler(ctx, connectorID, bankAccount); err != nil {
			switch {
			case errors.Is(err, manager.ErrValidation):
				return nil, errors.Wrap(ErrValidation, err.Error())
			case errors.Is(err, manager.ErrConnectorNotFound):
				return nil, errors.Wrap(ErrValidation, err.Error())
			default:
				return nil, err
			}
		}

		adjustments, err := s.store.GetBankAccountAdjustments(ctx, bankAccount.ID)
		if err != nil {
			return nil, newStorageError(err, "fetching bank account")
		}

		bankAccount.Adjustments = adjustments
	}

	return bankAccount, nil
}

type ForwardBankAccountToConnectorRequest struct {
	ConnectorID string `json:"connectorID"`
}

func (f *ForwardBankAccountToConnectorRequest) Validate() error {
	if f.ConnectorID == "" {
		return errors.New("connectorID must be provided")
	}

	return nil
}

func (s *Service) ForwardBankAccountToConnector(ctx context.Context, id string, req *ForwardBankAccountToConnectorRequest) (*models.BankAccount, error) {
	bankAccountID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidID, err.Error())
	}

	connectorID, err := models.ConnectorIDFromString(req.ConnectorID)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	connector, err := s.store.GetConnector(ctx, connectorID)
	if err != nil && errors.Is(err, storage.ErrNotFound) {
		return nil, errors.Wrap(ErrValidation, "connector not installed")
	} else if err != nil {
		return nil, newStorageError(err, "getting connector")
	}

	handlers, ok := s.connectorHandlers[connector.Provider]
	if !ok || handlers.BankAccountHandler == nil {
		return nil, errors.Wrap(ErrValidation, "no bank account handler for connector")
	}

	bankAccount, err := s.store.GetBankAccount(ctx, bankAccountID, true)
	if err != nil {
		return nil, newStorageError(err, "fetching bank account")
	}

	for _, adjustment := range bankAccount.Adjustments {
		if adjustment.ConnectorID == connectorID {
			return nil, errors.Wrap(ErrValidation, "bank account already forwarded to connector")
		}
	}

	if err := handlers.BankAccountHandler(ctx, connectorID, bankAccount); err != nil {
		switch {
		case errors.Is(err, manager.ErrValidation):
			return nil, errors.Wrap(ErrValidation, err.Error())
		case errors.Is(err, manager.ErrConnectorNotFound):
			return nil, errors.Wrap(ErrValidation, err.Error())
		default:
			return nil, err
		}
	}

	adjustments, err := s.store.GetBankAccountAdjustments(ctx, bankAccount.ID)
	if err != nil {
		return nil, newStorageError(err, "fetching bank account")
	}
	bankAccount.Adjustments = adjustments

	return bankAccount, err
}
