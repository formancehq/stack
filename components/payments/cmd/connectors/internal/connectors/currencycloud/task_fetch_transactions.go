package currencycloud

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	paymentsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments"))...)
)

func taskFetchTransactions(logger logging.Logger, client *client.Client, config Config) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		return ingestTransactions(ctx, logger, connectorID, client, ingester, metricsRegistry)
	}
}

func ingestTransactions(ctx context.Context, logger logging.Logger, connectorID models.ConnectorID,
	client *client.Client, ingester ingestion.Ingester, metricsRegistry metrics.MetricsRegistry,
) error {
	now := time.Now()
	defer func() {
		metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs)
	}()

	page := 1
	for {
		if page < 0 {
			break
		}

		logger.Info("Fetching transactions")

		transactions, nextPage, err := client.GetTransactions(ctx, page)
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
			return err
		}

		page = nextPage

		batch := ingestion.PaymentBatch{}

		for _, transaction := range transactions {
			logger.Info(transaction)

			var amount big.Float
			_, ok := amount.SetString(transaction.Amount)
			if !ok {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
				return fmt.Errorf("failed to parse amount %s", transaction.Amount)
			}

			var amountInt big.Int
			amount.Mul(&amount, big.NewFloat(100)).Int(&amountInt)

			var rawData json.RawMessage

			rawData, err = json.Marshal(transaction)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
				return fmt.Errorf("failed to marshal transaction: %w", err)
			}

			paymentType := matchTransactionType(transaction.Type)

			batchElement := ingestion.PaymentBatchElement{
				Payment: &models.Payment{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: transaction.ID,
							Type:      paymentType,
						},
						ConnectorID: connectorID,
					},
					Reference:   transaction.ID,
					Type:        paymentType,
					ConnectorID: connectorID,
					Status:      matchTransactionStatus(transaction.Status),
					Scheme:      models.PaymentSchemeOther,
					Amount:      &amountInt,
					Asset:       models.Asset(fmt.Sprintf("%s/2", transaction.Currency)),
					RawData:     rawData,
				},
			}

			switch paymentType {
			case models.PaymentTypePayOut:
				batchElement.Payment.SourceAccountID = &models.AccountID{
					Reference:   transaction.AccountID,
					ConnectorID: connectorID,
				}
			default:
				batchElement.Payment.DestinationAccountID = &models.AccountID{
					Reference:   transaction.AccountID,
					ConnectorID: connectorID,
				}
			}

			batch = append(batch, batchElement)
		}

		err = ingester.IngestPayments(ctx, connectorID, batch, struct{}{})
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
			return err
		}
		metricsRegistry.ConnectorObjects().Add(ctx, int64(len(batch)), paymentsAttrs)
	}

	return nil
}

func matchTransactionType(transactionType string) models.PaymentType {
	switch transactionType {
	case "credit":
		return models.PaymentTypePayOut
	case "debit":
		return models.PaymentTypePayIn
	}

	return models.PaymentTypeOther
}

func matchTransactionStatus(transactionStatus string) models.PaymentStatus {
	switch transactionStatus {
	case "completed":
		return models.PaymentStatusSucceeded
	case "pending":
		return models.PaymentStatusPending
	case "deleted":
		return models.PaymentStatusFailed
	}

	return models.PaymentStatusOther
}
