package wise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/wise/client"
	"github.com/formancehq/payments/internal/models"
)

type paymentsState struct {
	Offset int
}

func (p Plugin) fetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	var oldState paymentsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}
	}

	var from client.Profile
	if req.FromPayload == nil {
		return models.FetchNextPaymentsResponse{}, errors.New("missing from payload when fetching payments")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}

	newState := paymentsState{
		Offset: oldState.Offset,
	}

	var payments []models.PSPPayment
	hasMore := false
	for {
		pagedTransfers, err := p.client.GetTransfers(ctx, from.ID, newState.Offset, req.PageSize)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}

		if len(pagedTransfers) == 0 {
			break
		}

		for _, transfer := range pagedTransfers {
			payment, err := fromTransferToPayment(transfer)
			if err != nil {
				return models.FetchNextPaymentsResponse{}, err
			}

			payments = append(payments, *payment)
			newState.Offset++

			if len(payments) >= req.PageSize {
				break
			}
		}

		if len(pagedTransfers) < req.PageSize {
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

func fromTransferToPayment(from client.Transfer) (*models.PSPPayment, error) {
	raw, err := json.Marshal(from)
	if err != nil {
		return nil, err
	}

	precision, ok := supportedCurrenciesWithDecimal[from.TargetCurrency]
	if !ok {
		return nil, nil
	}

	amount, err := currency.GetAmountWithPrecisionFromString(from.TargetValue.String(), precision)
	if err != nil {
		return nil, err
	}

	p := models.PSPPayment{
		Reference: fmt.Sprintf("%d", from.ID),
		CreatedAt: from.CreatedAt,
		Type:      models.PAYMENT_TYPE_TRANSFER,
		Amount:    amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, from.TargetCurrency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    matchTransferStatus(from.Status),
		Raw:       raw,
	}

	if from.SourceBalanceID != 0 {
		p.SourceAccountReference = pointer.For(fmt.Sprintf("%d", from.SourceBalanceID))
	}

	if from.DestinationBalanceID != 0 {
		p.DestinationAccountReference = pointer.For(fmt.Sprintf("%d", from.DestinationBalanceID))
	}

	return &p, nil
}

func matchTransferStatus(status string) models.PaymentStatus {
	switch status {
	case "incoming_payment_waiting", "incoming_payment_initiated", "processing", "funds_converted", "bounced_back":
		return models.PAYMENT_STATUS_PENDING
	case "outgoing_payment_sent":
		return models.PAYMENT_STATUS_SUCCEEDED
	case "funds_refunded", "charged_back":
		return models.PAYMENT_STATUS_FAILED
	case "cancelled":
		return models.PAYMENT_STATUS_CANCELLED
	}

	return models.PAYMENT_STATUS_OTHER
}
