package moneycorp

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/connectors/moneycorp/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	pageSize = 100
)

var (
	paymentsAttrs = append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments"))
)

func taskFetchTransactions(logger logging.Logger, client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info("Fetching transactions for account", accountID)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), metric.WithAttributes(paymentsAttrs...))
		}()

		for page := 0; ; page++ {
			pagedTransactions, err := client.GetTransactions(ctx, accountID, page, pageSize)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, metric.WithAttributes(paymentsAttrs...))
				return err
			}

			if len(pagedTransactions) == 0 {
				break
			}

			if err := ingestBatch(ctx, ingester, pagedTransactions); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, metric.WithAttributes(paymentsAttrs...))
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedTransactions)), metric.WithAttributes(paymentsAttrs...))

			if len(pagedTransactions) < pageSize {
				break
			}
		}

		return nil
	}
}

func ingestBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	transactions []*client.Transaction,
) error {
	batch := ingestion.PaymentBatch{}
	for _, transaction := range transactions {
		rawData, err := json.Marshal(transaction)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}

		paymentType, shouldBeRecorded := matchPaymentType(transaction.Attributes.Type, transaction.Attributes.Direction)
		if !shouldBeRecorded {
			continue
		}

		createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", transaction.Attributes.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to parse transaction date: %w", err)
		}

		var amount big.Float
		_, ok := amount.SetString(transaction.Attributes.Amount.String())
		if !ok {
			return fmt.Errorf("failed to parse amount %s", transaction.Attributes.Amount.String())
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(currency.GetPrecision(transaction.Attributes.Currency))))).Int(&amountInt)

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: transaction.ID,
						Type:      paymentType,
					},
					Provider: models.ConnectorProviderMoneycorp,
				},
				CreatedAt: createdAt,
				Reference: transaction.ID,
				Amount:    &amountInt,
				Asset:     currency.FormatAsset(transaction.Attributes.Currency),
				Type:      paymentType,
				Status:    models.PaymentStatusSucceeded,
				Scheme:    models.PaymentSchemeOther,
				RawData:   rawData,
			},
		}

		switch paymentType {
		case models.PaymentTypePayIn:
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference: strconv.Itoa(int(transaction.Attributes.AccountID)),
				Provider:  models.ConnectorProviderMoneycorp,
			}
		case models.PaymentTypePayOut:
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference: strconv.Itoa(int(transaction.Attributes.AccountID)),
				Provider:  models.ConnectorProviderMoneycorp,
			}
		default:
			if transaction.Attributes.Direction == "Debit" {
				batchElement.Payment.SourceAccountID = &models.AccountID{
					Reference: strconv.Itoa(int(transaction.Attributes.AccountID)),
					Provider:  models.ConnectorProviderMoneycorp,
				}
			} else {
				batchElement.Payment.DestinationAccountID = &models.AccountID{
					Reference: strconv.Itoa(int(transaction.Attributes.AccountID)),
					Provider:  models.ConnectorProviderMoneycorp,
				}
			}
		}

		batch = append(batch, batchElement)
	}

	return ingester.IngestPayments(ctx, batch, struct{}{})
}

func matchPaymentType(transactionType string, transactionDirection string) (models.PaymentType, bool) {
	switch transactionType {
	case "Transfer":
		return models.PaymentTypeTransfer, true
	case "Payment", "Exchange", "Charge", "Refund":
		switch transactionDirection {
		case "Debit":
			return models.PaymentTypePayOut, true
		case "Credit":
			return models.PaymentTypePayIn, true
		}
	}

	return models.PaymentTypeOther, false
}
