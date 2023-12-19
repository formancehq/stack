package modulr

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchTransactions(logger logging.Logger, client *client.Client, accountID string) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		logger.Info("Fetching transactions for account", accountID)

		transactions, err := client.GetTransactions(ctx, accountID)
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

			paymentType := matchTransactionType(transaction.Type)

			precision, ok := supportedCurrenciesWithDecimal[transaction.Account.Currency]
			if !ok {
				logger.Errorf("currency %s is not supported", transaction.Account.Currency)
				continue
			}

			var amount big.Float
			_, ok = amount.SetString(transaction.Amount.String())
			if !ok {
				return fmt.Errorf("failed to parse amount %s", transaction.Amount.String())
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

		if err := ingester.IngestPayments(ctx, connectorID, batch, struct{}{}); err != nil {
			return err
		}

		return nil
	}
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
