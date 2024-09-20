package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/publish"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/pkg/errors"
)

type CreateAccountRequest struct {
	Reference    string            `json:"reference"`
	ConnectorID  string            `json:"connectorID"`
	CreatedAt    time.Time         `json:"createdAt"`
	DefaultAsset string            `json:"defaultAsset"`
	AccountName  string            `json:"accountName"`
	Type         string            `json:"type"`
	Metadata     map[string]string `json:"metadata"`
}

func (r *CreateAccountRequest) Validate() error {
	if r.Reference == "" {
		return errors.New("reference is required")
	}

	if r.ConnectorID == "" {
		return errors.New("connectorID is required")
	}

	if r.CreatedAt.IsZero() || r.CreatedAt.After(time.Now()) {
		return errors.New("createdAt is empty or in the future")
	}

	if r.AccountName == "" {
		return errors.New("accountName is required")
	}

	if r.Type == "" {
		return errors.New("type is required")
	}

	_, err := models.AccountTypeFromString(r.Type)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*models.Account, error) {
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

	accountType, err := models.AccountTypeFromString(req.Type)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	raw, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	account := &models.Account{
		ID: models.AccountID{
			Reference:   req.Reference,
			ConnectorID: connectorID,
		},
		ConnectorID:  connectorID,
		CreatedAt:    req.CreatedAt,
		Reference:    req.Reference,
		DefaultAsset: models.Asset(req.DefaultAsset),
		AccountName:  req.AccountName,
		Type:         accountType,
		Metadata:     req.Metadata,
		RawData:      raw,
	}

	err = s.store.UpsertAccounts(ctx, []*models.Account{account})
	if err != nil {
		return nil, newStorageError(err, "creating account")
	}

	err = s.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, s.messages.NewEventSavedAccounts(connectorID.Provider, account)))
	if err != nil {
		return nil, errors.Wrap(err, "publishing message")
	}

	return account, nil
}

func (s *Service) ListAccounts(ctx context.Context, q storage.ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error) {
	cursor, err := s.store.ListAccounts(ctx, q)
	return cursor, newStorageError(err, "listing accounts")
}

func (s *Service) GetAccount(
	ctx context.Context,
	accountID string,
) (*models.Account, error) {
	_, err := models.AccountIDFromString(accountID)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	account, err := s.store.GetAccount(ctx, accountID)
	return account, newStorageError(err, "getting account")
}
