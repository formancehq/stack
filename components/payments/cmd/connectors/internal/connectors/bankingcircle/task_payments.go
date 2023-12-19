package bankingcircle

import (
	"context"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskInitiatePayment(logger logging.Logger, bankingCircleClient *client.Client, transferID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
	) error {
		logger.Info("initiate payment for transfer-initiation %s", transferID)

		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)
		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			return err
		}

		var paymentID *models.PaymentID
		defer func() {
			if err != nil {
				ctx, cancel := contextutil.Detached(ctx)
				defer cancel()
				if err := ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusFailed, err.Error(), transfer.Attempts, time.Now()); err != nil {
					logger.Error("failed to update transfer initiation status: %v", err)
				}
			}
		}()

		err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessing, "", transfer.Attempts, time.Now())
		if err != nil {
			return err
		}

		logger.Info("initiate payment between", transfer.SourceAccountID, " and %s", transfer.DestinationAccountID)

		if transfer.SourceAccount == nil {
			err = errors.New("no source account provided")
			return err
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

		amount := big.NewFloat(0).SetInt(transfer.Amount)
		amount = amount.Quo(amount, big.NewFloat(math.Pow(10, float64(precision))))

		var sourceAccount *client.Account
		sourceAccount, err = bankingCircleClient.GetAccount(ctx, transfer.SourceAccountID.Reference)
		if err != nil {
			return err
		}
		if len(sourceAccount.AccountIdentifiers) == 0 {
			err = errors.New("no source account identifiers provided")
			return err
		}

		var destinationAccount *client.Account
		destinationAccount, err = bankingCircleClient.GetAccount(ctx, transfer.DestinationAccountID.Reference)
		if err != nil {
			return err
		}
		if len(destinationAccount.AccountIdentifiers) == 0 {
			err = errors.New("no destination account identifiers provided")
			return err
		}

		var connectorPaymentID string
		var paymentType models.PaymentType
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		switch transfer.DestinationAccount.Type {
		case models.AccountTypeInternal:
			// Transfer between internal accounts
			var resp *client.PaymentResponse
			resp, err = bankingCircleClient.InitiateTransferOrPayouts(ctx, &client.PaymentRequest{
				IdempotencyKey:         transfer.ID.Reference,
				RequestedExecutionDate: transfer.ScheduledAt,
				DebtorAccount: client.PaymentAccount{
					Account:              sourceAccount.AccountIdentifiers[0].Account,
					FinancialInstitution: sourceAccount.AccountIdentifiers[0].FinancialInstitution,
					Country:              sourceAccount.AccountIdentifiers[0].Country,
				},
				DebtorReference:    transfer.Description,
				CurrencyOfTransfer: curr,
				Amount: struct {
					Currency string     "json:\"currency\""
					Amount   *big.Float "json:\"amount\""
				}{
					Currency: curr,
					Amount:   amount,
				},
				ChargeBearer: "SHA",
				CreditorAccount: &client.PaymentAccount{
					Account:              destinationAccount.AccountIdentifiers[0].Account,
					FinancialInstitution: destinationAccount.AccountIdentifiers[0].FinancialInstitution,
					Country:              destinationAccount.AccountIdentifiers[0].Country,
				},
			})
			if err != nil {
				return err
			}

			connectorPaymentID = resp.PaymentID
			paymentType = models.PaymentTypeTransfer
		case models.AccountTypeExternal:
			// Payout to an external account
			var resp *client.PaymentResponse
			resp, err = bankingCircleClient.InitiateTransferOrPayouts(ctx, &client.PaymentRequest{
				IdempotencyKey:         transfer.ID.Reference,
				RequestedExecutionDate: transfer.ScheduledAt,
				DebtorAccount: client.PaymentAccount{
					Account:              sourceAccount.AccountIdentifiers[0].Account,
					FinancialInstitution: sourceAccount.AccountIdentifiers[0].FinancialInstitution,
					Country:              sourceAccount.AccountIdentifiers[0].Country,
				},
				DebtorReference:    transfer.Description,
				CurrencyOfTransfer: curr,
				Amount: struct {
					Currency string     "json:\"currency\""
					Amount   *big.Float "json:\"amount\""
				}{
					Currency: curr,
					Amount:   amount,
				},
				ChargeBearer: "SHA",
				CreditorAccount: &client.PaymentAccount{
					Account:              destinationAccount.AccountIdentifiers[0].Account,
					FinancialInstitution: destinationAccount.AccountIdentifiers[0].FinancialInstitution,
					Country:              destinationAccount.AccountIdentifiers[0].Country,
				},
			})
			if err != nil {
				return err
			}

			connectorPaymentID = resp.PaymentID
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
}

func taskUpdatePaymentStatus(
	logger logging.Logger,
	bankingCircleClient *client.Client,
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
		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, false)
		if err != nil {
			return err
		}
		logger.Info("attempt: ", attempt, " fetching status of ", pID)

		var status string
		switch transfer.Type {
		case models.TransferInitiationTypeTransfer:
			var resp *client.StatusResponse
			resp, err = bankingCircleClient.GetPaymentStatus(ctx, paymentID.Reference)
			if err != nil {
				return err
			}

			status = resp.Status
		case models.TransferInitiationTypePayout:
			var resp *client.StatusResponse
			resp, err = bankingCircleClient.GetPaymentStatus(ctx, paymentID.Reference)
			if err != nil {
				return err
			}

			status = resp.Status
		}

		switch status {
		case "PendingApproval", "PendingProcessing", "Hold", "Approved", "ScaPending":
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
		case "Processed":
			err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusProcessed, "", transfer.Attempts, time.Now())
			if err != nil {
				return err
			}

			return nil
		case "Unknown", "ScaExpired", "ScaFailed", "MissingFunding",
			"PendingCancellation", "PendingCancellationApproval", "DeclinedByApprover",
			"Rejected", "Cancelled", "Reversed", "ScaDeclined":
			err = ingester.UpdateTransferInitiationPaymentsStatus(ctx, transfer, paymentID, models.TransferInitiationStatusFailed, "", transfer.Attempts, time.Now())
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
