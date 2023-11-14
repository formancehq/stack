package bankingcircle

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle/client"
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

func taskFetchPayments(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs)
		}()

		for page := 1; ; page++ {
			pagedPayments, err := client.GetPayments(ctx, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
				return err
			}

			if len(pagedPayments) == 0 {
				break
			}

			if err := ingestBatch(ctx, connectorID, ingester, pagedPayments); err != nil {
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
	connectorID models.ConnectorID,
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
		amount.SetFloat64(paymentEl.Transfer.Amount.Amount)

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(100)).Int(&amountInt)

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: paymentEl.PaymentID,
						Type:      paymentType,
					},
					ConnectorID: connectorID,
				},
				Reference:   paymentEl.PaymentID,
				Type:        paymentType,
				ConnectorID: connectorID,
				Status:      matchPaymentStatus(paymentEl.Status),
				Scheme:      models.PaymentSchemeOther,
				Amount:      &amountInt,
				Asset:       models.Asset(paymentEl.Transfer.Amount.Currency + "/2"),
				RawData:     raw,
			},
		}

		if paymentEl.DebtorInformation.AccountID != "" {
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference:   paymentEl.DebtorInformation.AccountID,
				ConnectorID: connectorID,
			}
		}

		if paymentEl.CreditorInformation.AccountID != "" {
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference:   paymentEl.CreditorInformation.AccountID,
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
	case "Own":
		return models.PaymentTypeTransfer
	}

	return models.PaymentTypeOther
}
