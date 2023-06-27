package mangopay

import (
	"context"
	"encoding/json"
	"fmt"

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

		transactions, err := client.GetAllTransactions(ctx, userID)
		if err != nil {
			return err
		}

		batch := ingestion.PaymentBatch{}
		for _, transaction := range transactions {
			logger.Info(transaction)

			rawData, err := json.Marshal(transaction)
			if err != nil {
				return fmt.Errorf("failed to marshal transaction: %w", err)
			}

			batchElement := ingestion.PaymentBatchElement{
				Payment: &models.Payment{
					CreatedAt: transaction.CreationDate,
					Reference: transaction.Id,
					Amount:    int64(transaction.DebitedFunds.Amount * 100),
					Type:      matchPaymentType(transaction.Type),
					Status:    matchPaymentStatus(transaction.Status),
					Scheme:    models.PaymentSchemeOther,
					Asset:     models.PaymentAsset(transaction.DebitedFunds.Currency + "/2"),
					RawData:   rawData,
				},
			}

			batch = append(batch, batchElement)
		}

		return ingester.IngestPayments(ctx, batch, struct{}{})
	}
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
