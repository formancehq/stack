package service

import (
	"context"
	"math/big"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type MockStore struct {
	isConnectorInstalled bool
}

func (m *MockStore) WithIsConnectorInstalled(isConnectorInstalled bool) *MockStore {
	m.isConnectorInstalled = isConnectorInstalled
	return m
}

func (m *MockStore) Ping() error {
	return nil
}

func (m *MockStore) IsConnectorInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error) {
	return m.isConnectorInstalled, nil
}

func (m *MockStore) ListBalances(ctx context.Context, q storage.ListBalancesQuery) (*api.Cursor[models.Balance], error) {
	return nil, nil
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

func (m *MockStore) ListAccounts(ctx context.Context, q storage.ListAccountsQuery) (*api.Cursor[models.Account], error) {
	return nil, nil
}

func (m *MockStore) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return nil, nil
}

func (m *MockStore) ListBankAccounts(ctx context.Context, q storage.ListBankAccountQuery) (*api.Cursor[models.BankAccount], error) {
	return nil, nil
}

func (m *MockStore) GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	return nil, nil
}

func (m *MockStore) UpsertPayments(ctx context.Context, payments []*models.Payment) error {
	return nil
}

func (m *MockStore) ListPayments(ctx context.Context, q storage.ListPaymentsQuery) (*api.Cursor[models.Payment], error) {
	return nil, nil
}

func (m *MockStore) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	return nil, nil
}

func (m *MockStore) UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error {
	return nil
}

func (m *MockStore) ListTransferInitiations(ctx context.Context, q storage.ListTransferInitiationsQuery) (*api.Cursor[models.TransferInitiation], error) {
	return nil, nil
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

func (m *MockStore) ListPools(ctx context.Context, q storage.ListPoolsQuery) (*api.Cursor[models.Pool], error) {
	return nil, nil
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

type MockPublisher struct {
	errorToSend  error
	messagesChan chan *message.Message
}

func (m *MockPublisher) WithError(err error) *MockPublisher {
	m.errorToSend = err
	return m
}

func (m *MockPublisher) WithMessagesChan(messagesChan chan *message.Message) *MockPublisher {
	m.messagesChan = messagesChan
	return m
}

func (m *MockPublisher) Publish(topic string, messages ...*message.Message) error {
	if m.errorToSend != nil {
		return m.errorToSend
	}

	if m.messagesChan != nil {
		for _, msg := range messages {
			m.messagesChan <- msg
		}
	}

	return nil
}

func (m *MockPublisher) Close() error {
	return nil
}
