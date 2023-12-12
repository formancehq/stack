package atlar

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	atlar_client "github.com/get-momo/atlar-v1-go-client/client"
	"github.com/get-momo/atlar-v1-go-client/client/transactions"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	paymentsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments"))...)
)

func FetchPaymentsTask(config Config, client *atlar_client.Rest, account string) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs)
		}()

		// Pagination works by cursor token.
		params := transactions.GetV1TransactionsParams{
			Context: ctx,
			Limit:   pointer.For(int64(config.ApiConfig.PageSize)),
		}
		for token := ""; ; {
			limit := int64(config.PageSize)
			params.Token = &token
			params.Limit = &limit
			pagedTransactions, err := client.Transactions.GetV1Transactions(&params)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
				return err
			}

			token = pagedTransactions.Payload.NextToken

			if err := ingestPaymentsBatch(ctx, connectorID, ingester, metricsRegistry, pagedTransactions); err != nil {
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedTransactions.Payload.Items)), paymentsAttrs)

			if token == "" {
				break
			}
		}

		return nil
	}
}

func ingestPaymentsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	metricsRegistry metrics.MetricsRegistry,
	pagedTransactions *transactions.GetV1TransactionsOK,
) error {
	batch := ingestion.PaymentBatch{}

	for _, item := range pagedTransactions.Payload.Items {
		raw, err := json.Marshal(item)
		if err != nil {
			return err
		}

		paymentType := determinePaymentType(item)

		itemAmount := item.Amount
		precision := supportedCurrenciesWithDecimal[*itemAmount.Currency]

		var amount big.Float
		_, ok := amount.SetString(*itemAmount.StringValue)
		if !ok {
			return fmt.Errorf("failed to parse amount %s", *itemAmount.StringValue)
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

		createdAt, err := ParseAtlarTimestamp(item.Created)
		if err != nil {
			return err
		}

		paymentId := models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.ID,
				Type:      paymentType,
			},
			ConnectorID: connectorID,
		}

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID:            paymentId,
				Reference:     item.ID,
				Type:          paymentType,
				ConnectorID:   connectorID,
				CreatedAt:     createdAt,
				Status:        determinePaymentStatus(item),
				Scheme:        determinePaymentScheme(item),
				Amount:        &amountInt,
				InitialAmount: &amountInt,
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, *item.Amount.Currency),
				Metadata:      ExtractPaymentMetadata(paymentId, item),
				RawData:       raw,
			},
			Update: true,
		}

		if amountInt.Cmp(big.NewInt(0)) >= 0 {
			// DEBIT
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference:   *item.Account.ID,
				ConnectorID: connectorID,
			}
		} else {
			// CREDIT
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference:   *item.Account.ID,
				ConnectorID: connectorID,
			}
		}

		batch = append(batch, batchElement)
	}

	if err := ingester.IngestPayments(ctx, connectorID, batch, struct{}{}); err != nil {
		return err
	}

	return nil
}

func determinePaymentType(item *atlar_models.Transaction) models.PaymentType {
	if *item.Amount.Value >= 0 {
		return models.PaymentTypePayIn
	} else {
		return models.PaymentTypePayOut
	}
}

func determinePaymentStatus(item *atlar_models.Transaction) models.PaymentStatus {
	if item.Reconciliation.Status == atlar_models.ReconciliationDetailsStatusEXPECTED {
		// A payment initiated by the owner of the accunt through the Atlar API,
		// which was not yet reconciled with a payment from the statement
		return models.PaymentStatusPending
	}
	if item.Reconciliation.Status == atlar_models.ReconciliationDetailsStatusBOOKED {
		// A payment comissioned with the bank, which was not yet reconciled with a
		// payment from the statement
		return models.PaymentStatusPending
	}
	if item.Reconciliation.Status == atlar_models.ReconciliationDetailsStatusRECONCILED {
		return models.PaymentStatusSucceeded
	}
	return models.PaymentStatusOther
}

