package ingestion

import (
	"context"
	"time"

	"github.com/formancehq/payments/internal/app/models"
)

func (i *DefaultIngester) UpdateTransferInitiationStatus(ctx context.Context, id models.TransferInitiationID, status models.TransferInitiationStatus, errorMessage string, updatedAt time.Time) error {
	return i.repo.UpdateTransferInitiationStatus(ctx, id, status, errorMessage, updatedAt)
}

func (i *DefaultIngester) UpdateTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID models.PaymentID, updatedAt time.Time) error {
	return i.repo.UpdateTransferInitiationPaymentID(ctx, id, paymentID, updatedAt)
}
