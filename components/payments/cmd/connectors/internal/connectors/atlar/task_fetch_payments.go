package atlar

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client/transactions"
	atlar_models "github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	paymentsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments"))...)
)

func FetchPaymentsTask(config Config, account string, client *client.Client) task.Task {
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
			params.Token = &token
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

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: item.ID,
						Type:      paymentType,
					},
					ConnectorID: connectorID,
				},
				Reference:   item.ID,
				Type:        paymentType,
				ConnectorID: connectorID,
				Status:      determinePaymentStatus(item),
				Scheme:      determinePaymentScheme(item),
				Amount:      &amountInt,
				Asset:       currency.FormatAsset(supportedCurrenciesWithDecimal, *item.Amount.Currency),
				RawData:     raw,
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
	return models.PaymentStatusOther
}

func determinePaymentScheme(item *atlar_models.Transaction) models.PaymentScheme {
	return models.PaymentSchemeOther
}
