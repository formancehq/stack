package currencycloud

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchTransactionsState struct {
	LastUpdatedAt time.Time
}

func taskFetchTransactions(client *client.Client, config Config) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"currencycloud.taskFetchTransactions",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchTransactionsState{})

		newState, err := ingestTransactions(ctx, connectorID, client, ingester, state)
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

func ingestTransactions(
	ctx context.Context,
	connectorID models.ConnectorID,
	client *client.Client,
	ingester ingestion.Ingester,
	state fetchTransactionsState,
) (fetchTransactionsState, error) {
	newState := fetchTransactionsState{
		LastUpdatedAt: state.LastUpdatedAt,
	}

	page := 1
	for {
		if page < 0 {
			break
		}

		transactions, nextPage, err := client.GetTransactions(ctx, page, state.LastUpdatedAt)
		if err != nil {
			return fetchTransactionsState{}, err
		}

		page = nextPage

		batch := ingestion.PaymentBatch{}

		for _, transaction := range transactions {
			switch transaction.UpdatedAt.Compare(state.LastUpdatedAt) {
			case -1, 0:
				continue
			default:
			}

			precision, ok := supportedCurrenciesWithDecimal[transaction.Currency]
			if !ok {
				continue
			}

			var amount big.Float
			_, ok = amount.SetString(transaction.Amount)
			if !ok {
				return fetchTransactionsState{}, fmt.Errorf("failed to parse amount %s", transaction.Amount)
			}
			var amountInt big.Int
			amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

			var rawData json.RawMessage

			rawData, err = json.Marshal(transaction)
			if err != nil {
				return fetchTransactionsState{}, fmt.Errorf("failed to marshal transaction: %w", err)
			}

			paymentType := matchTransactionType(transaction.Type)

			batchElement := ingestion.PaymentBatchElement{
				Payment: &models.Payment{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: transaction.ID,
							Type:      paymentType,
						},
						ConnectorID: connectorID,
					},
					Reference:     transaction.ID,
					CreatedAt:     transaction.CreatedAt,
					Type:          paymentType,
					ConnectorID:   connectorID,
					Status:        matchTransactionStatus(transaction.Status),
					Scheme:        models.PaymentSchemeOther,
					Amount:        &amountInt,
					InitialAmount: &amountInt,
					Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transaction.Currency),
					RawData:       rawData,
				},
			}

			switch paymentType {
			case models.PaymentTypePayOut:
				batchElement.Payment.SourceAccountID = &models.AccountID{
					Reference:   transaction.AccountID,
					ConnectorID: connectorID,
				}
			default:
				batchElement.Payment.DestinationAccountID = &models.AccountID{
					Reference:   transaction.AccountID,
					ConnectorID: connectorID,
				}
			}

			batch = append(batch, batchElement)

			newState.LastUpdatedAt = transaction.UpdatedAt
		}

		err = ingester.IngestPayments(ctx, batch)
		if err != nil {
			return fetchTransactionsState{}, err
		}
	}

	return newState, nil
}

func matchTransactionType(transactionType string) models.PaymentType {
	switch transactionType {
	case "credit":
		return models.PaymentTypePayOut
	case "debit":
		return models.PaymentTypePayIn
	}

	return models.PaymentTypeOther
}

func matchTransactionStatus(transactionStatus string) models.PaymentStatus {
	switch transactionStatus {
	case "completed":
		return models.PaymentStatusSucceeded
	case "pending":
		return models.PaymentStatusPending
	case "deleted":
		return models.PaymentStatusFailed
	}

	return models.PaymentStatusOther
}
