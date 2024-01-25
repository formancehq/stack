package mangopay

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchTransactionsState struct {
	// Mangopay only allows us to sort/filter by creation date.
	// So in order to have every last updates of transactions, we need to
	// keep track of the first transaction with created status in order to
	// refetch all transactions created after this one.
	// Example:
	// - SUCCEEDED
	// - FAILED
	// - CREATED -> We want to keep track of the creation date of this transaction since we want its updates
	// - SUCCEEDED
	// - CREATED
	// - SUCCEEDED
	FirstCreatedTransactionCreationDate time.Time `json:"first_created_transaction_creation_date"`
}

func taskFetchTransactions(client *client.Client, walletsID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"mangopay.taskFetchTransactions",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("walletsID", walletsID),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchTransactionsState{})

		newState, err := fetchTransactions(ctx, client, walletsID, connectorID, ingester, state)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, newState); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchTransactions(
	ctx context.Context,
	client *client.Client,
	walletsID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	state fetchTransactionsState,
) (fetchTransactionsState, error) {
	newState := fetchTransactionsState{}

	var firstCreatedCreationDate time.Time
	var lastCreationDate time.Time
	for page := 1; ; page++ {
		pagedPayments, err := client.GetTransactions(ctx, walletsID, page, pageSize, state.FirstCreatedTransactionCreationDate)
		if err != nil {
			return fetchTransactionsState{}, err
		}

		if len(pagedPayments) == 0 {
			break
		}

		batch := ingestion.PaymentBatch{}
		for _, payment := range pagedPayments {
			batchElement, err := processPayment(ctx, connectorID, payment)
			if err != nil {
				return fetchTransactionsState{}, err
			}

			if batchElement.Payment != nil {
				// State update
				if firstCreatedCreationDate.IsZero() &&
					batchElement.Payment.Status == models.PaymentStatusPending {
					firstCreatedCreationDate = batchElement.Payment.CreatedAt
				}

				lastCreationDate = batchElement.Payment.CreatedAt
			}

			batch = append(batch, batchElement)
		}

		err = ingester.IngestPayments(ctx, batch)
		if err != nil {
			return fetchTransactionsState{}, err
		}

		if len(pagedPayments) < pageSize {
			break
		}
	}

	newState.FirstCreatedTransactionCreationDate = firstCreatedCreationDate
	if newState.FirstCreatedTransactionCreationDate.IsZero() {
		// No new created payments, let's set the last creation date to the last
		// transaction we fetched.
		newState.FirstCreatedTransactionCreationDate = lastCreationDate
	}

	return newState, nil
}

func processPayment(
	ctx context.Context,
	connectorID models.ConnectorID,
	payment *client.Payment,
) (ingestion.PaymentBatchElement, error) {
	rawData, err := json.Marshal(payment)
	if err != nil {
		return ingestion.PaymentBatchElement{}, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	paymentType := matchPaymentType(payment.Type)
	paymentStatus := matchPaymentStatus(payment.Status)

	var amount big.Int
	_, ok := amount.SetString(payment.DebitedFunds.Amount.String(), 10)
	if !ok {
		return ingestion.PaymentBatchElement{}, fmt.Errorf("failed to parse amount %s", payment.DebitedFunds.Amount.String())
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
			Status:        paymentStatus,
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

	return batchElement, nil
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
