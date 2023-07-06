package mangopay

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchTransactions(logger logging.Logger, client *client.Client, userID string) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
	) error {
		logger.Info("Fetching transactions for user", userID)

		for page := 1; ; page++ {
			pagedPayments, err := client.GetTransactions(ctx, userID, page)
			if err != nil {
				return err
			}

			if len(pagedPayments) == 0 {
				break
			}

			if err := ingestBatch(ctx, ingester, pagedPayments); err != nil {
				return err
			}
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
				Amount:    payment.DebitedFunds.Amount,
				Type:      paymentType,
				Status:    matchPaymentStatus(payment.Status),
				Scheme:    models.PaymentSchemeOther,
				Asset:     currency.FormatAsset(payment.DebitedFunds.Currency),
				RawData:   rawData,
			},
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
