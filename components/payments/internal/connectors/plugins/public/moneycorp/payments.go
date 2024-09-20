package moneycorp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/moneycorp/client"
	"github.com/formancehq/payments/internal/models"
)

type paymentsState struct {
	LastCreatedAt time.Time `json:"lastCreatedAt"`
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
		LastCreatedAt: oldState.LastCreatedAt,
	}

	var payments []models.PSPPayment
	hasMore := false
	for page := 0; ; page++ {
		pagedTransactions, err := p.client.GetTransactions(ctx, from.Reference, page, req.PageSize, oldState.LastCreatedAt)
		if err != nil {
			// retryable error already handled by the client
			return models.FetchNextPaymentsResponse{}, err
		}

		if len(pagedTransactions) == 0 {
			break
		}

		for _, transaction := range pagedTransactions {
			createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", transaction.Attributes.CreatedAt)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, fmt.Errorf("failed to parse transaction date: %v", err)
			}

			switch createdAt.Compare(oldState.LastCreatedAt) {
			case -1, 0:
				continue
			default:
			}

			newState.LastCreatedAt = createdAt

			payment, err := transactionToPayment(transaction)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, err
			}

			if payment != nil {
				payments = append(payments, *payment)
			}

			if len(payments) == req.PageSize {
				break
			}
		}

		if len(pagedTransactions) < req.PageSize {
			break
		}

		if len(payments) == req.PageSize {
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

func transactionToPayment(transaction *client.Transaction) (*models.PSPPayment, error) {
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

	c, err := currency.GetPrecision(supportedCurrenciesWithDecimal, transaction.Attributes.Currency)
	if err != nil {
		return nil, err
	}

	amount, err := currency.GetAmountWithPrecisionFromString(transaction.Attributes.Amount.String(), c)
	if err != nil {
		return nil, err
	}

	payment := models.PSPPayment{
		Reference: transaction.ID,
		CreatedAt: createdAt,
		Type:      paymentType,
		Amount:    amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, transaction.Attributes.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    models.PAYMENT_STATUS_SUCCEEDED,
		Metadata:  map[string]string{},
		Raw:       rawData,
	}

	switch paymentType {
	case models.PAYMENT_TYPE_PAYIN:
		payment.DestinationAccountReference = pointer.For(strconv.Itoa(int(transaction.Attributes.AccountID)))
	case models.PAYMENT_TYPE_PAYOUT:
		payment.SourceAccountReference = pointer.For(strconv.Itoa(int(transaction.Attributes.AccountID)))
	default:
		if transaction.Attributes.Direction == "Debit" {
			payment.SourceAccountReference = pointer.For(strconv.Itoa(int(transaction.Attributes.AccountID)))
		} else {
			payment.DestinationAccountReference = pointer.For(strconv.Itoa(int(transaction.Attributes.AccountID)))
		}
	}

	return &payment, nil
}

func matchPaymentType(transactionType string, transactionDirection string) (models.PaymentType, bool) {
	switch transactionType {
	case "Transfer":
		return models.PAYMENT_TYPE_TRANSFER, true
	case "Payment", "Exchange", "Charge", "Refund":
		switch transactionDirection {
		case "Debit":
			return models.PAYMENT_TYPE_PAYOUT, true
		case "Credit":
			return models.PAYMENT_TYPE_PAYIN, true
		}
	}

	return models.PAYMENT_TYPE_OTHER, false
}
