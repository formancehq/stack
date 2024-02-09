package ingestion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/internal/models"
)

type MockStore struct {
	paymentIDsNotModified map[string]struct{}
}

func NewMockStore() *MockStore {
	return &MockStore{
		paymentIDsNotModified: make(map[string]struct{}),
	}
}

func (m *MockStore) WithPaymentIDsNotModified(paymentsIDs []models.PaymentID) *MockStore {
	for _, id := range paymentsIDs {
		m.paymentIDsNotModified[id.String()] = struct{}{}
	}
	return m
}

func (m *MockStore) UpsertAccounts(ctx context.Context, accounts []*models.Account) error {
	return nil
}

func (m *MockStore) UpsertPayments(ctx context.Context, payments []*models.Payment) ([]*models.PaymentID, error) {
	ids := make([]*models.PaymentID, 0, len(payments))
	for _, payment := range payments {
		if _, ok := m.paymentIDsNotModified[payment.ID.String()]; !ok {
			ids = append(ids, &payment.ID)
		}
	}

	return ids, nil
}

func (m *MockStore) UpsertPaymentsAdjustments(ctx context.Context, adjustments []*models.PaymentAdjustment) error {
	return nil
}

func (m *MockStore) UpsertPaymentsMetadata(ctx context.Context, metadata []*models.PaymentMetadata) error {
	return nil
}

func (m *MockStore) InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error {
	return nil
}

func (m *MockStore) UpdateTaskState(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, state json.RawMessage) error {
	return nil
}

func (m *MockStore) UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, adjustment *models.TransferInitiationAdjustment) error {
	return nil
}

func (m *MockStore) UpdateTransferReversalStatus(ctx context.Context, transfer *models.TransferInitiation, transferReversal *models.TransferReversal, adjustment *models.TransferInitiationAdjustment) error {
	return nil
}

func (m *MockStore) AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, updatedAt time.Time, metadata map[string]string) error {
	return nil
}

func (m *MockStore) AddBankAccountRelatedAccount(ctx context.Context, adjustment *models.BankAccountRelatedAccount) error {
	return nil
}

type MockPublisher struct {
	messages chan *message.Message
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		messages: make(chan *message.Message, 100),
	}
}

func (m *MockPublisher) Publish(topic string, messages ...*message.Message) error {
	for _, msg := range messages {
		m.messages <- msg
	}

	return nil
}

func (m *MockPublisher) Close() error {
	close(m.messages)
	return nil
}
