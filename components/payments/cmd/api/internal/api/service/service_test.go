package service

import (
	"context"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

type MockStore struct{}

func (m *MockStore) Ping() error {
	return nil
}
func (m *MockStore) ListBalances(ctx context.Context, query storage.BalanceQuery) ([]*models.Balance, storage.PaginationDetails, error) {
	return nil, storage.PaginationDetails{}, nil
}

func (m *MockStore) GetBalancesAt(ctx context.Context, accountID models.AccountID, atTime time.Time) ([]*models.Balance, error) {
	return []*models.Balance{
		{
			AccountID: accountID,
			Asset:     "EUR/2",
			Balance:   big.NewInt(100),
		},
		{
			AccountID: accountID,
			Asset:     "USD/2",
			Balance:   big.NewInt(150),
		},
	}, nil
}

func (m *MockStore) ListAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Account, storage.PaginationDetails, error) {
	return nil, storage.PaginationDetails{}, nil
}

func (m *MockStore) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return nil, nil
}

func (m *MockStore) ListBankAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.BankAccount, storage.PaginationDetails, error) {
	return nil, storage.PaginationDetails{}, nil
}

func (m *MockStore) GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	return nil, nil
}

func (m *MockStore) ListPayments(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Payment, storage.PaginationDetails, error) {
	return nil, storage.PaginationDetails{}, nil
}

func (m *MockStore) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	return nil, nil
}

func (m *MockStore) UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error {
	return nil
}

func (m *MockStore) ListTransferInitiations(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.TransferInitiation, storage.PaginationDetails, error) {
	return nil, storage.PaginationDetails{}, nil
}

func (m *MockStore) GetTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error) {
	return nil, nil
}

func (m *MockStore) CreatePool(ctx context.Context, pool *models.Pool) error {
	return nil
}

func (m *MockStore) AddAccountsToPool(ctx context.Context, poolAccounts []*models.PoolAccounts) error {
	return nil
}

func (m *MockStore) AddAccountToPool(ctx context.Context, poolAccount *models.PoolAccounts) error {
	return nil
}

func (m *MockStore) RemoveAccountFromPool(ctx context.Context, poolAccount *models.PoolAccounts) error {
	return nil
}

func (m *MockStore) ListPools(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Pool, storage.PaginationDetails, error) {
	return nil, storage.PaginationDetails{}, nil
}

func (m *MockStore) GetPool(ctx context.Context, poolID uuid.UUID) (*models.Pool, error) {
	return &models.Pool{
		ID:   poolID,
		Name: "test",
		PoolAccounts: []*models.PoolAccounts{
			{
				PoolID: poolID,
				AccountID: models.AccountID{
					Reference: "acc1",
					ConnectorID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  models.ConnectorProviderDummyPay,
					},
				},
			},
			{
				PoolID: poolID,
				AccountID: models.AccountID{
					Reference: "acc2",
					ConnectorID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  models.ConnectorProviderDummyPay,
					},
				},
			},
		},
	}, nil
}

func (m *MockStore) DeletePool(ctx context.Context, poolID uuid.UUID) error {
	return nil
}
