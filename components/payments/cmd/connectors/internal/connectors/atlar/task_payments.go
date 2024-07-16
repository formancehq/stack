package atlar

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/get-momo/atlar-v1-go-client/client/credit_transfers"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
	"go.opentelemetry.io/otel/attribute"
)

func InitiatePaymentTask(config Config, client *client.Client, transferID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		storageReader storage.Reader,
		ingester ingestion.Ingester,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		ctx, span := connectors.StartSpan(
			ctx,
			"atlar.taskInitiatePayment",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("transferID", transferID),
			attribute.String("reference", transferInitiationID.Reference),
		)
		defer span.End()

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			return err
		}

		var paymentID *models.PaymentID
		defer func() {
			if err != nil {
				otel.RecordError(span, err)
				ctx, cancel := contextutil.Detached(ctx)
				defer cancel()
				if err := ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusFailed, err.Error(), time.Now()); err != nil {
					otel.RecordError(span, err)
				}
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

		currency, precision, err := currency.GetCurrencyAndPrecisionFromAsset(supportedCurrenciesWithDecimal, transfer.Asset)
		if err != nil {
			return err
		}

		paymentSchemeType := "SCT" // SEPA Credit Transfer
		remittanceInformationType := "UNSTRUCTURED"
		remittanceInformationValue := transfer.Description
		amount := atlar_models.AmountInput{
			Currency:    &currency,
			Value:       transfer.Amount.Int64(),
			StringValue: amountToString(*transfer.Amount, precision),
		}
		date := transfer.ScheduledAt
		if date.IsZero() {
			date = time.Now()
		}
		dateString := date.Format(time.DateOnly)

		createPaymentRequest := atlar_models.CreatePaymentRequest{
			SourceAccountID:              &transfer.SourceAccount.Reference,
			DestinationExternalAccountID: &transfer.DestinationAccount.Reference,
			Amount:                       &amount,
			Date:                         &dateString,
			ExternalID:                   serializeAtlarPaymentExternalID(transfer.ID.Reference, transfer.CountRetries()),
			PaymentSchemeType:            &paymentSchemeType,
			RemittanceInformation: &atlar_models.RemittanceInformation{
				Type:  &remittanceInformationType,
				Value: &remittanceInformationValue,
			},
		}

		requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
		defer cancel()
		var postCreditTransferResponse *credit_transfers.PostV1CreditTransfersCreated
		postCreditTransferResponse, err = client.PostV1CreditTransfers(requestCtx, &createPaymentRequest)
		if err != nil {
			return err
		}

		paymentID = &models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: postCreditTransferResponse.Payload.Reconciliation.ExpectedTransactionID,
				Type:      models.PaymentTypePayOut,
			},
			ConnectorID: connectorID,
		}

		var taskDescriptor models.TaskDescriptor
		taskDescriptor, err = models.EncodeTaskDescriptor(TaskDescriptor{
			Name:       fmt.Sprintf("Update transfer initiation status of transfer %s", transfer.ID.String()),
			Key:        taskNameUpdatePaymentStatus,
			TransferID: transfer.ID.String(),
			PaymentID:  paymentID.String(),
			Attempt:    1,
		})
		if err != nil {
			return err
		}

		err = scheduler.Schedule(ctx, taskDescriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		return nil
	}
}

func ValidateTransferInitiation(transfer *models.TransferInitiation) error {
	if transfer == nil {
		return errors.New("transfer cannot be nil")
	}
	if transfer.Type.String() != "PAYOUT" {
		return errors.New("this connector only supports type PAYOUT")
	}
	return nil
}

