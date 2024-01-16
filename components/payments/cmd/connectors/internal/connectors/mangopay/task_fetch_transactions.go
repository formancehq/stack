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
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func taskFetchTransactions(client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("mangopay.taskFetchTransactions")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
			attribute.String("userID", userID),
		)

		if err := fetchTransactions(ctx, client, userID, connectorID, ingester); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchTransactions(
	ctx context.Context,
	client *client.Client,
	userID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
) error {
	for page := 1; ; page++ {
		pagedPayments, err := client.GetTransactions(ctx, userID, page)
		if err != nil {
			return err
		}

		if len(pagedPayments) == 0 {
			break
		}

		if err := ingestBatch(ctx, connectorID, ingester, pagedPayments); err != nil {
			return err
		}
	}

	return nil
}

func ingestBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
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
					ConnectorID: connectorID,
				},
				CreatedAt:     time.Unix(payment.CreationDate, 0),
				Reference:     payment.Id,
				Amount:        &amount,
				InitialAmount: &amount,
				ConnectorID:   connectorID,
				Type:          paymentType,
				Status:        matchPaymentStatus(payment.Status),
				Scheme:        models.PaymentSchemeOther,
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, payment.DebitedFunds.Currency),
				RawData:       rawData,
			},
		}

		if payment.DebitedWalletID != "" {
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference:   payment.DebitedWalletID,
				ConnectorID: connectorID,
			}
		}

		if payment.CreditedWalletID != "" {
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference:   payment.CreditedWalletID,
				ConnectorID: connectorID,
			}
		}

		batch = append(batch, batchElement)
	}

	return ingester.IngestPayments(ctx, connectorID, batch, struct{}{})
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
