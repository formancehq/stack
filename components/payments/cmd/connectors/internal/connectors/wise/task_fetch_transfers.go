package wise

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise/client"
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

func taskFetchTransfers(wiseClient *client.Client, profileID uint64) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		logger logging.Logger,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs)
		}()

		transfers, err := wiseClient.GetTransfers(ctx, &client.Profile{
			ID: profileID,
		})
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
			return err
		}

		if len(transfers) == 0 {
			logger.Info("No transfers found")

			return nil
		}

		var (
			// accountBatch ingestion.AccountBatch
			paymentBatch ingestion.PaymentBatch
		)

		for _, transfer := range transfers {
			logger.Info(transfer)

			var rawData json.RawMessage

			rawData, err = json.Marshal(transfer)
			if err != nil {
				return fmt.Errorf("failed to marshal transfer: %w", err)
			}

			precision, ok := supportedCurrenciesWithDecimal[transfer.TargetCurrency]
			if !ok {
				logger.Errorf("currency %s is not supported", transfer.TargetCurrency)
				metricsRegistry.ConnectorCurrencyNotSupported().Add(ctx, 1, metric.WithAttributes(connectorAttrs...))
				continue
			}

			var amount big.Float
			_, ok = amount.SetString(transfer.TargetValue.String())
			if !ok {
				return fmt.Errorf("failed to parse amount %s", transfer.TargetValue.String())
			}

			var amountInt big.Int
			amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

			batchElement := ingestion.PaymentBatchElement{
				Payment: &models.Payment{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: fmt.Sprintf("%d", transfer.ID),
							Type:      models.PaymentTypeTransfer,
						},
						ConnectorID: connectorID,
					},
					CreatedAt:   transfer.CreatedAt,
					Reference:   fmt.Sprintf("%d", transfer.ID),
					ConnectorID: connectorID,
					Type:        models.PaymentTypeTransfer,
					Status:      matchTransferStatus(transfer.Status),
					Scheme:      models.PaymentSchemeOther,
					Amount:      &amountInt,
					Asset:       currency.FormatAsset(supportedCurrenciesWithDecimal, transfer.TargetCurrency),
					RawData:     rawData,
				},
			}

			if transfer.SourceBalanceID != 0 {
				batchElement.Payment.SourceAccountID = &models.AccountID{
					Reference:   fmt.Sprintf("%d", transfer.SourceBalanceID),
					ConnectorID: connectorID,
				}
			}

			if transfer.DestinationBalanceID != 0 {
				batchElement.Payment.DestinationAccountID = &models.AccountID{
					Reference:   fmt.Sprintf("%d", transfer.DestinationBalanceID),
					ConnectorID: connectorID,
				}
			}

			paymentBatch = append(paymentBatch, batchElement)
		}

		if err := ingester.IngestPayments(ctx, connectorID, paymentBatch, struct{}{}); err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
			return err
		}
		metricsRegistry.ConnectorObjects().Add(ctx, int64(len(paymentBatch)), paymentsAttrs)

		return nil
	}
}

func matchTransferStatus(status string) models.PaymentStatus {
	switch status {
	case "incoming_payment_waiting", "processing":
		return models.PaymentStatusPending
	case "funds_converted", "outgoing_payment_sent":
		return models.PaymentStatusSucceeded
	case "bounced_back", "funds_refunded":
		return models.PaymentStatusFailed
	case "cancelled":
		return models.PaymentStatusCancelled
	}

	return models.PaymentStatusOther
}
