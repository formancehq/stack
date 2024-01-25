package bankingcircle

import (
	"context"
	"encoding/json"
	"math"
	"math/big"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func taskFetchPayments(
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"bankingcircle.taskFetchPayments",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		if err := fetchPayments(ctx, client, connectorID, ingester); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchPayments(
	ctx context.Context,
	client *client.Client,
	connectorID models.ConnectorID,
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

	for _, paymentEl := range payments {
		raw, err := json.Marshal(paymentEl)
		if err != nil {
			return err
		}

		paymentType := matchPaymentType(paymentEl.Classification)

		precision, ok := supportedCurrenciesWithDecimal[paymentEl.Transfer.Amount.Currency]
		if !ok {
			continue
		}

		var amount big.Float
		amount.SetFloat64(paymentEl.Transfer.Amount.Amount)
		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

		batchElement := ingestion.PaymentBatchElement{
			Payment: &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: paymentEl.PaymentID,
						Type:      paymentType,
					},
					ConnectorID: connectorID,
				},
				Reference:     paymentEl.PaymentID,
				Type:          paymentType,
				ConnectorID:   connectorID,
				Status:        matchPaymentStatus(paymentEl.Status),
				Scheme:        models.PaymentSchemeOther,
				Amount:        &amountInt,
				InitialAmount: &amountInt,
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, paymentEl.Transfer.Amount.Currency),
				RawData:       raw,
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

	if err := ingester.IngestPayments(ctx, batch); err != nil {
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
