package ingestion

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
)

func (i *DefaultIngester) UpdateTransferInitiationPaymentsStatus(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, updatedAt time.Time) error {
	adjustment := &models.TransferInitiationAdjustment{
		ID:                   uuid.New(),
		TransferInitiationID: tf.ID,
		CreatedAt:            updatedAt.UTC(),
		Status:               status,
		Error:                errorMessage,
	}

	tf.RelatedAdjustments = append([]*models.TransferInitiationAdjustment{adjustment}, tf.RelatedAdjustments...)

	if err := i.store.UpdateTransferInitiationPaymentsStatus(ctx, tf.ID, paymentID, adjustment); err != nil {
		return err
	}

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

func (i *DefaultIngester) AddTransferInitiationPaymentID(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, updatedAt time.Time) error {
	if paymentID == nil {
		return fmt.Errorf("payment id is nil")
	}

	tf.RelatedPayments = append(tf.RelatedPayments, &models.TransferInitiationPayment{
		TransferInitiationID: tf.ID,
		PaymentID:            *paymentID,
		CreatedAt:            updatedAt.UTC(),
		Status:               models.TransferInitiationStatusProcessing,
	})

	if err := i.store.AddTransferInitiationPaymentID(ctx, tf.ID, paymentID, updatedAt, tf.Metadata); err != nil {
		return err
	}

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
