package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/genericclient"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
)

type paymentsState struct {
	LastUpdatedAtFrom time.Time `json:"last_updated_at_from"`
}

func (p Plugin) fetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	var oldState paymentsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}
	}

	newState := paymentsState{
		LastUpdatedAtFrom: oldState.LastUpdatedAtFrom,
	}

	payments := make([]models.PSPPayment, 0, req.PageSize)
	hasMore := false
	for page := 0; ; page++ {
		pagedPayments, err := p.client.ListTransactions(ctx, int64(page), int64(req.PageSize), oldState.LastUpdatedAtFrom)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}

		if len(pagedPayments) == 0 {
			break
		}

		for _, payment := range pagedPayments {
			switch payment.UpdatedAt.Compare(oldState.LastUpdatedAtFrom) {
			case -1, 0:
				// Payment already ingested, skip
				continue
			default:
			}

			raw, err := json.Marshal(payment)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, err
			}

			paymentType := matchPaymentType(payment.Type)
			paymentStatus := matchPaymentStatus(payment.Status)

			var amount big.Int
			_, ok := amount.SetString(payment.Amount, 10)
			if !ok {
				return models.FetchNextPaymentsResponse{}, fmt.Errorf("failed to parse amount: %s", payment.Amount)
			}

			p := models.PSPPayment{
				Reference: payment.Id,
				CreatedAt: payment.CreatedAt,
				Type:      paymentType,
				Amount:    &amount,
				Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, payment.Currency),
				Scheme:    models.PAYMENT_SCHEME_OTHER,
				Status:    paymentStatus,
				Metadata:  payment.Metadata,
				Raw:       raw,
			}

			if payment.RelatedTransactionID != nil {
				p.Reference = *payment.RelatedTransactionID
			}

			if payment.SourceAccountID != nil {
				p.SourceAccountReference = payment.SourceAccountID
			}

			if payment.DestinationAccountID != nil {
				p.DestinationAccountReference = payment.DestinationAccountID
			}

			payments = append(payments, p)

			newState.LastUpdatedAtFrom = payment.UpdatedAt

			if len(payments) >= req.PageSize {
				break
			}
		}

		if len(pagedPayments) < req.PageSize {
			break
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

func matchPaymentType(
	transactionType genericclient.TransactionType,
) models.PaymentType {
	switch transactionType {
	case genericclient.PAYIN:
		return models.PAYMENT_TYPE_PAYIN
	case genericclient.PAYOUT:
		return models.PAYMENT_TYPE_PAYOUT
	case genericclient.TRANSFER:
		return models.PAYMENT_TYPE_TRANSFER
	default:
		return models.PAYMENT_TYPE_OTHER
	}
}

func matchPaymentStatus(
	status genericclient.TransactionStatus,
) models.PaymentStatus {
	switch status {
	case genericclient.PENDING:
		return models.PAYMENT_STATUS_PENDING
	case genericclient.FAILED:
		return models.PAYMENT_STATUS_FAILED
	case genericclient.SUCCEEDED:
		return models.PAYMENT_STATUS_SUCCEEDED
	default:
		return models.PAYMENT_STATUS_OTHER
	}
}
