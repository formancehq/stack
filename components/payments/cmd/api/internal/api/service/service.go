package service

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error
	IsConnectorInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error)
	ListAccounts(ctx context.Context, q storage.ListAccountsQuery) (*api.Cursor[models.Account], error)
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	ListBalances(ctx context.Context, q storage.ListBalancesQuery) (*api.Cursor[models.Balance], error)
	GetBalancesAt(ctx context.Context, accountID models.AccountID, at time.Time) ([]*models.Balance, error)
	ListBankAccounts(ctx context.Context, q storage.ListBankAccountQuery) (*api.Cursor[models.BankAccount], error)
	GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error)
	UpsertPayments(ctx context.Context, payments []*models.Payment) error
	ListPayments(ctx context.Context, q storage.ListPaymentsQuery) (*api.Cursor[models.Payment], error)
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error
	ListTransferInitiations(ctx context.Context, q storage.ListTransferInitiationsQuery) (*api.Cursor[models.TransferInitiation], error)
	GetTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	CreatePool(ctx context.Context, pool *models.Pool) error
	AddAccountToPool(ctx context.Context, poolAccount *models.PoolAccounts) error
	AddAccountsToPool(ctx context.Context, poolAccounts []*models.PoolAccounts) error
	RemoveAccountFromPool(ctx context.Context, poolAccount *models.PoolAccounts) error
	ListPools(ctx context.Context, q storage.ListPoolsQuery) (*api.Cursor[models.Pool], error)
	GetPool(ctx context.Context, poolID uuid.UUID) (*models.Pool, error)
	DeletePool(ctx context.Context, poolID uuid.UUID) error
}

type Service struct {
	store     Store
	publisher message.Publisher
	messages  *messages.Messages
}

func New(store Store, publisher message.Publisher, messages *messages.Messages) *Service {
	return &Service{
		store:     store,
		publisher: publisher,
		messages:  messages,
	}
}
