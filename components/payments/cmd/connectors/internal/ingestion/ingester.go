package ingestion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
)

type Ingester interface {
	IngestAccounts(ctx context.Context, batch AccountBatch) error
	IngestPayments(ctx context.Context, batch PaymentBatch) error
	IngestBalances(ctx context.Context, batch BalanceBatch, checkIfAccountExists bool) error
	UpdateTaskState(ctx context.Context, state any) error
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, updatedAt time.Time) error
	UpdateTransferReversalStatus(ctx context.Context, transfer *models.TransferInitiation, transferReversal *models.TransferReversal) error
	AddTransferInitiationPaymentID(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, updatedAt time.Time) error
	LinkBankAccountWithAccount(ctx context.Context, bankAccount *models.BankAccount, accountID *models.AccountID) error
}

type DefaultIngester struct {
	provider    models.ConnectorProvider
	connectorID models.ConnectorID
	store       Store
	descriptor  models.TaskDescriptor
	publisher   message.Publisher
	messages    *messages.Messages
}

type Store interface {
	UpsertAccounts(ctx context.Context, accounts []*models.Account) error
	UpsertPayments(ctx context.Context, payments []*models.Payment) ([]*models.PaymentID, error)
	UpsertPaymentsAdjustments(ctx context.Context, paymentsAdjustment []*models.PaymentAdjustment) error
	UpsertPaymentsMetadata(ctx context.Context, paymentsMetadata []*models.PaymentMetadata) error
	InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error
	UpdateTaskState(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, state json.RawMessage) error
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, adjustment *models.TransferInitiationAdjustment) error
	UpdateTransferReversalStatus(ctx context.Context, transfer *models.TransferInitiation, transferReversal *models.TransferReversal, adjustment *models.TransferInitiationAdjustment) error
	AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, updatedAt time.Time, metadata map[string]string) error
	AddBankAccountRelatedAccount(ctx context.Context, adjustment *models.BankAccountRelatedAccount) error
}

func NewDefaultIngester(
	provider models.ConnectorProvider,
	connectorID models.ConnectorID,
	descriptor models.TaskDescriptor,
	repo Store,
	publisher message.Publisher,
	messages *messages.Messages,
) *DefaultIngester {
	return &DefaultIngester{
		provider:    provider,
		connectorID: connectorID,
		descriptor:  descriptor,
		store:       repo,
		publisher:   publisher,
		messages:    messages,
	}
}

var _ Ingester = (*DefaultIngester)(nil)
