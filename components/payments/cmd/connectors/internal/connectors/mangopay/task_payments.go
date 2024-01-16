package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	bankWireRefPatternRegexp = regexp.MustCompile("[a-zA-Z0-9 ]*")
)

func taskInitiatePayment(mangopayClient *client.Client, transferID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		span := trace.SpanFromContext(ctx)
		span.SetName("mangopay.taskInitiatePayment")
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

		if err := initiatePayment(ctx, mangopayClient, transfer, connectorID, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func initiatePayment(
	ctx context.Context,
	mangopayClient *client.Client,
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

	// No need to modify the amount since it's already in the correct format
	// and precision (checked before during API call)
	var curr string
	curr, _, err = currency.GetCurrencyAndPrecisionFromAsset(supportedCurrenciesWithDecimal, transfer.Asset)
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
	mangopayClient *client.Client,
	transferID string,
	pID string,
	attempt int,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		paymentID := models.MustPaymentIDFromString(pID)
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		span := trace.SpanFromContext(ctx)
		span.SetName("mangopay.taskUpdatePaymentStatus")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("transferID", transferID),
			attribute.String("paymentID", paymentID.String()),
			attribute.String("reference", transferInitiationID.Reference),
		)

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, false)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := udpatePaymentStatus(ctx, mangopayClient, transfer, paymentID, connectorID, attempt, ingester, scheduler, storageReader); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func udpatePaymentStatus(
	ctx context.Context,
	mangopayClient *client.Client,
	transfer *models.TransferInitiation,
	paymentID *models.PaymentID,
	connectorID models.ConnectorID,
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
		resp, err = mangopayClient.GetWalletTransfer(ctx, paymentID.Reference)
		if err != nil {
			return err
		}

		status = resp.Status
		resultMessage = resp.ResultMessage
	case models.TransferInitiationTypePayout:
		var resp *client.PayoutResponse
		resp, err = mangopayClient.GetPayout(ctx, paymentID.Reference)
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
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}
	case "SUCCEEDED":
		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessed, "", time.Now())
		if err != nil {
			return err
		}

		return nil
	case "FAILED":
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
