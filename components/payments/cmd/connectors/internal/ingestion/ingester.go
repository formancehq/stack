package ingestion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

type Ingester interface {
	IngestAccounts(ctx context.Context, batch AccountBatch) error
	IngestPayments(ctx context.Context, connectorID models.ConnectorID, batch PaymentBatch, commitState any) error
	IngestBalances(ctx context.Context, batch BalanceBatch, checkIfAccountExists bool) error
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
	AddTransferInitiationPaymentID(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, updatedAt time.Time) error
	LinkBankAccountWithAccount(ctx context.Context, bankAccount *models.BankAccount, accountID *models.AccountID) error
}

type DefaultIngester struct {
	provider   models.ConnectorProvider
	store      Store
	descriptor models.TaskDescriptor
	publisher  message.Publisher
}

type Store interface {
	UpsertAccounts(ctx context.Context, accounts []*models.Account) error
	UpsertPayments(ctx context.Context, payments []*models.Payment) ([]*models.PaymentID, error)
	InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error
	UpdateTaskState(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, state json.RawMessage) error
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
	AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, updatedAt time.Time) error
	LinkBankAccountWithAccount(ctx context.Context, id uuid.UUID, accountID *models.AccountID) error
}

func NewDefaultIngester(
	provider models.ConnectorProvider,
	descriptor models.TaskDescriptor,
	repo Store,
	publisher message.Publisher,
) *DefaultIngester {
	return &DefaultIngester{
		provider:   provider,
		descriptor: descriptor,
		store:      repo,
		publisher:  publisher,
	}
}

var _ Ingester = (*DefaultIngester)(nil)
