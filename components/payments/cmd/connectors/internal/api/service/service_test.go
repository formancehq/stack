package service

import (
	"context"
	"math/big"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

var (
	connectorDummyPay = models.Connector{
		ID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
		Name:      "c1",
		CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
		Provider:  models.ConnectorProviderDummyPay,
	}

	connectorBankingCircle = models.Connector{
		ID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderBankingCircle,
		},
		Name:      "c2",
		CreatedAt: time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
		Provider:  models.ConnectorProviderBankingCircle,
	}

	transferInitiationWaiting = models.TransferInitiation{
		ID: models.TransferInitiationID{
			Reference:   "ref1",
			ConnectorID: connectorDummyPay.ID,
		},
		CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
		ScheduledAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
		Description: "test",
		Type:        models.TransferInitiationTypeTransfer,
		SourceAccountID: &models.AccountID{
			Reference:   "acc1",
			ConnectorID: connectorDummyPay.ID,
		},
		DestinationAccountID: models.AccountID{
			Reference:   "acc2",
			ConnectorID: connectorDummyPay.ID,
		},
		Provider:      models.ConnectorProviderDummyPay,
		ConnectorID:   connectorDummyPay.ID,
		Amount:        big.NewInt(100),
		InitialAmount: big.NewInt(100),
		Asset:         "EUR/2",
		RelatedAdjustments: []*models.TransferInitiationAdjustment{
			{
				ID: uuid.New(),
				TransferInitiationID: models.TransferInitiationID{
					Reference:   "ref1",
					ConnectorID: connectorDummyPay.ID,
				},
				CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Status:    models.TransferInitiationStatusWaitingForValidation,
			},
		},
	}

	transferInitiationFailed = models.TransferInitiation{
		ID: models.TransferInitiationID{
			Reference:   "ref2",
			ConnectorID: connectorDummyPay.ID,
		},
		CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
		ScheduledAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
		Description: "test",
		Type:        models.TransferInitiationTypeTransfer,
		SourceAccountID: &models.AccountID{
			Reference:   "acc1",
			ConnectorID: connectorDummyPay.ID,
		},
		DestinationAccountID: models.AccountID{
			Reference:   "acc2",
			ConnectorID: connectorDummyPay.ID,
		},
		Provider:      models.ConnectorProviderDummyPay,
		ConnectorID:   connectorDummyPay.ID,
		Amount:        big.NewInt(100),
		InitialAmount: big.NewInt(100),
		Asset:         "EUR/2",
		RelatedAdjustments: []*models.TransferInitiationAdjustment{
			{
				ID: uuid.New(),
				TransferInitiationID: models.TransferInitiationID{
					Reference:   "ref2",
					ConnectorID: connectorDummyPay.ID,
				},
				CreatedAt: time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
				Status:    models.TransferInitiationStatusFailed,
				Error:     "some error",
			},
			{
				ID: uuid.New(),
				TransferInitiationID: models.TransferInitiationID{
					Reference:   "ref2",
					ConnectorID: connectorDummyPay.ID,
				},
				CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Status:    models.TransferInitiationStatusWaitingForValidation,
			},
		},
	}
)

type MockStore struct {
	listConnectorsNB int
}

func (m *MockStore) WithListConnectorsNB(nb int) *MockStore {
	m.listConnectorsNB = nb
	return m
}

func (m *MockStore) Ping() error {
	return nil
}

func (m *MockStore) GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error) {
	if connectorID == connectorDummyPay.ID {
		return &connectorDummyPay, nil
	} else if connectorID == connectorBankingCircle.ID {
		return &connectorBankingCircle, nil
	}

	return nil, storage.ErrNotFound
}

func (m *MockStore) ListConnectors(ctx context.Context) ([]*models.Connector, error) {
	return []*models.Connector{&connectorDummyPay, &connectorBankingCircle}, nil
}

func (m *MockStore) UpsertAccounts(ctx context.Context, accounts []*models.Account) error {
	return nil
}

func (m *MockStore) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return nil, nil
}

func (m *MockStore) CreateBankAccount(ctx context.Context, account *models.BankAccount) error {
	account.ID = uuid.New()
	return nil
}

func (m *MockStore) LinkBankAccountWithAccount(ctx context.Context, id uuid.UUID, accountID *models.AccountID) error {
	return nil
}

func (m *MockStore) FetchLinkedAccountForBankAccount(ctx context.Context, bankAccountID uuid.UUID) (*models.AccountID, error) {
	return nil, nil
}

func (m *MockStore) ListConnectorsByProvider(ctx context.Context, provider models.ConnectorProvider) ([]*models.Connector, error) {
	switch m.listConnectorsNB {
	case 0:
		return []*models.Connector{}, nil
	case 1:
		return []*models.Connector{&connectorDummyPay}, nil
	default:
		return []*models.Connector{&connectorDummyPay, &connectorBankingCircle}, nil
	}
}

func (m *MockStore) IsInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error) {
	if connectorID == connectorDummyPay.ID {
		return true, nil
	}

	return false, nil
}

func (m *MockStore) CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error {
	return nil
}

func (m *MockStore) ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error) {
	if id == transferInitiationWaiting.ID {
		tc := transferInitiationWaiting
		return &tc, nil
	} else if id == transferInitiationFailed.ID {
		tc := transferInitiationFailed
		return &tc, nil
	}

	return nil, storage.ErrNotFound
}

func (m *MockStore) UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, adjustment *models.TransferInitiationAdjustment) error {
	return nil
}

func (m *MockStore) DeleteTransferInitiation(ctx context.Context, id models.TransferInitiationID) error {
	return nil
}

func (m *MockStore) CreateTransferReversal(ctx context.Context, transferReversal *models.TransferReversal) error {
	return nil
}

type MockPublisher struct {
	errorToSend error
}

func (m *MockPublisher) WithError(err error) *MockPublisher {
	m.errorToSend = err
	return m
}

func (m *MockPublisher) Publish(topic string, messages ...*message.Message) error {
	return m.errorToSend
}

func (m *MockPublisher) Close() error {
	return nil
}
