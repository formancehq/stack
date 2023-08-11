package bankingcircle

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/bankingcircle/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	paymentsAttrs = append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments"))
)

func taskFetchPayments(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), metric.WithAttributes(paymentsAttrs...))
		}()

		for page := 1; ; page++ {
			pagedPayments, err := client.GetPayments(ctx, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, metric.WithAttributes(paymentsAttrs...))
				return err
			}

			if len(pagedPayments) == 0 {
				break
			}

			if err := ingestBatch(ctx, ingester, pagedPayments); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, metric.WithAttributes(paymentsAttrs...))
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedPayments)), metric.WithAttributes(paymentsAttrs...))
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

	for _, paymentEl := range payments {
		raw, err := json.Marshal(paymentEl)
		if err != nil {
			return err
		}

		paymentType := matchPaymentType(paymentEl.Classification)

		var amount big.Float
		_, ok := amount.SetString(paymentEl.Transfer.Amount.Amount.String())
		if !ok {
			return fmt.Errorf("failed to parse amount %s", paymentEl.Transfer.Amount.Amount.String())
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(100)).Int(&amountInt)

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: paymentEl.TransactionReference,
						Type:      paymentType,
					},
					Provider: models.ConnectorProviderBankingCircle,
				},
				Reference: paymentEl.TransactionReference,
				Type:      paymentType,
				Status:    matchPaymentStatus(paymentEl.Status),
				Scheme:    models.PaymentSchemeOther,
				Amount:    &amountInt,
				Asset:     models.Asset(paymentEl.Transfer.Amount.Currency + "/2"),
				RawData:   raw,
			},
		}

		if paymentEl.DebtorInformation.AccountID != "" {
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference: paymentEl.DebtorInformation.AccountID,
				Provider:  models.ConnectorProviderBankingCircle,
			}
		}

		if paymentEl.CreditorInformation.AccountID != "" {
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference: paymentEl.CreditorInformation.AccountID,
				Provider:  models.ConnectorProviderBankingCircle,
			}
		}

		batch = append(batch, batchElement)
	}

	if err := ingester.IngestPayments(ctx, batch, struct{}{}); err != nil {
		return err
	}

	return nil
}

func matchPaymentStatus(paymentStatus string) models.PaymentStatus {
	switch paymentStatus {
	case "Processed":
		return models.PaymentStatusSucceeded
	// On MissingFunding - the payment is still in progress.
	// If there will be funds available within 10 days - the payment will be processed.
	// Otherwise - it will be cancelled.
	case "PendingProcessing", "MissingFunding":
		return models.PaymentStatusPending
	case "Rejected", "Cancelled", "Reversed", "Returned":
		return models.PaymentStatusFailed
	}

	return models.PaymentStatusOther
}

func matchPaymentType(paymentType string) models.PaymentType {
	switch paymentType {
	case "Incoming":
		return models.PaymentTypePayIn
	case "Outgoing":
		return models.PaymentTypePayOut
	}

	return models.PaymentTypeOther
}
