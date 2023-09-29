package ingestion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/internal/app/models"
)

type Ingester interface {
	IngestAccounts(ctx context.Context, batch AccountBatch) error
	IngestPayments(ctx context.Context, batch PaymentBatch, commitState any) error
	IngestBalances(ctx context.Context, batch BalanceBatch, checkIfAccountExists bool) error
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
	AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, updatedAt time.Time) error
}

type DefaultIngester struct {
	repo       Repository
	provider   models.ConnectorProvider
	descriptor models.TaskDescriptor
	publisher  message.Publisher
}

type Repository interface {
	UpsertAccounts(ctx context.Context, provider models.ConnectorProvider, accounts []*models.Account) error
	UpsertPayments(ctx context.Context, provider models.ConnectorProvider, payments []*models.Payment) error
	InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error
	UpdateTaskState(ctx context.Context, provider models.ConnectorProvider, descriptor models.TaskDescriptor, state json.RawMessage) error
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
	AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, updatedAt time.Time) error
}

func NewDefaultIngester(
	provider models.ConnectorProvider,
	descriptor models.TaskDescriptor,
	repo Repository,
	publisher message.Publisher,
) *DefaultIngester {
	return &DefaultIngester{
		provider:   provider,
		descriptor: descriptor,
		repo:       repo,
		publisher:  publisher,
	}
}
