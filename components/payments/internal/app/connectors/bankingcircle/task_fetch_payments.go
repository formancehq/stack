package bankingcircle

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/app/connectors/bankingcircle/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchPayments(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		for page := 1; ; page++ {
			pagedPayments, err := client.GetPayments(ctx, page)
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

	for _, paymentEl := range payments {
		raw, err := json.Marshal(paymentEl)
		if err != nil {
			return err
		}

		paymentType := matchPaymentType(paymentEl.Classification)

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
				Amount:    int64(paymentEl.Transfer.Amount.Amount * 100),
				Asset:     models.PaymentAsset(paymentEl.Transfer.Amount.Currency + "/2"),
				RawData:   raw,
			},
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
