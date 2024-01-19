package stripe

import (
	"context"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func reversePaymentTask(transferReversalID string, stripeClient *client.DefaultClient) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		reversalID := models.MustTransferReversalIDFromString(transferReversalID)

		span := trace.SpanFromContext(ctx)
		span.SetName("stripe.reversePaymentTask")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("transferReversalID", transferReversalID),
			attribute.String("reference", reversalID.Reference),
		)

		transferReversal, err := getTransferReversal(ctx, storageReader, reversalID)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		transfer, err := getTransfer(ctx, storageReader, transferReversal.TransferInitiationID, true)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := reversePayment(ctx, stripeClient, transfer, transferReversal, connectorID, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func reversePayment(
	ctx context.Context,
	stripeClient *client.DefaultClient,
	transfer *models.TransferInitiation,
	transferReversal *models.TransferReversal,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	storageReader storage.Reader,
) error {
	var err error
	defer func() {
		if err != nil {
			ctx, cancel := contextutil.Detached(ctx)
			defer cancel()

			transferReversal.Status = models.TransferReversalStatusFailed
			transferReversal.Error = err.Error()
			transferReversal.UpdatedAt = time.Now().UTC()

			_ = ingester.UpdateTransferReversalStatus(ctx, transfer, transferReversal)
		}
	}()

	c := client.Client(stripeClient)
	// If source account is nil, or equal to root (which is a special
	// account we create for stripe for the balance platform), we don't need
	// to set the stripe account.
	if transfer.SourceAccount != nil && transfer.SourceAccount.Reference != rootAccountReference {
		c = c.ForAccount(transfer.SourceAccountID.Reference)
	}

	_, err = c.ReverseTransfer(ctx, &client.CreateTransferReversalRequest{
		TransferID:  transfer.ID.Reference,
		Amount:      transferReversal.Amount.Int64(),
		Description: transferReversal.Description,
		Metadata:    transferReversal.Metadata,
	})
	if err != nil {
		return err
	}

	transferReversal.Status = models.TransferReversalStatusProcessed
	if err = ingester.UpdateTransferReversalStatus(ctx, transfer, transferReversal); err != nil {
		return err
	}

	return nil
}

func getTransferReversal(
	ctx context.Context,
	reader storage.Reader,
	transferReversalID models.TransferReversalID,
) (*models.TransferReversal, error) {
	transferReversal, err := reader.GetTransferReversal(ctx, transferReversalID)
	if err != nil {
		return nil, err
	}

	return transferReversal, nil
}
