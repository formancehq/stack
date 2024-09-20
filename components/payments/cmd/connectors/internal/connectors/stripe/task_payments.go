package stripe

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/formancehq/go-libs/contextutil"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/stripe/stripe-go/v72"
	"go.opentelemetry.io/otel/attribute"
)

const (
	transferIDKey string = "transfer_id"
)

func initiatePaymentTask(transferID string, stripeClient *client.DefaultClient) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		ctx, span := connectors.StartSpan(
			ctx,
			"stripe.initiatePaymentTask",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("transferID", transferID),
			attribute.String("reference", transferInitiationID.Reference),
		)
		defer span.End()

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := initiatePayment(ctx, stripeClient, transfer, connectorID, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func initiatePayment(
	ctx context.Context,
	stripeClient *client.DefaultClient,
	transfer *models.TransferInitiation,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	storageReader storage.Reader,
) error {
	var err error
	var paymentID *models.PaymentID
	defer func() {
		if err != nil {
			ctx, cancel := contextutil.Detached(ctx)
			defer cancel()
			_ = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusFailed, err.Error(), time.Now())
		}
	}()

	err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessing, "", time.Now())
	if err != nil {
		return err
	}

	if transfer.SourceAccount != nil {
		if transfer.SourceAccount.Type == models.AccountTypeExternal {
			err = errors.New("payin not implemented: source account must be an internal account")
			return err
		}
	}

	var curr string
	curr, _, err = currency.GetCurrencyAndPrecisionFromAsset(supportedCurrenciesWithDecimal, transfer.Asset)
	if err != nil {
		return err
	}

	c := client.Client(stripeClient)
	// If source account is nil, or equal to root (which is a special
	// account we create for stripe for the balance platform), we don't need
	// to set the stripe account.
	if transfer.SourceAccount != nil && transfer.SourceAccount.Reference != rootAccountReference {
		c = c.ForAccount(transfer.SourceAccountID.Reference)
	}

	var connectorPaymentID string
	var paymentType models.PaymentType
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	switch transfer.DestinationAccount.Type {
	case models.AccountTypeInternal:
		// Transfer between internal accounts
		var resp *stripe.Transfer
		resp, err = c.CreateTransfer(ctx, &client.CreateTransferRequest{
			IdempotencyKey: fmt.Sprintf("%s_%d", transfer.ID.Reference, len(transfer.RelatedAdjustments)),
			Amount:         transfer.Amount.Int64(),
			Currency:       curr,
			Destination:    transfer.DestinationAccountID.Reference,
			Description:    transfer.Description,
		})
		if err != nil {
			return err
		}

		if transfer.Metadata == nil {
			transfer.Metadata = make(map[string]string)
		}
		transfer.Metadata[transferIDKey] = resp.ID
		connectorPaymentID = resp.BalanceTransaction.ID
		paymentType = models.PaymentTypeTransfer
	case models.AccountTypeExternal:
		// Payout to an external account
		var resp *stripe.Payout
		resp, err = c.CreatePayout(ctx, &client.CreatePayoutRequest{
			IdempotencyKey: fmt.Sprintf("%s_%d", transfer.ID.Reference, len(transfer.RelatedAdjustments)),
			Amount:         transfer.Amount.Int64(),
			Currency:       curr,
			Destination:    transfer.DestinationAccountID.Reference,
			Description:    transfer.Description,
		})
		if err != nil {
			return err
		}

		connectorPaymentID = resp.BalanceTransaction.ID
		paymentType = models.PaymentTypePayOut
	}

	paymentID = &models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: connectorPaymentID,
			Type:      paymentType,
		},
		ConnectorID: connectorID,
	}

	err = ingester.AddTransferInitiationPaymentID(ctx, transfer, paymentID, time.Now())
	if err != nil {
		return err
	}

	taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name:       "Update transfer initiation status",
		Key:        taskNameUpdatePaymentStatus,
		TransferID: transfer.ID.String(),
		PaymentID:  paymentID.String(),
		Attempt:    1,
	})
	if err != nil {
		return err
	}

	ctx, _ = contextutil.DetachedWithTimeout(ctx, 10*time.Second)
	err = scheduler.Schedule(ctx, taskDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	})
	if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
		return err
	}

	return nil
}

func updatePaymentStatusTask(
	transferID string,
	pID string,
	attempt int,
	stripeClient *client.DefaultClient,
) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		paymentID := models.MustPaymentIDFromString(pID)
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		ctx, span := connectors.StartSpan(
			ctx,
			"stripe.updatePaymentStatusTask",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("transferID", transferID),
			attribute.String("paymentID", pID),
			attribute.Int("attempt", attempt),
			attribute.String("reference", transferInitiationID.Reference),
		)
		defer span.End()

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, false)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := updatePaymentStatus(ctx, stripeClient, transfer, paymentID, attempt, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func updatePaymentStatus(
	ctx context.Context,
	stripeClient *client.DefaultClient,
	transfer *models.TransferInitiation,
	paymentID *models.PaymentID,
	attempt int,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	storageReader storage.Reader,
) error {
	var err error
	var status stripe.PayoutFailureCode
	var resultMessage string
	switch transfer.Type {
	case models.TransferInitiationTypeTransfer:
		// Nothing to do
		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessed, "", time.Now())
		if err != nil {
			return err
		}

		return nil

	case models.TransferInitiationTypePayout:
		var resp *stripe.Payout
		resp, err = stripeClient.GetPayout(ctx, paymentID.Reference)
		if err != nil {
			return err
		}

		status = resp.FailureCode
		resultMessage = resp.FailureMessage
	}

	if status == "" {
		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessed, "", time.Now())
		if err != nil {
			return err
		}

		return nil
	}

	err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusFailed, resultMessage, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func getTransfer(
	ctx context.Context,
	reader storage.Reader,
	transferID models.TransferInitiationID,
	expand bool,
) (*models.TransferInitiation, error) {
	transfer, err := reader.ReadTransferInitiation(ctx, transferID)
	if err != nil {
		return nil, err
	}

	if expand {
		if transfer.SourceAccountID != nil {
			sourceAccount, err := reader.GetAccount(ctx, transfer.SourceAccountID.String())
			if err != nil {
				return nil, err
			}
			transfer.SourceAccount = sourceAccount
		}

		destinationAccount, err := reader.GetAccount(ctx, transfer.DestinationAccountID.String())
		if err != nil {
			return nil, err
		}
		transfer.DestinationAccount = destinationAccount
	}

	return transfer, nil
}
