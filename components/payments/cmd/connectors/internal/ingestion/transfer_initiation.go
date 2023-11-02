package ingestion

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

func (i *DefaultIngester) UpdateTransferInitiationPaymentsStatus(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error {
	tf.Status = status
	tf.Error = errorMessage
	tf.Attempts = attempts
	tf.UpdatedAt = updatedAt

	if err := i.repo.UpdateTransferInitiationPaymentsStatus(ctx, tf.ID, paymentID, tf.Status, tf.Error, tf.Attempts, tf.UpdatedAt); err != nil {
		return err
	}

	if err := i.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			messages.NewEventSavedTransferInitiations(tf),
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

	tf.RelatedPayments = append(tf.RelatedPayments, &models.TransferInitiationPayments{
		TransferInitiationID: tf.ID,
		PaymentID:            *paymentID,
		CreatedAt:            updatedAt,
		Status:               models.TransferInitiationStatusProcessing,
	})

	if err := i.repo.AddTransferInitiationPaymentID(ctx, tf.ID, paymentID, updatedAt); err != nil {
		return err
	}

	if err := i.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			messages.NewEventSavedTransferInitiations(tf),
		),
	); err != nil {
		return err
	}

	return nil
}
