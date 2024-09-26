package currencycloud

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/currencycloud/client"
	"github.com/formancehq/payments/internal/models"
)

type paymentsState struct {
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
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
		LastUpdatedAt: oldState.LastUpdatedAt,
	}

	var payments []models.PSPPayment
	hasMore := false
	page := 1
	for {
		pagedTransactions, nextPage, err := p.client.GetTransactions(ctx, page, req.PageSize, newState.LastUpdatedAt)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}

		if len(pagedTransactions) == 0 {
			break
		}

		for _, transaction := range pagedTransactions {
			switch transaction.UpdatedAt.Compare(newState.LastUpdatedAt) {
			case -1, 0:
				continue
			default:
			}

			payment, err := transactionToPayment(transaction)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, err
			}

			if payment != nil {
				payments = append(payments, *payment)
			}

			newState.LastUpdatedAt = transaction.UpdatedAt

			if len(payments) >= req.PageSize {
				break
			}
		}

		if len(payments) >= req.PageSize {
			hasMore = true
			break
		}

		if nextPage == -1 {
			break
		}

		page = nextPage
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

func transactionToPayment(transaction client.Transaction) (*models.PSPPayment, error) {
	raw, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	precision, ok := supportedCurrenciesWithDecimal[transaction.Currency]
	if !ok {
		return nil, nil
	}

	amount, err := currency.GetAmountWithPrecisionFromString(transaction.Amount.String(), precision)
	if err != nil {
		return nil, err
	}

	paymentType := matchTransactionType(transaction.Type)

	payment := &models.PSPPayment{
		Reference: transaction.ID,
		CreatedAt: transaction.CreatedAt,
		Type:      paymentType,
		Amount:    amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, transaction.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    matchTransactionStatus(transaction.Status),
		Raw:       raw,
	}

	switch paymentType {
	case models.PAYMENT_TYPE_PAYOUT:
		payment.SourceAccountReference = &transaction.AccountID
	default:
		payment.DestinationAccountReference = &transaction.AccountID
	}

	return payment, nil
}

func matchTransactionType(transactionType string) models.PaymentType {
	switch transactionType {
	case "credit":
		return models.PAYMENT_TYPE_PAYIN
	case "debit":
		return models.PAYMENT_TYPE_PAYOUT
	}
	return models.PAYMENT_TYPE_OTHER
}

func matchTransactionStatus(transactionStatus string) models.PaymentStatus {
	switch transactionStatus {
	case "completed":
		return models.PAYMENT_STATUS_SUCCEEDED
	case "pending":
		return models.PAYMENT_STATUS_PENDING
	case "deleted":
		return models.PAYMENT_STATUS_FAILED
	}
	return models.PAYMENT_STATUS_OTHER
}
