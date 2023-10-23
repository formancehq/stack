package ingestion

import (
	"context"
	"time"

	"github.com/formancehq/payments/internal/models"
)

func (i *DefaultIngester) UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error {
	return i.repo.UpdateTransferInitiationPaymentsStatus(ctx, id, paymentID, status, errorMessage, attempts, updatedAt)
}

func (i *DefaultIngester) AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, updatedAt time.Time) error {
	return i.repo.AddTransferInitiationPaymentID(ctx, id, paymentID, updatedAt)
}
