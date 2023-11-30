package service

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/messages"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/pkg/errors"
)

type CreateBankAccountRequest struct {
	AccountNumber string `json:"accountNumber"`
	IBAN          string `json:"iban"`
	SwiftBicCode  string `json:"swiftBicCode"`
	Country       string `json:"country"`
	ConnectorID   string `json:"connectorID"`
	Name          string `json:"name"`
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

	if connector.Provider != models.ConnectorProviderBankingCircle {
		// For now, bank accounts can only be created for BankingCircle
		// in the future, we will support other providers
		return nil, errors.Wrap(ErrValidation, "provider not supported")
	}

	bankAccount := &models.BankAccount{
		CreatedAt:     time.Now().UTC(),
		AccountNumber: req.AccountNumber,
		IBAN:          req.IBAN,
		SwiftBicCode:  req.SwiftBicCode,
		Country:       req.Country,
		ConnectorID:   connectorID,
		Name:          req.Name,
	}
	err = s.store.CreateBankAccount(ctx, bankAccount)
	if err != nil {
		return nil, newStorageError(err, "creating bank account")
	}

	// BankingCircle does not have external accounts so we need to create
	// one by hand
	if connector.Provider == models.ConnectorProviderBankingCircle {
		accountID := models.AccountID{
			Reference:   bankAccount.ID.String(),
			ConnectorID: connector.ID,
		}
		err = s.store.UpsertAccounts(ctx, []*models.Account{
			{
				ID:          accountID,
				CreatedAt:   time.Now(),
				Reference:   bankAccount.ID.String(),
				ConnectorID: connector.ID,
				AccountName: bankAccount.Name,
				Type:        models.AccountTypeExternalFormance,
			},
		})
		if err != nil {
			return nil, newStorageError(err, "upserting accounts")
		}

		err = s.store.LinkBankAccountWithAccount(ctx, bankAccount.ID, &accountID)
		if err != nil {
			return nil, newStorageError(err, "linking bank account with account")
		}

		bankAccount.AccountID = &accountID
	}

	if err := s.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			messages.NewEventSavedBankAccounts(bankAccount),
		),
	); err != nil {
		return nil, errors.Wrap(ErrPublish, err.Error())
	}

	return bankAccount, nil
}
