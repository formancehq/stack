package modulr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/modulr/client"
	"github.com/formancehq/payments/internal/models"
)

type paymentsState struct {
	LastTransactionTime time.Time `json:"lastTransactionTime"`
}

func (p Plugin) fetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	var oldState paymentsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}
	}

	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextPaymentsResponse{}, errors.New("missing from payload when fetching payments")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}

	newState := paymentsState{
		LastTransactionTime: oldState.LastTransactionTime,
	}

	var payments []models.PSPPayment
	hasMore := false
	for page := 0; ; page++ {
		pagedTransactions, err := p.client.GetTransactions(ctx, from.Reference, page, req.PageSize, oldState.LastTransactionTime)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}

		if len(pagedTransactions.Content) == 0 {
			break
		}

		for _, transaction := range pagedTransactions.Content {
			createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", transaction.TransactionDate)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, err
			}

			switch createdTime.Compare(oldState.LastTransactionTime) {
			case -1, 0:
				// Account already ingested, skip
				continue
			default:
			}

			payment, err := transactionToPayment(transaction, from)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, err
			}

			if payment != nil {
				payments = append(payments, *payment)
			}

			newState.LastTransactionTime = createdTime

			if len(payments) >= req.PageSize {
				break
			}
		}

		if len(payments) >= req.PageSize {
			hasMore = true
			break
		}
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}

	return models.FetchNextPaymentsResponse{
		Payments: payments,
		NewState: payload,
		HasMore:  hasMore,
	}, nil
}

func transactionToPayment(transaction client.Transaction, from models.PSPAccount) (*models.PSPPayment, error) {
	raw, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	paymentType := matchTransactionType(transaction.Type)

	precision, ok := supportedCurrenciesWithDecimal[transaction.Account.Currency]
	if !ok {
		return nil, nil
	}

	amount, err := currency.GetAmountWithPrecisionFromString(transaction.Amount.String(), precision)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount %s: %w", transaction.Amount, err)
	}

	createdAt, err := time.Parse("2006-01-02T15:04:05-0700", transaction.PostedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse posted date %s: %w", transaction.PostedDate, err)
	}

	payment := &models.PSPPayment{
		Reference: transaction.ID,
		CreatedAt: createdAt,
		Type:      paymentType,
		Amount:    amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, transaction.Account.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    models.PAYMENT_STATUS_SUCCEEDED,
		Raw:       raw,
	}

	switch paymentType {
	case models.PAYMENT_TYPE_PAYIN:
		payment.DestinationAccountReference = &from.Reference
	case models.PAYMENT_TYPE_PAYOUT:
		payment.SourceAccountReference = &from.Reference
	default:
		if transaction.Credit {
			payment.DestinationAccountReference = &from.Reference
		} else {
			payment.SourceAccountReference = &from.Reference
		}
	}

	return payment, nil
}

func matchTransactionType(transactionType string) models.PaymentType {
	if transactionType == "PI_REV" ||
		transactionType == "PO_REV" ||
		transactionType == "ADHOC" ||
		transactionType == "INT_INTERC" {
		return models.PAYMENT_TYPE_OTHER
	}

	if strings.HasPrefix(transactionType, "PI_") {
		return models.PAYMENT_TYPE_PAYIN
	}

	if strings.HasPrefix(transactionType, "PO_") {
		return models.PAYMENT_TYPE_PAYOUT
	}

	return models.PAYMENT_TYPE_OTHER
}
