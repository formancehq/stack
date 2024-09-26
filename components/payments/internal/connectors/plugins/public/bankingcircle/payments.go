package bankingcircle

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/bankingcircle/client"
	"github.com/formancehq/payments/internal/models"
)

type paymentsState struct {
	LatestStatusChangedTimestamp time.Time `json:"latestStatusChangedTimestamp"`
}

func (p Plugin) fetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	var oldState paymentsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}
	}

	newState := paymentsState{
		LatestStatusChangedTimestamp: oldState.LatestStatusChangedTimestamp,
	}

	var payments []models.PSPPayment
	hasMore := false
	for page := 1; ; page++ {
		pagedPayments, err := p.client.GetPayments(ctx, page, req.PageSize)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}

		if len(pagedPayments) == 0 {
			break
		}

		for _, payment := range pagedPayments {
			switch payment.LatestStatusChangedTimestamp.Compare(oldState.LatestStatusChangedTimestamp) {
			case -1, 0:
				continue
			default:
			}

			p, err := translatePayment(payment)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, err
			}

			if p != nil {
				payments = append(payments, *p)
			}

			if payment.LatestStatusChangedTimestamp.After(newState.LatestStatusChangedTimestamp) {
				newState.LatestStatusChangedTimestamp = payment.LatestStatusChangedTimestamp
			}

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
	}, err
}

func translatePayment(from client.Payment) (*models.PSPPayment, error) {
	raw, err := json.Marshal(from)
	if err != nil {
		return nil, err
	}

	paymentType := matchPaymentType(from.Classification)

	precision, ok := supportedCurrenciesWithDecimal[from.Transfer.Amount.Currency]
	if !ok {
		return nil, nil
	}

	amount, err := currency.GetAmountWithPrecisionFromString(from.Transfer.Amount.Amount.String(), precision)
	if err != nil {
		return nil, err
	}

	payment := models.PSPPayment{
		Reference: from.PaymentID,
		CreatedAt: from.ProcessedTimestamp,
		Type:      paymentType,
		Amount:    amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, from.Transfer.Amount.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    matchPaymentStatus(from.Status),
		Raw:       raw,
	}

	if from.DebtorInformation.AccountID != "" {
		payment.SourceAccountReference = &from.DebtorInformation.AccountID
	}

	if from.CreditorInformation.AccountID != "" {
		payment.DestinationAccountReference = &from.CreditorInformation.AccountID
	}

	return &payment, nil
}

func matchPaymentStatus(paymentStatus string) models.PaymentStatus {
	switch paymentStatus {
	case "Processed":
		return models.PAYMENT_STATUS_SUCCEEDED
	// On MissingFunding - the payment is still in progress.
	// If there will be funds available within 10 days - the payment will be processed.
	// Otherwise - it will be cancelled.
	case "PendingProcessing", "MissingFunding":
		return models.PAYMENT_STATUS_PENDING
	case "Rejected", "Cancelled", "Reversed", "Returned":
		return models.PAYMENT_STATUS_FAILED
	}

	return models.PAYMENT_STATUS_OTHER
}

func matchPaymentType(paymentType string) models.PaymentType {
	switch paymentType {
	case "Incoming":
		return models.PAYMENT_TYPE_PAYIN
	case "Outgoing":
		return models.PAYMENT_TYPE_PAYOUT
	case "Own":
		return models.PAYMENT_TYPE_TRANSFER
	}

	return models.PAYMENT_TYPE_OTHER
}
