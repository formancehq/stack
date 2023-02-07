package ingestion

import (
	"context"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/google/uuid"
)

func (i *DefaultIngester) GetTransfer(ctx context.Context, transferID uuid.UUID) (models.Transfer, error) {
	return i.repo.GetTransfer(ctx, transferID)
}

func (i *DefaultIngester) UpdateTransferStatus(ctx context.Context, transferID uuid.UUID,
	status models.TransferStatus, reference, err string,
) error {
	return i.repo.UpdateTransferStatus(ctx, transferID, status, reference, err)
}
