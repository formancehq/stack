package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/generic/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/genericclient"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

type fetchTransactionsState struct {
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

func taskFetchTransactions(client *client.Client, config *Config) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		resolver task.StateResolver,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"generic.taskFetchTransactions",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		state := task.MustResolveTo(ctx, resolver, fetchTransactionsState{})

		newState, err := ingestTransactions(ctx, connectorID, client, ingester, scheduler, state)
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
	scheduler task.Scheduler,
	state fetchTransactionsState,
) (fetchTransactionsState, error) {
	newState := fetchTransactionsState{
		LastUpdatedAt: state.LastUpdatedAt,
	}

	for page := 1; ; page++ {
		transactions, err := client.ListTransactions(ctx, int64(page), pageSize, state.LastUpdatedAt)
		if err != nil {
			return fetchTransactionsState{}, err
		}

		if len(transactions) == 0 {
			break
		}

		paymentBatch := make([]ingestion.PaymentBatchElement, 0, len(transactions))
		for _, transaction := range transactions {
			elt, err := translate(ctx, connectorID, transaction)
			if err != nil {
				return fetchTransactionsState{}, err
			}

			paymentBatch = append(paymentBatch, elt)

			newState.LastUpdatedAt = transaction.UpdatedAt
		}

		if err := ingester.IngestPayments(ctx, ingestion.PaymentBatch(paymentBatch)); err != nil {
			return fetchTransactionsState{}, errors.Wrap(task.ErrRetryable, err.Error())
		}
	}

	return newState, nil
}

func translate(
	ctx context.Context,
	connectorID models.ConnectorID,
	transaction genericclient.Transaction,
) (ingestion.PaymentBatchElement, error) {
	paymentType := matchPaymentType(transaction.Type)
	paymentStatus := matchPaymentStatus(transaction.Status)

	var amount big.Int
	_, ok := amount.SetString(transaction.Amount, 10)
	if !ok {
		return ingestion.PaymentBatchElement{}, fmt.Errorf("failed to parse amount: %s", transaction.Amount)
	}

	raw, err := json.Marshal(transaction)
	if err != nil {
		return ingestion.PaymentBatchElement{}, err
	}

	paymentID := models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: transaction.Id,
			Type:      paymentType,
		},
		ConnectorID: connectorID,
	}
	elt := ingestion.PaymentBatchElement{
		Payment: &models.Payment{
			ID:            paymentID,
			ConnectorID:   connectorID,
			CreatedAt:     transaction.CreatedAt,
			Reference:     transaction.Id,
			Amount:        &amount,
			InitialAmount: &amount,
			Type:          paymentType,
			Status:        paymentStatus,
			Scheme:        models.PaymentSchemeOther,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transaction.Currency),
			RawData:       raw,
		},
	}

	if transaction.SourceAccountID != nil && *transaction.SourceAccountID != "" {
		elt.Payment.SourceAccountID = &models.AccountID{
			Reference:   *transaction.SourceAccountID,
			ConnectorID: connectorID,
		}
	}

	if transaction.DestinationAccountID != nil && *transaction.DestinationAccountID != "" {
		elt.Payment.DestinationAccountID = &models.AccountID{
			Reference:   *transaction.DestinationAccountID,
			ConnectorID: connectorID,
		}
	}

	for k, v := range transaction.Metadata {
		elt.Payment.Metadata = append(elt.Payment.Metadata, &models.PaymentMetadata{
			PaymentID: paymentID,
			CreatedAt: transaction.CreatedAt,
			Key:       k,
			Value:     v,
			Changelog: []models.MetadataChangelog{
				{
					CreatedAt: transaction.CreatedAt,
					Value:     v,
				},
			},
		})

	}

	return elt, nil
}

func matchPaymentType(
	transactionType genericclient.TransactionType,
) models.PaymentType {
	switch transactionType {
	case genericclient.PAYIN:
		return models.PaymentTypePayIn
	case genericclient.PAYOUT:
		return models.PaymentTypePayOut
	case genericclient.TRANSFER:
		return models.PaymentTypeTransfer
	default:
		return models.PaymentTypeOther
	}
}

func matchPaymentStatus(
	status genericclient.TransactionStatus,
) models.PaymentStatus {
	switch status {
	case genericclient.PENDING:
		return models.PaymentStatusPending
	case genericclient.FAILED:
		return models.PaymentStatusFailed
	case genericclient.SUCCEEDED:
		return models.PaymentStatusSucceeded
	default:
		return models.PaymentStatusOther
	}
}
