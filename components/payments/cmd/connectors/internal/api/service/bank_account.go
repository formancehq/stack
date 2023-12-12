package service

import (
	"context"
	"time"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
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

	if c.ConnectorID == "" {
		return errors.New("connectorID must be provided")
	}

	return nil
}

func (s *Service) CreateBankAccount(ctx context.Context, req *CreateBankAccountRequest) (*models.BankAccount, error) {
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

	bankAccount := &models.BankAccount{
		CreatedAt:     time.Now().UTC(),
		AccountNumber: req.AccountNumber,
		IBAN:          req.IBAN,
		SwiftBicCode:  req.SwiftBicCode,
		Country:       req.Country,
		ConnectorID:   connectorID,
		Name:          req.Name,
		Metadata:      req.Metadata,
	}
	err = s.store.CreateBankAccount(ctx, bankAccount)
	if err != nil {
		return nil, newStorageError(err, "creating bank account")
	}

	if err := handlers.BankAccountHandler(ctx, bankAccount); err != nil {
		switch {
		case errors.Is(err, manager.ErrValidation):
			return nil, errors.Wrap(ErrValidation, err.Error())
		case errors.Is(err, manager.ErrConnectorNotFound):
			return nil, errors.Wrap(ErrValidation, err.Error())
		default:
			return nil, err
		}
	}

	return bankAccount, nil
}
