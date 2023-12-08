package atlar

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	atlar_client "github.com/get-momo/atlar-v1-go-client/client"
	"github.com/get-momo/atlar-v1-go-client/client/credit_transfers"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
)

func InitiatePaymentTask(config Config, client *atlar_client.Rest, transferID string) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		storageReader storage.Reader,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)
		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			return err
		}

		err = validateTransferInitiation(transfer)
		if err != nil {
			return err
		}

		// client := createAtlarClient(&config)

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
		logger.WithContext(ctx).Infof("date is %s", date.Format(time.RFC3339Nano))

		createPaymentRequest := atlar_models.CreatePaymentRequest{
			SourceAccountID:              &transfer.SourceAccount.Reference,
			DestinationExternalAccountID: &transfer.DestinationAccount.Reference,
			Amount:                       &amount,
			Date:                         &dateString,
			ExternalID:                   transfer.ID.Reference, // TODO: don't know if that's correct
			PaymentSchemeType:            &paymentSchemeType,
			RemittanceInformation: &atlar_models.RemittanceInformation{
				Type:  &remittanceInformationType,
				Value: &remittanceInformationValue,
			},
		}

		logger.Debug(createPaymentRequest, "createPaymentRequest")

		postCreditTransfersParams := credit_transfers.PostV1CreditTransfersParams{
			Context:        ctx,
			CreditTransfer: &createPaymentRequest,
		}
		_, err = client.CreditTransfers.PostV1CreditTransfers(&postCreditTransfersParams)
		if err != nil {
			return err
		}

		// return connectors.ErrNotImplemented
		return nil
	}
}

func validateTransferInitiation(transfer *models.TransferInitiation) error {
	if transfer == nil {
		return errors.New("transfer cannot be nil")
	}
	if transfer.Description == "" {
		return errors.New("description of transfer initiation can not be empty")
	}
	if transfer.Type.String() != "TRANSFER" {
		return errors.New("this connector only supports type TRANSFER")
	}
	return nil
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
