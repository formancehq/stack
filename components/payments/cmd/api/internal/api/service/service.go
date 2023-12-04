package service

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error
	IsConnectorInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error)
	ListAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Account, storage.PaginationDetails, error)
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	ListBalances(ctx context.Context, balanceQuery storage.BalanceQuery) ([]*models.Balance, storage.PaginationDetails, error)
	GetBalancesAt(ctx context.Context, accountID models.AccountID, at time.Time) ([]*models.Balance, error)
	ListBankAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.BankAccount, storage.PaginationDetails, error)
	GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error)
	UpsertPayments(ctx context.Context, payments []*models.Payment) error
	ListPayments(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Payment, storage.PaginationDetails, error)
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error
	ListTransferInitiations(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.TransferInitiation, storage.PaginationDetails, error)
	GetTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	CreatePool(ctx context.Context, pool *models.Pool) error
	AddAccountToPool(ctx context.Context, poolAccount *models.PoolAccounts) error
	AddAccountsToPool(ctx context.Context, poolAccounts []*models.PoolAccounts) error
	RemoveAccountFromPool(ctx context.Context, poolAccount *models.PoolAccounts) error
	ListPools(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Pool, storage.PaginationDetails, error)
	GetPool(ctx context.Context, poolID uuid.UUID) (*models.Pool, error)
	DeletePool(ctx context.Context, poolID uuid.UUID) error
}

type Service struct {
	store     Store
	publisher message.Publisher
}

func New(store Store, publisher message.Publisher) *Service {
	return &Service{
		store:     store,
		publisher: publisher,
	}
}
