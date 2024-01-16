package wise

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func taskFetchTransfers(wiseClient *client.Client, profileID uint64) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("wise.taskFetchTransfers")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("profileID", strconv.FormatUint(profileID, 10)),
		)

		if err := fetchTransfers(ctx, wiseClient, profileID, connectorID, scheduler, ingester); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchTransfers(
	ctx context.Context,
	wiseClient *client.Client,
	profileID uint64,
	connectorID models.ConnectorID,
	scheduler task.Scheduler,
	ingester ingestion.Ingester,
) error {
	transfers, err := wiseClient.GetTransfers(ctx, &client.Profile{
		ID: profileID,
	})
	if err != nil {
		return err
	}

	if len(transfers) == 0 {
		return nil
	}

	var (
		// accountBatch ingestion.AccountBatch
		paymentBatch ingestion.PaymentBatch
	)

	for _, transfer := range transfers {

		var rawData json.RawMessage

		rawData, err = json.Marshal(transfer)
		if err != nil {
			return fmt.Errorf("failed to marshal transfer: %w", err)
		}

		precision, ok := supportedCurrenciesWithDecimal[transfer.TargetCurrency]
		if !ok {
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
		return err
	}

	return nil
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
