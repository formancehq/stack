package wise

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
)

func taskInitiatePayment(
	wiseClient *client.Client,
	transferID string,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		span := trace.SpanFromContext(ctx)
		span.SetName("wise.taskInitiatePayment")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("transferID", transferID),
			attribute.String("reference", transferInitiationID.Reference),
		)

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := initiatePayment(ctx, wiseClient, transfer, connectorID, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func initiatePayment(
	ctx context.Context,
	wiseClient *client.Client,
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

	if transfer.SourceAccount == nil {
		err = errors.New("missing source account")
		return err
	}

	if transfer.SourceAccount.Type == models.AccountTypeExternal {
		err = errors.New("payin not implemented: source account must be an internal account")
		return err
	}

	profileID, ok := transfer.SourceAccount.Metadata["profile_id"]
	if !ok || profileID == "" {
		err = errors.New("missing user_id in source account metadata")
		return err
	}

	var curr string
	var precision int
	curr, precision, err = currency.GetCurrencyAndPrecisionFromAsset(supportedCurrenciesWithDecimal, transfer.Asset)
	if err != nil {
		return err
	}

	amount := big.NewFloat(0).SetInt(transfer.Amount)
	amount = amount.Quo(amount, big.NewFloat(math.Pow(10, float64(precision))))

	quote, err := wiseClient.CreateQuote(ctx, profileID, curr, amount)
	if err != nil {
		return err
	}

	var connectorPaymentID uint64
	var paymentType models.PaymentType
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	switch transfer.DestinationAccount.Type {
	case models.AccountTypeInternal:
		// Transfer between internal accounts
		destinationAccount, err := strconv.ParseUint(transfer.DestinationAccount.Metadata["profile_id"], 10, 64)
		if err != nil {
			return err
		}

		var resp *client.Transfer
		resp, err = wiseClient.CreateTransfer(ctx, quote, destinationAccount, fmt.Sprintf("%s_%d", transfer.ID.Reference, len(transfer.RelatedAdjustments)))
		if err != nil {
			return err
		}

		connectorPaymentID = resp.ID
		paymentType = models.PaymentTypeTransfer
	case models.AccountTypeExternal:
		// Payout to an external account

		destinationAccount, err := strconv.ParseUint(transfer.DestinationAccount.Reference, 10, 64)
		if err != nil {
			return err
		}

		var resp *client.Payout
		resp, err = wiseClient.CreatePayout(ctx, quote, destinationAccount, fmt.Sprintf("%s_%d", transfer.ID.Reference, len(transfer.RelatedAdjustments)))
		if err != nil {
			return err
		}

		connectorPaymentID = resp.ID
		paymentType = models.PaymentTypePayOut
	}

	paymentID = &models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: strconv.FormatUint(connectorPaymentID, 10),
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

func taskUpdatePaymentStatus(
	wiseClient *client.Client,
	transferID string,
	pID string,
	attempt int,
) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		paymentID := models.MustPaymentIDFromString(pID)
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		span := trace.SpanFromContext(ctx)
		span.SetName("wise.taskUpdatePaymentStatus")
		span.SetAttributes(
			attribute.String("transferID", transferID),
			attribute.String("paymentID", pID),
			attribute.Int("attempt", attempt),
			attribute.String("reference", transferInitiationID.Reference),
		)

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, false)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := updatePaymentStatus(ctx, wiseClient, transfer, paymentID, attempt, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func updatePaymentStatus(
	ctx context.Context,
	wiseClient *client.Client,
	transfer *models.TransferInitiation,
	paymentID *models.PaymentID,
	attempt int,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	storageReader storage.Reader,
) error {
	var err error
	var status string
	switch transfer.Type {
	case models.TransferInitiationTypeTransfer:
		var resp *client.Transfer
		resp, err = wiseClient.GetTransfer(ctx, paymentID.Reference)
		if err != nil {
			return err
		}

		status = resp.Status
	case models.TransferInitiationTypePayout:
		var resp *client.Payout
		resp, err = wiseClient.GetPayout(ctx, paymentID.Reference)
		if err != nil {
			return err
		}

		status = resp.Status
	}

	switch status {
	case "incoming_payment_waiting",
		"incoming_payment_initiated",
		"processing",
		"funds_converted",
		"bounced_back",
		"unknown":
		taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:       "Update transfer initiation status",
			Key:        taskNameUpdatePaymentStatus,
			TransferID: transfer.ID.String(),
			Attempt:    attempt + 1,
		})
		if err != nil {
			return err
		}

		err = scheduler.Schedule(ctx, taskDescriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_IN_DURATION,
			Duration:       2 * time.Minute,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}
	case "outgoing_payment_sent", "funds_refunded":
		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessed, "", time.Now())
		if err != nil {
			return err
		}

		return nil
	case "charged_back", "cancelled":
		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusFailed, "", time.Now())
		if err != nil {
			return err
		}

		return nil
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
