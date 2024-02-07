package modulr

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fetchTransactionsState struct {
	LastTransactionTime time.Time `json:"last_transaction_time"`
}

func (s *fetchTransactionsState) UpdateLatest(latest *client.Transaction) error {
	transactionTime, err := time.Parse("2006-01-02T15:04:05.999-0700", latest.TransactionDate)
	if err != nil {
		return err
	}
	if transactionTime.After(s.LastTransactionTime) {
		s.LastTransactionTime = transactionTime
	}
	return nil
}

func (s *fetchTransactionsState) FindLatest(transactions []*client.Transaction) error {
	for _, transaction := range transactions {
		if err := s.UpdateLatest(transaction); err != nil {
			return err
		}
	}
	return nil
}

func (s *fetchTransactionsState) IsNew(transaction *client.Transaction) (bool, error) {
	transactionTime, err := time.Parse("2006-01-02T15:04:05.999-0700", transaction.TransactionDate)
	if err != nil {
		return false, err
	}
	return transactionTime.After(s.LastTransactionTime), nil
}

func (s *fetchTransactionsState) FilterNew(transactions []*client.Transaction) ([]*client.Transaction, error) {
	result := make([]*client.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		isNew, err := s.IsNew(transaction)
		if err != nil {
			return nil, err
		}
		if !isNew {
			continue
		}
		result = append(result, transaction)
	}
	return result, nil
}

func (s *fetchTransactionsState) GetFilterValue() string {
	if s.LastTransactionTime.IsZero() {
		return ""
	}
	return s.LastTransactionTime.Format("2006-01-02T15:04:05-0700")
}

func taskFetchTransactions(config Config, client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"modulr.taskFetchTransactions",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("accountID", accountID),
		)
		defer span.End()

		state, err := fetchTransactions(
			ctx,
			config,
			client,
			accountID,
			connectorID,
			ingester,
			task.MustResolveTo(ctx, resolver, fetchTransactionsState{}),
		)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := ingester.UpdateTaskState(ctx, state); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func fetchTransactions(
	ctx context.Context,
	config Config,
	client *client.Client,
	accountID string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	state fetchTransactionsState,
) (fetchTransactionsState, error) {
	newState := state
	for page := 0; ; page++ {
		pagedTransactions, err := client.GetTransactions(
			ctx,
			accountID,
			page,
			config.PageSize,
			state.GetFilterValue(),
		)
		if err != nil {
			return newState, err
		}

		if len(pagedTransactions.Content) == 0 {
			break
		}

		transactions, err := state.FilterNew(pagedTransactions.Content)
		if err != nil {
			return newState, err
		}

		batch, err := toBatch(connectorID, accountID, transactions)
		if err != nil {
			return newState, err
		}

		if err := ingester.IngestPayments(ctx, batch); err != nil {
			return newState, err
		}

		if err := newState.FindLatest(transactions); err != nil {
			return newState, err
		}

		if len(pagedTransactions.Content) < config.PageSize {
			break
		}

		if page+1 >= pagedTransactions.TotalPages {
			// Modulr paging starts at 0, so the last page is TotalPages - 1.
			break
		}
	}

	return newState, nil
}

func toBatch(
	connectorID models.ConnectorID,
	accountID string,
	transactions []*client.Transaction,
) (ingestion.PaymentBatch, error) {
	batch := ingestion.PaymentBatch{}

	for _, transaction := range transactions {

		rawData, err := json.Marshal(transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal transaction: %w", err)
		}

		paymentType := matchTransactionType(transaction.Type)

		precision, ok := supportedCurrenciesWithDecimal[transaction.Account.Currency]
		if !ok {
			continue
		}

		var amount big.Float
		_, ok = amount.SetString(transaction.Amount.String())
		if !ok {
			return nil, fmt.Errorf("failed to parse amount %s", transaction.Amount.String())
		}

		createdAt, err := time.Parse(timeTemplate, transaction.PostedDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse posted date %s: %w", transaction.PostedDate, err)
		}

		var amountInt big.Int
		amount.Mul(&amount, big.NewFloat(math.Pow(10, float64(precision)))).Int(&amountInt)

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
				Type:          paymentType,
				Status:        models.PaymentStatusSucceeded,
				Scheme:        models.PaymentSchemeOther,
				Amount:        &amountInt,
				InitialAmount: &amountInt,
				Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transaction.Account.Currency),
				RawData:       rawData,
			},
		}

		switch paymentType {
		case models.PaymentTypePayIn:
			batchElement.Payment.DestinationAccountID = &models.AccountID{
				Reference:   accountID,
				ConnectorID: connectorID,
			}
		case models.PaymentTypePayOut:
			batchElement.Payment.SourceAccountID = &models.AccountID{
				Reference:   accountID,
				ConnectorID: connectorID,
			}
		default:
			if transaction.Credit {
				batchElement.Payment.DestinationAccountID = &models.AccountID{
					Reference:   accountID,
					ConnectorID: connectorID,
				}
			} else {
				batchElement.Payment.SourceAccountID = &models.AccountID{
					Reference:   accountID,
					ConnectorID: connectorID,
				}
			}
		}

		batch = append(batch, batchElement)
	}

	return batch, nil
}

func matchTransactionType(transactionType string) models.PaymentType {
	if transactionType == "PI_REV" ||
		transactionType == "PO_REV" ||
		transactionType == "ADHOC" ||
		transactionType == "INT_INTERC" {
		return models.PaymentTypeOther
	}

	if strings.HasPrefix(transactionType, "PI_") {
		return models.PaymentTypePayIn
	}

	if strings.HasPrefix(transactionType, "PO_") {
		return models.PaymentTypePayOut
	}

	return models.PaymentTypeOther
}
