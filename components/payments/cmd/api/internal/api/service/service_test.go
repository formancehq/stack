package service

import (
	"context"

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
