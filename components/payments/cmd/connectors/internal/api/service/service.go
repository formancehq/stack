package service

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

type PaymentHandler func(ctx context.Context, transfer *models.TransferInitiation) error

type Store interface {
	Ping() error

	GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error)
	ListConnectors(ctx context.Context) ([]*models.Connector, error)

	UpsertAccounts(ctx context.Context, accounts []*models.Account) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)

	CreateBankAccount(ctx context.Context, account *models.BankAccount) error
	LinkBankAccountWithAccount(ctx context.Context, id uuid.UUID, accountID *models.AccountID) error

	ListConnectorsByProvider(ctx context.Context, provider models.ConnectorProvider) ([]*models.Connector, error)
	IsInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error)

	CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
	DeleteTransferInitiation(ctx context.Context, id models.TransferInitiationID) error
}

type Service struct {
	store           Store
	publisher       message.Publisher
	paymentHandlers map[models.ConnectorProvider]PaymentHandler
}

func New(
	store Store,
	publisher message.Publisher,
	paymentHandlers map[models.ConnectorProvider]PaymentHandler,
) *Service {
	return &Service{
		store:           store,
		publisher:       publisher,
		paymentHandlers: paymentHandlers,
	}
}
