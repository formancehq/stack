package moneycorp

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func taskInitiatePayment(moneycorpClient *client.Client, transferID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		span := trace.SpanFromContext(ctx)
		span.SetName("moneycorp.taskInitiatePayment")
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

		if err := initiatePayment(ctx, moneycorpClient, transfer, connectorID, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func initiatePayment(
	ctx context.Context,
	moneycorpClient *client.Client,
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
		err = errors.New("no source account provided")
	}

	if transfer.SourceAccount.Type == models.AccountTypeExternal {
		err = errors.New("payin not implemented: source account must be an internal account")
		return err
	}

	var curr string
	var precision int
	curr, precision, err = currency.GetCurrencyAndPrecisionFromAsset(supportedCurrenciesWithDecimal, transfer.Asset)
	if err != nil {
		return err
	}

	transferAmount := big.NewFloat(0).SetInt(transfer.Amount)
	amount, accuracy := transferAmount.Quo(transferAmount, big.NewFloat(math.Pow(10, float64(precision)))).Float64()
	if accuracy != big.Exact {
		return errors.New("amount is not accurate, psp does not support big ints")
	}

	var connectorPaymentID string
	var paymentType models.PaymentType
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	switch transfer.DestinationAccount.Type {
	case models.AccountTypeInternal:
		// Transfer between internal accounts
		var resp *client.TransferResponse
		resp, err = moneycorpClient.InitiateTransfer(ctx, &client.TransferRequest{
			SourceAccountID:    transfer.SourceAccountID.Reference,
			IdempotencyKey:     fmt.Sprintf("%s_%d", transfer.ID.Reference, len(transfer.RelatedAdjustments)),
			ReceivingAccountID: transfer.DestinationAccountID.Reference,
			TransferAmount:     amount,
			TransferCurrency:   curr,
			TransferReference:  transfer.Description,
			ClientReference:    transfer.Description,
		})
		if err != nil {
			return err
		}

		connectorPaymentID = resp.ID
		paymentType = models.PaymentTypeTransfer
	case models.AccountTypeExternal:
		// Payout to an external account
		var resp *client.PayoutResponse
		resp, err = moneycorpClient.InitiatePayout(ctx, &client.PayoutRequest{
			SourceAccountID:  transfer.SourceAccountID.Reference,
			IdempotencyKey:   fmt.Sprintf("%s_%d", transfer.ID.Reference, len(transfer.RelatedAdjustments)),
			RecipientID:      transfer.DestinationAccountID.Reference,
			PaymentAmount:    amount,
			PaymentCurrency:  curr,
			PaymentMethgod:   "Standard",
			PaymentReference: transfer.Description,
			ClientReference:  transfer.Description,
		})
		if err != nil {
			return err
		}

		connectorPaymentID = resp.ID
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

func taskUpdatePaymentStatus(
	moneycorpClient *client.Client,
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
		span.SetName("moneycorp.taskUpdatePaymentStatus")
		span.SetAttributes(
			attribute.String("transferID", transferID),
			attribute.String("paymentID", pID),
			attribute.Int("attempt", attempt),
			attribute.String("reference", transferInitiationID.Reference),
		)

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := updatePaymentStatus(ctx, moneycorpClient, transfer, paymentID, attempt, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func updatePaymentStatus(
	ctx context.Context,
	moneycorpClient *client.Client,
	transfer *models.TransferInitiation,
	paymentID *models.PaymentID,
	attempt int,
	ingester ingestion.Ingester,
	scheduler task.Scheduler,
	storageReader storage.Reader,
) error {
	var err error
	var status string
	var resultMessage string
	switch transfer.Type {
	case models.TransferInitiationTypeTransfer:
		var resp *client.TransferResponse
		resp, err = moneycorpClient.GetTransfer(ctx, transfer.SourceAccount.Reference, paymentID.Reference)
		if err != nil {
			return err
		}

		status = resp.Attributes.TransferStatus
	case models.TransferInitiationTypePayout:
		var resp *client.PayoutResponse
		resp, err = moneycorpClient.GetPayout(ctx, transfer.SourceAccount.Reference, paymentID.Reference)
		if err != nil {
			return err
		}

		status = resp.Attributes.PaymentStatus
	}

	switch status {
	case "Awaiting Dispatch":
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
	case "Cleared", "Sent":
		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessed, "", time.Now())
		if err != nil {
			return err
		}

		return nil
	case "Unauthorised", "Failed", "Cancelled":
		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusFailed, resultMessage, time.Now())
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
