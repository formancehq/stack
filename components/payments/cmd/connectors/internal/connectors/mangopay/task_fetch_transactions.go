package mangopay

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
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

func taskFetchTransactions(logger logging.Logger, client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info("Fetching transactions for user", userID)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs)
		}()

		for page := 1; ; page++ {
			pagedPayments, err := client.GetTransactions(ctx, userID, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
				return err
			}

			if len(pagedPayments) == 0 {
				break
			}

			if err := ingestBatch(ctx, ingester, pagedPayments); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedPayments)), paymentsAttrs)
		}

		return nil
	}
}

func ingestBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	payments []*client.Payment,
) error {
	batch := ingestion.PaymentBatch{}
	for _, payment := range payments {
		rawData, err := json.Marshal(payment)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}

		paymentType := matchPaymentType(payment.Type)

		var amount big.Int
		_, ok := amount.SetString(payment.DebitedFunds.Amount.String(), 10)
		if !ok {
			return fmt.Errorf("failed to parse amount %s", payment.DebitedFunds.Amount.String())
		}

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: payment.Id,
						Type:      paymentType,
					},
					Provider: models.ConnectorProviderMangopay,
				},
				CreatedAt: time.Unix(payment.CreationDate, 0),
				Reference: payment.Id,
				Amount:    &amount,
				Type:      paymentType,
				Status:    matchPaymentStatus(payment.Status),
				Scheme:    models.PaymentSchemeOther,
				Asset:     currency.FormatAsset(payment.DebitedFunds.Currency),
				RawData:   rawData,
			},
		}

		if payment.DebitedWalletID != "" {
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference: payment.DebitedWalletID,
				Provider:  models.ConnectorProviderMangopay,
			}
		}

		if payment.CreditedWalletID != "" {
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference: payment.CreditedWalletID,
				Provider:  models.ConnectorProviderMangopay,
			}
		}

		batch = append(batch, batchElement)
	}

	return ingester.IngestPayments(ctx, batch, struct{}{})
}

func matchPaymentType(paymentType string) models.PaymentType {
	switch paymentType {
	case "PAYIN":
		return models.PaymentTypePayIn
	case "PAYOUT":
		return models.PaymentTypePayOut
	case "TRANSFER":
		return models.PaymentTypeTransfer
	}

	return models.PaymentTypeOther
}

func matchPaymentStatus(paymentStatus string) models.PaymentStatus {
	switch paymentStatus {
	case "CREATED":
		return models.PaymentStatusPending
	case "SUCCEEDED":
		return models.PaymentStatusSucceeded
	case "FAILED":
		return models.PaymentStatusFailed
	}

	return models.PaymentStatusOther
}
