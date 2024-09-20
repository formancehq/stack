package ingestion

import (
	"context"
	"math/big"

	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/google/uuid"
)

func (i *DefaultIngester) UpdateTransferReversalStatus(ctx context.Context, tf *models.TransferInitiation, transferReversal *models.TransferReversal) error {
	finalAmount := new(big.Int)
	isFullyReversed := transferReversal.Status == models.TransferReversalStatusProcessed &&
		finalAmount.Sub(tf.Amount, transferReversal.Amount).Cmp(big.NewInt(0)) == 0

	adjustment := &models.TransferInitiationAdjustment{
		ID:                   uuid.New(),
		TransferInitiationID: transferReversal.TransferInitiationID,
		CreatedAt:            transferReversal.UpdatedAt.UTC(),
		Status:               transferReversal.Status.ToTransferInitiationStatus(isFullyReversed),
		Error:                transferReversal.Error,
		Metadata:             transferReversal.Metadata,
	}

	if err := i.store.UpdateTransferReversalStatus(ctx, tf, transferReversal, adjustment); err != nil {
		return err
	}

	tf.RelatedAdjustments = append([]*models.TransferInitiationAdjustment{adjustment}, tf.RelatedAdjustments...)

	if err := i.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			i.messages.NewEventSavedTransferInitiations(tf),
		),
	); err != nil {
		return err
	}

	return nil
}
