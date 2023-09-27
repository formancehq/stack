package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/storage"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	initiateTransferAttrs    = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "initiate_transfer"))...)
	initiatePayoutAttrs      = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "initiate_payout"))...)
	bankWireRefPatternRegexp = regexp.MustCompile("[a-zA-Z0-9 ]*")
)

func taskInitiatePayment(logger logging.Logger, mangopayClient *client.Client, transferID string) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info("initiate payment for transfer-initiation %s", transferID)

		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		attrs := metric.WithAttributes(connectorAttrs...)
		var err error
		defer func() {
			if err != nil {
				ctx, cancel := contextutil.Detached(ctx)
				defer cancel()
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, attrs)
				if err := ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusFailed, err.Error(), time.Now()); err != nil {
					logger.Error("failed to update transfer initiation status: %v", err)
				}
			}
		}()

		err = ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusProcessing, "", time.Now())
		if err != nil {
			return err
		}

		var transfer *models.TransferInitiation
		transfer, err = getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			return err
		}

		attrs = initiateTransferAttrs
		if transfer.Type == models.TransferInitiationTypePayout {
			attrs = initiatePayoutAttrs
		}

		logger.Info("initiate payment between", transfer.SourceAccountID, " and %s", transfer.DestinationAccountID)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), attrs)
		}()

		if transfer.SourceAccount == nil {
			err = errors.New("no source account provided")
			return err
		}

		if transfer.SourceAccount.Type == models.AccountTypeExternal {
			err = errors.New("payin not implemented: source account must be an internal account")
			return err
		}

		userID, ok := transfer.SourceAccount.Metadata["user_id"]
		if !ok || userID == "" {
			err = errors.New("missing user_id in source account metadata")
			return err
		}

		var curr string
		curr, _, err = currency.GetCurrencyAndPrecisionFromAsset(transfer.Asset)
		if err != nil {
			return err
		}

		bankWireRef := ""
		if len(transfer.Description) <= 12 && bankWireRefPatternRegexp.MatchString(transfer.Description) {
			bankWireRef = transfer.Description
		}

		var connectorPaymentID string
		var paymentType models.PaymentType
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		switch transfer.DestinationAccount.Type {
		case models.AccountTypeInternal:
			// Transfer between internal accounts
			var resp *client.TransferResponse
			resp, err = mangopayClient.InitiateWalletTransfer(ctx, &client.TransferRequest{
				AuthorID: userID,
				DebitedFunds: client.Funds{
					Currency: curr,
					Amount:   json.Number(transfer.Amount.String()),
				},
				Fees: client.Funds{
					Currency: curr,
					Amount:   "0",
				},
				DebitedWalletID:  transfer.SourceAccountID.Reference,
				CreditedWalletID: transfer.DestinationAccountID.Reference,
			})
			if err != nil {
				return err
			}

			connectorPaymentID = resp.ID
			paymentType = models.PaymentTypeTransfer
		case models.AccountTypeExternal:
			// Payout to an external account
			var resp *client.PayoutResponse
			resp, err = mangopayClient.InitiatePayout(ctx, &client.PayoutRequest{
				AuthorID: userID,
				DebitedFunds: client.Funds{
					Currency: curr,
					Amount:   json.Number(transfer.Amount.String()),
				},
				Fees: client.Funds{
					Currency: curr,
					Amount:   json.Number("0"),
				},
				DebitedWalletID: transfer.SourceAccountID.Reference,
				BankAccountID:   transfer.DestinationAccountID.Reference,
				BankWireRef:     bankWireRef,
			})
			if err != nil {
				return err
			}

			connectorPaymentID = resp.ID
			paymentType = models.PaymentTypePayOut
		}
		metricsRegistry.ConnectorObjects().Add(ctx, 1, attrs)

		err = ingester.UpdateTransferInitiationPaymentID(ctx, transferInitiationID, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: connectorPaymentID,
				Type:      paymentType,
			},
			Provider: models.ConnectorProviderMangopay,
		}, time.Now())
		if err != nil {
			return err
		}

		taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:       "Update transfer initiation status",
			Key:        taskNameUpdatePaymentStatus,
			TransferID: transfer.ID.String(),
			Attempt:    1,
		})
		if err != nil {
			return err
		}

		ctx, _ = contextutil.DetachedWithTimeout(ctx, 10*time.Second)
		err = scheduler.Schedule(ctx, taskDescriptor, models.TaskSchedulerOptions{
			// We want to polling every c.cfg.PollingPeriod.Duration seconds the users
			// and their transactions.
			ScheduleOption: models.OPTIONS_RUN_NOW,
			// No need to restart this task, since the connector is not existing or
			// was uninstalled previously, the task does not exists in the database
			Restart: true,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		return nil
	}
}

var (
	updateTransferAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "update_transfer"))...)
	updatePayoutAttrs   = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "update_payout"))...)
)

func taskUpdatePaymentStatus(
	logger logging.Logger,
	mangopayClient *client.Client,
	transferID string,
	attempt int,
) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)
		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, false)
		if err != nil {
			return err
		}
		logger.Info("attempt: ", attempt, " fetching status of ", transfer.PaymentID)

		attrs := updateTransferAttrs
		if transfer.Type == models.TransferInitiationTypePayout {
			attrs = updatePayoutAttrs
		}

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), attrs)
		}()

		defer func() {
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, attrs)
			}
		}()

		var status string
		var resultMessage string
		switch transfer.Type {
		case models.TransferInitiationTypeTransfer:
			var resp *client.TransferResponse
			resp, err = mangopayClient.GetWalletTransfer(ctx, transfer.PaymentID.Reference)
			if err != nil {
				return err
			}

			status = resp.Status
			resultMessage = resp.ResultMessage
		case models.TransferInitiationTypePayout:
			var resp *client.PayoutResponse
			resp, err = mangopayClient.GetPayout(ctx, transfer.PaymentID.Reference)
			if err != nil {
				return err
			}

			status = resp.Status
			resultMessage = resp.ResultMessage
		}

		switch status {
		case "CREATED":
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
				Restart:        true,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}
		case "SUCCEEDED":
			err = ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusProcessed, "", time.Now())
			if err != nil {
				return err
			}

			return nil
		case "FAILED":
			err = ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusFailed, resultMessage, time.Now())
			if err != nil {
				return err
			}

			return nil
		}

		return nil
	}
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
		if transfer.SourceAccountID.Reference != "" {
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
