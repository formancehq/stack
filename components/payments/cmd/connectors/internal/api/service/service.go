package service

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

type InitiatePaymentHandler func(ctx context.Context, transfer *models.TransferInitiation) error
type ReversePaymentHandler func(ctx context.Context, transfer *models.TransferReversal) error
type BankAccountHandler func(ctx context.Context, connectorID models.ConnectorID, bankAccount *models.BankAccount) error

type Store interface {
	Ping() error

	GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error)
	ListConnectors(ctx context.Context) ([]*models.Connector, error)

	UpsertAccounts(ctx context.Context, accounts []*models.Account) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)

	CreateBankAccount(ctx context.Context, account *models.BankAccount) error
	UpdateBankAccountMetadata(ctx context.Context, id uuid.UUID, metadata map[string]string) error
	GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error)
	GetBankAccountRelatedAccounts(ctx context.Context, id uuid.UUID) ([]*models.BankAccountRelatedAccount, error)

	ListConnectorsByProvider(ctx context.Context, provider models.ConnectorProvider) ([]*models.Connector, error)
	IsInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error)

	CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, adjustment *models.TransferInitiationAdjustment) error
	DeleteTransferInitiation(ctx context.Context, id models.TransferInitiationID) error

	CreateTransferReversal(ctx context.Context, transferReversal *models.TransferReversal) error
}

type Service struct {
	store             Store
	publisher         message.Publisher
	messages          *messages.Messages
	connectorHandlers map[models.ConnectorProvider]*ConnectorHandlers
}

type ConnectorHandlers struct {
	InitiatePaymentHandler InitiatePaymentHandler
	ReversePaymentHandler  ReversePaymentHandler
	BankAccountHandler     BankAccountHandler
}

func New(
	store Store,
	publisher message.Publisher,
	messages *messages.Messages,
	connectorHandlers map[models.ConnectorProvider]*ConnectorHandlers,
) *Service {
	return &Service{
		store:             store,
		publisher:         publisher,
		connectorHandlers: connectorHandlers,
		messages:          messages,
	}
}
