package ingestion

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/formancehq/payments/internal/app/models"
)

type Ingester interface {
	IngestPayments(ctx context.Context, batch PaymentBatch, commitState any) error
	IngestAccounts(ctx context.Context, batch AccountBatch) error
	IngestBalances(ctx context.Context, batch BalanceBatch, checkIfAccountExists bool) error
	GetTransfer(ctx context.Context, transferID uuid.UUID) (models.Transfer, error)
	UpdateTransferStatus(ctx context.Context, transferID uuid.UUID, status models.TransferStatus, reference, err string) error
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

	GetTransfer(ctx context.Context, transferID uuid.UUID) (models.Transfer, error)
	UpdateTransferStatus(ctx context.Context, transferID uuid.UUID, status models.TransferStatus, reference, err string) error
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
