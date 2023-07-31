package modulr

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/modulr/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
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
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs...)
		}()

		transactions, err := client.GetTransactions(accountID)
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs...)
			return err
		}

		batch := ingestion.PaymentBatch{}

		for _, transaction := range transactions {
			logger.Info(transaction)

			rawData, err := json.Marshal(transaction)
			if err != nil {
				return fmt.Errorf("failed to marshal transaction: %w", err)
			}

			paymentType := matchTransactionType(transaction.Type)

			var amount big.Float
			_, ok := amount.SetString(transaction.Amount.String())
			if !ok {
				return fmt.Errorf("failed to parse amount %s", transaction.Amount.String())
			}

			var amountInt big.Int
			amount.Mul(&amount, big.NewFloat(100)).Int(&amountInt)

			batchElement := ingestion.PaymentBatchElement{
				Payment: &models.Payment{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: transaction.ID,
							Type:      paymentType,
						},
						Provider: models.ConnectorProviderModulr,
					},
					Reference: transaction.ID,
					Type:      paymentType,
					Status:    models.PaymentStatusSucceeded,
					Scheme:    models.PaymentSchemeOther,
					Amount:    &amountInt,
					Asset:     models.PaymentAsset(fmt.Sprintf("%s/2", transaction.Account.Currency)),
					RawData:   rawData,
				},
			}

			switch paymentType {
			case models.PaymentTypePayIn:
				batchElement.Payment.DestinationAccountID = &models.AccountID{
					Reference: accountID,
					Provider:  models.ConnectorProviderModulr,
				}
			case models.PaymentTypePayOut:
				batchElement.Payment.SourceAccountID = &models.AccountID{
					Reference: accountID,
					Provider:  models.ConnectorProviderModulr,
				}
			default:
				if transaction.Credit {
					batchElement.Payment.DestinationAccountID = &models.AccountID{
						Reference: accountID,
						Provider:  models.ConnectorProviderModulr,
					}
				} else {
					batchElement.Payment.SourceAccountID = &models.AccountID{
						Reference: accountID,
						Provider:  models.ConnectorProviderModulr,
					}
				}
			}

			batch = append(batch, batchElement)
		}

		if err := ingester.IngestPayments(ctx, batch, struct{}{}); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs...)
			return err
		}

		metricsRegistry.ConnectorObjects().Add(ctx, int64(len(transactions)), paymentsAttrs...)
		return nil
	}
}

func matchTransactionType(transactionType string) models.PaymentType {
	if transactionType == "PI_REV" ||
		transactionType == "PO_REV" ||
		transactionType == "ADHOC" ||
		transactionType == "INT_INTERC" {
		return models.PaymentTypeOther
	}

	if strings.HasPrefix(transactionType, "PI_") {
		return models.PaymentTypePayIn
	}

	if strings.HasPrefix(transactionType, "PO_") {
		return models.PaymentTypePayOut
	}

	return models.PaymentTypeOther
}