func determinePaymentScheme(item *atlar_models.Transaction) models.PaymentScheme {
	// item.Characteristics.BankTransactionCode.Domain
	// item.Characteristics.BankTransactionCode.Family
	// TODO: fees and interest -> models.PaymentSchemeOther with additional info on metadata. Will need example transactions for that

	if *item.Amount.Value > 0 {
		return models.PaymentSchemeSepaDebit
	} else if *item.Amount.Value < 0 {
		return models.PaymentSchemeSepaCredit
	}
	return models.PaymentSchemeSepa
}

func ExtractPaymentMetadata(paymentId models.PaymentID, transaction *atlar_models.Transaction) []*models.Metadata {
	result := []*models.Metadata{}
	if transaction.Date != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "date", transaction.Date))
	}
	if transaction.ValueDate != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "valueDate", transaction.ValueDate))
	}
	result = append(result, ComputePaymentMetadata(paymentId, "remittanceInformation/type", *transaction.RemittanceInformation.Type))
	result = append(result, ComputePaymentMetadata(paymentId, "remittanceInformation/value", *transaction.RemittanceInformation.Value))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/domain", transaction.Characteristics.BankTransactionCode.Domain))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/familiy", transaction.Characteristics.BankTransactionCode.Family))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/subfamiliy", transaction.Characteristics.BankTransactionCode.Subfamily))
	result = append(result, ComputePaymentMetadata(paymentId, "btc/description", transaction.Characteristics.BankTransactionCode.Description))
	result = append(result, ComputePaymentMetadataBool(paymentId, "returned", transaction.Characteristics.Returned))
	if transaction.CounterpartyDetails != nil && transaction.CounterpartyDetails.Name != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "counterparty/nane", transaction.CounterpartyDetails.Name))
		if transaction.CounterpartyDetails.ExternalAccount != nil && transaction.CounterpartyDetails.ExternalAccount.Identifier != nil {
			result = append(result, ComputePaymentMetadata(paymentId, "counterparty/bank/bic", transaction.CounterpartyDetails.ExternalAccount.Bank.Bic))
			result = append(result, ComputePaymentMetadata(paymentId, "counterparty/bank/name", transaction.CounterpartyDetails.ExternalAccount.Bank.Name))
			result = append(result, ComputePaymentMetadata(paymentId,
				fmt.Sprintf("counterparty/identifier/%s", transaction.CounterpartyDetails.ExternalAccount.Identifier.Type),
				transaction.CounterpartyDetails.ExternalAccount.Identifier.Number))
		}
	}
	if transaction.Characteristics.Returned {
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/code", transaction.Characteristics.ReturnReason.Code))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/description", transaction.Characteristics.ReturnReason.Description))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/domain", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Domain))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/family", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Family))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/subfamily", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Subfamily))
		result = append(result, ComputePaymentMetadata(paymentId, "returnReason/btc/description", transaction.Characteristics.ReturnReason.OriginalBankTransactionCode.Description))
	}
	result = append(result, ComputePaymentMetadata(paymentId, "reconciliation/status", transaction.Reconciliation.Status))
	result = append(result, ComputePaymentMetadata(paymentId, "reconciliation/transactableId", transaction.Reconciliation.TransactableID))
	result = append(result, ComputePaymentMetadata(paymentId, "reconciliation/transactableType", transaction.Reconciliation.TransactableType))
	if transaction.Characteristics.CurrencyExchange != nil {
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/sourceCurrency", transaction.Characteristics.CurrencyExchange.SourceCurrency))
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/targetCurrency", transaction.Characteristics.CurrencyExchange.TargetCurrency))
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/exchangeRate", transaction.Characteristics.CurrencyExchange.ExchangeRate))
		result = append(result, ComputePaymentMetadata(paymentId, "currencyExchange/unitCurrency", transaction.Characteristics.CurrencyExchange.UnitCurrency))
	}
	if transaction.CounterpartyDetails.MandateReference != "" {
		result = append(result, ComputePaymentMetadata(paymentId, "mandateReference", transaction.CounterpartyDetails.MandateReference))
	}

	return result
}
