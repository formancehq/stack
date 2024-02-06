package moneycorp

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchTransactionsState struct {
	LastCreatedAt time.Time `json:"last_created_at"`
}

func taskFetchTransactions(client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"moneycorp.taskFetchTransactions",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("accountID", accountID),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchTransactionsState{})

		newState, err := fetchTransactions(ctx, client, accountID, connectorID, ingester, state)
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
	accountID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	state fetchTransactionsState,
) (fetchTransactionsState, error) {
	newState := fetchTransactionsState{
		LastCreatedAt: state.LastCreatedAt,
	}

	for page := 0; ; page++ {
		pagedTransactions, err := client.GetTransactions(ctx, accountID, page, pageSize, state.LastCreatedAt)
		if err != nil {
			return fetchTransactionsState{}, err
		}

		if len(pagedTransactions) == 0 {
			break
		}

		batch := ingestion.PaymentBatch{}
		for _, transaction := range pagedTransactions {
			createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", transaction.Attributes.CreatedAt)
			if err != nil {
				return fetchTransactionsState{}, fmt.Errorf("failed to parse transaction date: %w", err)
			}

			switch createdAt.Compare(state.LastCreatedAt) {
			case -1, 0:
				continue
			default:
			}

			newState.LastCreatedAt = createdAt

			batchElement, err := toPaymentBatch(connectorID, transaction)
			if err != nil {
				return fetchTransactionsState{}, err
			}

			if batchElement == nil {
				continue
			}

			batch = append(batch, *batchElement)
		}

		if err := ingester.IngestPayments(ctx, batch); err != nil {
			return fetchTransactionsState{}, err
		}

		if len(pagedTransactions) < pageSize {
			break
		}
	}

	return fetchTransactionsState{}, nil
}

func toPaymentBatch(
	connectorID models.ConnectorID,
	transaction *client.Transaction,
) (*ingestion.PaymentBatchElement, error) {
	rawData, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	paymentType, shouldBeRecorded := matchPaymentType(transaction.Attributes.Type, transaction.Attributes.Direction)
	if !shouldBeRecorded {
		return nil, nil
	}

	createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", transaction.Attributes.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction date: %w", err)
	}

	var amount big.Float
	_, ok := amount.SetString(transaction.Attributes.Amount.String())
	if !ok {
		return nil, fmt.Errorf("failed to parse amount %s", transaction.Attributes.Amount.String())
	}

	c, err := currency.GetPrecision(supportedCurrenciesWithDecimal, transaction.Attributes.Currency)
	if err != nil {
		return nil, err
	}

	var amountInt big.Int
	amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(c)))).Int(&amountInt)

	batchElement := ingestion.PaymentBatchElement{
		Payment: &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: transaction.ID,
					Type:      paymentType,
				},
				ConnectorID: connectorID,
			},
			CreatedAt:     createdAt,
			Reference:     transaction.ID,
			ConnectorID:   connectorID,
			Amount:        &amountInt,
			InitialAmount: &amountInt,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transaction.Attributes.Currency),
			Type:          paymentType,
			Status:        models.PaymentStatusSucceeded,
			Scheme:        models.PaymentSchemeOther,
			RawData:       rawData,
		},
	}

	switch paymentType {
	case models.PaymentTypePayIn:
		batchElement.Payment.DestinationAccountID = &models.AccountID{
			Reference:   strconv.Itoa(int(transaction.Attributes.AccountID)),
			ConnectorID: connectorID,
		}
	case models.PaymentTypePayOut:
		batchElement.Payment.SourceAccountID = &models.AccountID{
			Reference:   strconv.Itoa(int(transaction.Attributes.AccountID)),
			ConnectorID: connectorID,
		}
	default:
		if transaction.Attributes.Direction == "Debit" {
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference:   strconv.Itoa(int(transaction.Attributes.AccountID)),
				ConnectorID: connectorID,
			}
		} else {
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference:   strconv.Itoa(int(transaction.Attributes.AccountID)),
				ConnectorID: connectorID,
			}
		}
	}

	return &batchElement, nil
}

func matchPaymentType(transactionType string, transactionDirection string) (models.PaymentType, bool) {
	switch transactionType {
	case "Transfer":
		return models.PaymentTypeTransfer, true
	case "Payment", "Exchange", "Charge", "Refund":
		switch transactionDirection {
		case "Debit":
			return models.PaymentTypePayOut, true
		case "Credit":
			return models.PaymentTypePayIn, true
		}
	}

	return models.PaymentTypeOther, false
}
