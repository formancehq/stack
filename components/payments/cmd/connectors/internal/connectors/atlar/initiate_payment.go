package atlar

import (
	"errors"
	"math/big"
	"strings"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/get-momo/atlar-v1-go-client/client/credit_transfers"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
)

func initiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation, config Config) error {
	err := validateTransferInitiation(transfer)
	if err != nil {
		return err
	}

	client := createAtlarClient(&config)
	detachedCtx, _ := contextutil.Detached(ctx.Context())

	currency, precision, err := currency.GetCurrencyAndPrecisionFromAsset(supportedCurrenciesWithDecimal, transfer.Asset)
	if err != nil {
		return err
	}

	paymentSchemeType := "SCT" // SEPA Credit Transfer
	remittanceInformationType := "UNSTRUCTURED"
	remittanceInformationValue := transfer.Description

	createPaymentRequest := atlar_models.CreatePaymentRequest{
		SourceAccountID:              &transfer.SourceAccount.Reference,
		DestinationExternalAccountID: &transfer.DestinationAccount.Reference,
		Amount: &atlar_models.AmountInput{
			Currency:    &currency,
			Value:       transfer.Amount.Int64(),
			StringValue: amountToString(*transfer.Amount, precision),
		},
		Date:              TimeToAtlarTimestamp(&transfer.ScheduledAt),
		ExternalID:        transfer.ID.Reference, // TODO: don't know if that's correct
		PaymentSchemeType: &paymentSchemeType,
		RemittanceInformation: &atlar_models.RemittanceInformation{
			Type:  &remittanceInformationType,
			Value: &remittanceInformationValue,
		},
	}

	postCounterpartiesParams := credit_transfers.PostV1CreditTransfersParams{
		Context:        detachedCtx,
		CreditTransfer: &createPaymentRequest,
	}
	_, err = client.CreditTransfers.PostV1CreditTransfers(&postCounterpartiesParams)
	if err != nil {
		return err
	}

	return connectors.ErrNotImplemented
}

func validateTransferInitiation(transfer *models.TransferInitiation) error {
	if transfer.Description == "" {
		return errors.New("description of transfer initiation can not be empty")
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