func UpdatePaymentStatusTask(
	config Config,
	client *client.Client,
	transferID string,
	stringPaymentID string,
	attempt int,
) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		storageReader storage.Reader,
		ingester ingestion.Ingester,
	) error {
		paymentID := models.MustPaymentIDFromString(stringPaymentID)
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		ctx, span := connectors.StartSpan(
			ctx,
			"atlar.taskUpdatePaymentStatus",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("transferID", transferID),
			attribute.String("paymentID", stringPaymentID),
			attribute.Int("attempt", attempt),
			attribute.String("reference", transferInitiationID.Reference),
		)
		defer span.End()

		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
		defer cancel()
		getCreditTransferResponse, err := client.GetV1CreditTransfersGetByExternalIDExternalID(
			requestCtx,
			serializeAtlarPaymentExternalID(transfer.ID.Reference, transfer.CountRetries()),
		)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		status := getCreditTransferResponse.Payload.Status
		// Status docs: https://docs.atlar.com/docs/payment-details#payment-states--events
		switch status {
		case "CREATED", "APPROVED", "PENDING_SUBMISSION", "SENT", "PENDING_AT_BANK", "ACCEPTED", "EXECUTED":
			taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:       fmt.Sprintf("Update transfer initiation status of transfer %s", transfer.ID.String()),
				Key:        taskNameUpdatePaymentStatus,
				TransferID: transfer.ID.String(),
				PaymentID:  paymentID.String(),
				Attempt:    attempt + 1,
			})
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			err = scheduler.Schedule(ctx, taskDescriptor, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_IN_DURATION,
				Duration:       config.TransferInitiationStatusPollingPeriod.Duration,
				RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				otel.RecordError(span, err)
				return err
			}
			return nil

		case "RECONCILED":
			err = ingestAtlarTransaction(ctx,
				ingester,
				connectorID,
				taskID,
				client,
				getCreditTransferResponse.Payload.Reconciliation.BookedTransactionID,
			)
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			paymentID = &models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: getCreditTransferResponse.Payload.Reconciliation.BookedTransactionID,
					Type:      models.PaymentTypePayOut,
				},
				ConnectorID: connectorID,
			}

			err = ingester.UpdateTransferInitiationPayment(ctx, transfer, paymentID, models.TransferInitiationStatusProcessed, "", time.Now())
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			return nil

		case "REJECTED", "FAILED", "RETURNED":
			err = ingester.UpdateTransferInitiationPaymentsStatus(
				ctx, transfer, paymentID, models.TransferInitiationStatusFailed,
				fmt.Sprintf("paymant initiation status is \"%s\"", status), time.Now(),
			)
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			return nil

		default:
			err := fmt.Errorf(
				"unknown status \"%s\" encountered while fetching payment initiation status of payment \"%s\"",
				status, getCreditTransferResponse.Payload.ID,
			)
			otel.RecordError(span, err)
			return err
		}
	}
}

func amountToString(amount big.Int, precision int) string {
	raw := amount.String()
	if precision < 0 {
		precision = 0
	}
	insertPosition := len(raw) - precision
	if insertPosition <= 0 {
		return "0." + strings.Repeat("0", -insertPosition) + raw
	}
	return raw[:insertPosition] + "." + raw[insertPosition:]
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

func serializeAtlarPaymentExternalID(ID string, attempts int) string {
	return fmt.Sprintf("%s_%d", ID, attempts)
}

var deserializeAtlarPaymentExternalIDRegex = regexp.MustCompile(`^([^\_]+)_([0-9]+)$`)

func deserializeAtlarPaymentExternalID(serialized string) (string, int, error) {
	var attempts int

	// Find matches in the input string
	matches := deserializeAtlarPaymentExternalIDRegex.FindStringSubmatch(serialized)
	if matches == nil || len(matches) != 3 {
		return "", 0, errors.New("cannot deserialize malformed externalID")
	}

	parsed, err := fmt.Sscanf(matches[2], "%d", &attempts)
	if err != nil {
		return "", 0, errors.New("cannot deserialize malformed externalID")
	}
	if parsed != 1 {
		return "", 0, errors.New("cannot deserialize malformed externalID")
	}
	return matches[1], attempts, nil
}

func ingestAtlarTransaction(
	ctx context.Context,
	ingester ingestion.Ingester,
	connectorID models.ConnectorID,
	taskID models.TaskID,
	client *client.Client,
	transactionId string,
) error {
	ctx, span := connectors.StartSpan(
		ctx,
		"atlar.taskUpdatePaymentStatus.ingestAtlarTransaction",
		attribute.String("connectorID", connectorID.String()),
		attribute.String("taskID", taskID.String()),
		attribute.String("transactionID", transactionId),
	)
	defer span.End()

	requestCtx, cancel := contextutil.DetachedWithTimeout(ctx, 30*time.Second)
	defer cancel()
	transactionResponse, err := client.GetV1TransactionsID(requestCtx, transactionId)
	if err != nil {
		otel.RecordError(span, err)
		return err
	}

	batchElement, err := atlarTransactionToPaymentBatchElement(ctx, connectorID, taskID, transactionResponse.Payload, client)
	if err != nil {
		otel.RecordError(span, err)
		return err
	}
	if batchElement == nil {
		return nil
	}

	batch := ingestion.PaymentBatch{*batchElement}

	err = ingester.IngestPayments(ctx, batch)
	if err != nil {
		otel.RecordError(span, err)
		return err
	}

	return nil
}
