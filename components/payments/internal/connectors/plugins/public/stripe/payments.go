package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	stripesdk "github.com/stripe/stripe-go/v79"
)

var (
	ErrInvalidPaymentSource       = errors.New("payment source is invalid")
	ErrUnsupportedAdjustment      = errors.New("unsupported adjustment")
	ErrUnsupportedTransactionType = errors.New("unsupported TransactionType")
	ErrUnsupportedCurrency        = errors.New("unsupported currency")
)

type PaymentState struct {
	LastID string `json:"lastID,omitempty"`
}

func (p *Plugin) fetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	var oldState PaymentState
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

	var payments []models.PSPPayment
	newState := PaymentState{}
	rawPayments, hasMore, err := p.client.GetPayments(ctx, &from.Reference, &oldState.LastID, PageLimit)
	if err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}

	for _, rawPayment := range rawPayments {
		newState.LastID = rawPayment.ID
		payment, err := p.translatePayment(&from.Reference, rawPayment)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, fmt.Errorf("failed to translate payment: %w", err)
		}

		// skip unsupported payments
		if payment == nil {
			continue
		}
		payments = append(payments, *payment)
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

func (p *Plugin) translatePayment(accountRef *string, balanceTransaction *stripesdk.BalanceTransaction) (payment *models.PSPPayment, err error) {
	if balanceTransaction.Source == nil {
		return nil, fmt.Errorf("payment source is invalid")
	}

	rawData, err := json.Marshal(balanceTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw data: %w", err)
	}
	metadata := make(map[string]string)

	switch balanceTransaction.Type {
	case stripesdk.BalanceTransactionTypeCharge:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		appendMetadata(metadata, balanceTransaction.Source.Charge.Metadata)
		if balanceTransaction.Source.Charge.PaymentIntent != nil {
			appendMetadata(metadata, balanceTransaction.Source.Charge.PaymentIntent.Metadata)
		}

		payment = &models.PSPPayment{
			Reference:                   balanceTransaction.ID,
			Type:                        models.PAYMENT_TYPE_PAYIN,
			Status:                      models.PAYMENT_STATUS_SUCCEEDED,
			Amount:                      big.NewInt(balanceTransaction.Source.Charge.Amount - balanceTransaction.Source.Charge.AmountRefunded),
			Asset:                       currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                      detailsToPaymentScheme(balanceTransaction.Source.Charge.PaymentMethodDetails),
			CreatedAt:                   time.Unix(balanceTransaction.Created, 0),
			DestinationAccountReference: accountRef,
			Metadata:                    metadata,
			Raw:                         rawData,
		}

	case stripesdk.BalanceTransactionTypeRefund:
		// Refund a charge
		// Created when a credit card charge refund is initiated.
		// If you authorize and capture separately and the capture amount is
		// less than the initial authorization, you see a balance transaction
		// of type charge for the full authorization amount and another balance
		// transaction of type refund for the uncaptured portion.
		// cf https://stripe.com/docs/reports/balance-transaction-types

		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		if balanceTransaction.Source.Refund.Charge == nil {
			return nil, fmt.Errorf("refund charge missing from refund payload: %q", balanceTransaction.ID)
		}

		appendMetadata(metadata, balanceTransaction.Source.Refund.Metadata)
		if balanceTransaction.Source.Refund.PaymentIntent != nil {
			appendMetadata(metadata, balanceTransaction.Source.Refund.PaymentIntent.Metadata)
		}

		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference:                   balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			Type:                        models.PAYMENT_TYPE_PAYIN,
			Status:                      models.PAYMENT_STATUS_REFUNDED,
			Amount:                      big.NewInt(balanceTransaction.Source.Refund.Amount),
			Asset:                       currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                      detailsToPaymentScheme(balanceTransaction.Source.Refund.Charge.PaymentMethodDetails),
			CreatedAt:                   time.Unix(balanceTransaction.Created, 0),
			DestinationAccountReference: accountRef,
			Raw:                         rawData,
			Metadata:                    metadata,
		}

	case stripesdk.BalanceTransactionTypeRefundFailure:
		// Refund a charge
		// Created when a credit card charge refund fails, and Stripe returns the funds to your balance.
		// This may occur if your customer’s bank or card issuer is unable to correctly process a refund
		// (e.g., due to a closed bank account or a problem with the card).
		// cf https://stripe.com/docs/reports/balance-transaction-types
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		if balanceTransaction.Source.Refund.Charge == nil {
			return nil, fmt.Errorf("refund charge missing from refund payload: %q", balanceTransaction.ID)
		}

		appendMetadata(metadata, balanceTransaction.Source.Refund.Metadata)
		if balanceTransaction.Source.Refund.PaymentIntent != nil {
			appendMetadata(metadata, balanceTransaction.Source.Refund.PaymentIntent.Metadata)
		}

		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference:                   balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			Type:                        models.PAYMENT_TYPE_PAYIN,
			Status:                      models.PAYMENT_STATUS_REFUNDED_FAILURE,
			Amount:                      big.NewInt(balanceTransaction.Source.Refund.Amount),
			Asset:                       currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                      detailsToPaymentScheme(balanceTransaction.Source.Refund.Charge.PaymentMethodDetails),
			CreatedAt:                   time.Unix(balanceTransaction.Created, 0),
			DestinationAccountReference: accountRef,
			Raw:                         rawData,
			Metadata:                    metadata,
		}

	case stripesdk.BalanceTransactionTypePayment:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		appendMetadata(metadata, balanceTransaction.Source.Charge.Metadata)
		if balanceTransaction.Source.Charge.PaymentIntent != nil {
			appendMetadata(metadata, balanceTransaction.Source.Charge.PaymentIntent.Metadata)
		}

		payment = &models.PSPPayment{
			Reference:                   balanceTransaction.ID,
			Type:                        models.PAYMENT_TYPE_PAYIN,
			Status:                      models.PAYMENT_STATUS_SUCCEEDED,
			Amount:                      big.NewInt(balanceTransaction.Source.Charge.Amount - balanceTransaction.Source.Charge.AmountRefunded),
			Asset:                       currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                      models.PAYMENT_SCHEME_OTHER,
			CreatedAt:                   time.Unix(balanceTransaction.Created, 0),
			DestinationAccountReference: accountRef,
			Raw:                         rawData,
			Metadata:                    metadata,
		}

	case stripesdk.BalanceTransactionTypePaymentRefund:
		// Refund a payment
		// Created when a local payment method refund is initiated.
		// Additionally, if your customer’s bank or card issuer is unable to correctly process a refund
		// (e.g., due to a closed bank account or a problem with the card) Stripe returns the funds to your balance.
		// The returned funds are represented as a Balance transaction with the type payment_refund.
		// cf https://stripe.com/docs/reports/balance-transaction-types
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		if balanceTransaction.Source.Refund.Charge == nil {
			return nil, fmt.Errorf("refund charge missing from payment refund payload: %q", balanceTransaction.ID)
		}

		appendMetadata(metadata, balanceTransaction.Source.Refund.Charge.Metadata)
		if balanceTransaction.Source.Refund.Charge.PaymentIntent != nil {
			appendMetadata(balanceTransaction.Source.Refund.Charge.PaymentIntent.Metadata)
		}

		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference:                   balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			Type:                        models.PAYMENT_TYPE_PAYIN,
			Status:                      models.PAYMENT_STATUS_REFUNDED,
			Amount:                      big.NewInt(balanceTransaction.Source.Refund.Amount),
			Asset:                       currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                      models.PAYMENT_SCHEME_OTHER,
			CreatedAt:                   time.Unix(balanceTransaction.Created, 0),
			DestinationAccountReference: accountRef,
			Raw:                         rawData,
			Metadata:                    metadata,
		}

	case stripesdk.BalanceTransactionTypePaymentFailureRefund:
		// Refund a payment
		// ACH, direct debit, and other delayed notification payment methods remain in a pending state
		// until they either succeed or fail. You’ll see a pending Balance transaction of type payment
		// when the payment is created. Another Balance transaction of type payment_failure_refund appears
		// if the pending payment later fails.
		// cf https://stripe.com/docs/reports/balance-transaction-types
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		if balanceTransaction.Source.Refund.Charge == nil {
			return nil, fmt.Errorf("refund charge missing from payment refund failure payload: %q", balanceTransaction.ID)
		}

		appendMetadata(metadata, balanceTransaction.Source.Refund.Charge.Metadata)
		if balanceTransaction.Source.Refund.Charge.PaymentIntent != nil {
			appendMetadata(metadata, balanceTransaction.Source.Refund.Charge.PaymentIntent.Metadata)
		}

		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference:                   balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			Type:                        models.PAYMENT_TYPE_PAYIN,
			Status:                      models.PAYMENT_STATUS_REFUNDED_FAILURE,
			Amount:                      big.NewInt(balanceTransaction.Source.Refund.Amount),
			Asset:                       currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                      models.PAYMENT_SCHEME_OTHER,
			CreatedAt:                   time.Unix(balanceTransaction.Created, 0),
			DestinationAccountReference: accountRef,
			Raw:                         rawData,
			Metadata:                    metadata,
		}

	case stripesdk.BalanceTransactionTypePayout:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Payout.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		appendMetadata(metadata, balanceTransaction.Source.Payout.Metadata)

		payment = &models.PSPPayment{
			Reference: balanceTransaction.ID,
			Type:      models.PAYMENT_TYPE_PAYOUT,
			Status:    convertPayoutStatus(balanceTransaction.Source.Payout.Status),
			Amount:    big.NewInt(balanceTransaction.Source.Payout.Amount),
			Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme: func() models.PaymentScheme {
				switch balanceTransaction.Source.Payout.Type {
				case stripesdk.PayoutTypeBank:
					return models.PAYMENT_SCHEME_SEPA_CREDIT
				case stripesdk.PayoutTypeCard:
					return destinationToPaymentScheme(balanceTransaction.Source.Payout.Destination)
				}
				return models.PAYMENT_SCHEME_UNKNOWN
			}(),
			CreatedAt:              time.Unix(balanceTransaction.Created, 0),
			SourceAccountReference: accountRef,
			Raw:                    rawData,
			Metadata:               metadata,
		}

	case stripesdk.BalanceTransactionTypePayoutFailure, stripesdk.BalanceTransactionTypePayoutCancel:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Payout.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		status := models.PAYMENT_STATUS_FAILED
		if balanceTransaction.Type == stripesdk.BalanceTransactionTypePayoutCancel {
			status = models.PAYMENT_STATUS_CANCELLED
		}

		appendMetadata(metadata, balanceTransaction.Source.Payout.Metadata)
		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference: balanceTransaction.Source.Payout.BalanceTransaction.ID,
			Type:      models.PAYMENT_TYPE_PAYOUT,
			Status:    status,
			Amount:    big.NewInt(balanceTransaction.Source.Payout.Amount),
			Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme: func() models.PaymentScheme {
				switch balanceTransaction.Source.Payout.Type {
				case stripesdk.PayoutTypeBank:
					return models.PAYMENT_SCHEME_SEPA_CREDIT
				case stripesdk.PayoutTypeCard:
					return destinationToPaymentScheme(balanceTransaction.Source.Payout.Destination)
				}
				return models.PAYMENT_SCHEME_UNKNOWN
			}(),
			CreatedAt:              time.Unix(balanceTransaction.Created, 0),
			SourceAccountReference: accountRef,
			Raw:                    rawData,
			Metadata:               metadata,
		}

	case stripesdk.BalanceTransactionTypeTransfer:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Transfer.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		appendMetadata(metadata, balanceTransaction.Source.Transfer.Metadata)
		payment = &models.PSPPayment{
			Reference:              balanceTransaction.ID,
			Type:                   models.PAYMENT_TYPE_TRANSFER,
			Status:                 models.PAYMENT_STATUS_SUCCEEDED,
			Amount:                 big.NewInt(balanceTransaction.Source.Transfer.Amount - balanceTransaction.Source.Transfer.AmountReversed),
			Asset:                  currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                 models.PAYMENT_SCHEME_OTHER,
			CreatedAt:              time.Unix(balanceTransaction.Created, 0),
			SourceAccountReference: accountRef,
			Raw:                    rawData,
			Metadata:               metadata,
		}

		if balanceTransaction.Source.Transfer.Destination != nil {
			payment.DestinationAccountReference = &balanceTransaction.Source.Transfer.Destination.ID
		}

	case stripesdk.BalanceTransactionTypeTransferRefund:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Transfer.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		appendMetadata(metadata, balanceTransaction.Source.Transfer.Metadata)
		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference:              balanceTransaction.Source.Transfer.BalanceTransaction.ID,
			Type:                   models.PAYMENT_TYPE_TRANSFER,
			Status:                 models.PAYMENT_STATUS_REFUNDED,
			Amount:                 big.NewInt(balanceTransaction.Amount),
			Asset:                  currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                 models.PAYMENT_SCHEME_OTHER,
			SourceAccountReference: accountRef,
			CreatedAt:              time.Unix(balanceTransaction.Created, 0),
			Raw:                    rawData,
			Metadata:               metadata,
		}

		if balanceTransaction.Source.Transfer.Destination != nil {
			payment.DestinationAccountReference = &balanceTransaction.Source.Transfer.Destination.ID
		}

	case stripesdk.BalanceTransactionTypeTransferCancel, stripesdk.BalanceTransactionTypeTransferFailure:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Transfer.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		status := models.PAYMENT_STATUS_FAILED
		if balanceTransaction.Type == stripesdk.BalanceTransactionTypeTransferCancel {
			status = models.PAYMENT_STATUS_CANCELLED
		}

		appendMetadata(metadata, balanceTransaction.Source.Transfer.Metadata)
		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference:              balanceTransaction.Source.Transfer.BalanceTransaction.ID,
			Type:                   models.PAYMENT_TYPE_TRANSFER,
			Status:                 status,
			Amount:                 big.NewInt(balanceTransaction.Amount),
			Asset:                  currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                 models.PAYMENT_SCHEME_OTHER,
			SourceAccountReference: accountRef,
			CreatedAt:              time.Unix(balanceTransaction.Created, 0),
			Raw:                    rawData,
			Metadata:               metadata,
		}

		if balanceTransaction.Source.Transfer.Destination != nil {
			payment.DestinationAccountReference = &balanceTransaction.Source.Transfer.Destination.ID
		}

	case stripesdk.BalanceTransactionTypeAdjustment:
		if balanceTransaction.Source.Dispute == nil {
			// We are only handling dispute adjustments
			return nil, ErrUnsupportedAdjustment
		}

		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Dispute.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return nil, fmt.Errorf("%w %q", ErrUnsupportedCurrency, transactionCurrency)
		}

		disputeStatus := convertDisputeStatus(balanceTransaction.Source.Dispute.Status)
		paymentStatus := models.PAYMENT_STATUS_DISPUTE
		switch disputeStatus {
		case models.PAYMENT_STATUS_DISPUTE_WON:
			paymentStatus = models.PAYMENT_STATUS_SUCCEEDED
		case models.PAYMENT_STATUS_DISPUTE_LOST:
			paymentStatus = models.PAYMENT_STATUS_FAILED
		}

		appendMetadata(metadata, balanceTransaction.Source.Dispute.Charge.Metadata)
		if balanceTransaction.Source.Dispute.Charge.PaymentIntent != nil {
			appendMetadata(metadata, balanceTransaction.Source.Dispute.Charge.PaymentIntent.Metadata)
		}

		payment = &models.PSPPayment{
			// ID of original transaction to ensure the refund is appended to the original record
			Reference:                   balanceTransaction.Source.Dispute.Charge.BalanceTransaction.ID,
			Type:                        models.PAYMENT_TYPE_PAYIN,
			Status:                      paymentStatus, // Dispute is occuring, we don't know the outcome yet
			Amount:                      big.NewInt(balanceTransaction.Amount),
			Asset:                       currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:                      detailsToPaymentScheme(balanceTransaction.Source.Dispute.Charge.PaymentMethodDetails),
			CreatedAt:                   time.Unix(balanceTransaction.Created, 0),
			DestinationAccountReference: accountRef,
			Raw:                         rawData,
			Metadata:                    metadata,
		}

	default:
		// TODO (laouji): standardize logging
		log.Printf("stripe fetch payments task %s %q", ErrUnsupportedTransactionType, balanceTransaction.Type)
		return nil, nil
	}

	return payment, err
}

func convertDisputeStatus(status stripesdk.DisputeStatus) models.PaymentStatus {
	switch status {
	case stripesdk.DisputeStatusNeedsResponse, stripesdk.DisputeStatusUnderReview:
		return models.PAYMENT_STATUS_DISPUTE
	case stripesdk.DisputeStatusLost:
		return models.PAYMENT_STATUS_DISPUTE_LOST
	case stripesdk.DisputeStatusWon:
		return models.PAYMENT_STATUS_DISPUTE_WON
	default:
		return models.PAYMENT_STATUS_DISPUTE
	}
}

func convertPayoutStatus(status stripesdk.PayoutStatus) models.PaymentStatus {
	switch status {
	case stripesdk.PayoutStatusCanceled:
		return models.PAYMENT_STATUS_CANCELLED
	case stripesdk.PayoutStatusFailed:
		return models.PAYMENT_STATUS_FAILED
	case stripesdk.PayoutStatusInTransit, stripesdk.PayoutStatusPending:
		return models.PAYMENT_STATUS_PENDING
	case stripesdk.PayoutStatusPaid:
		return models.PAYMENT_STATUS_SUCCEEDED
	}

	return models.PAYMENT_STATUS_OTHER
}

func appendMetadata(s map[string]string, vs ...map[string]string) {
	for _, in := range vs {
		if in == nil {
			continue
		}
		for k, v := range in {
			s[k] = v
		}
	}
}

func detailsToPaymentScheme(details *stripesdk.ChargePaymentMethodDetails) models.PaymentScheme {
	if details == nil || details.Card == nil {
		return models.PAYMENT_SCHEME_UNKNOWN
	}
	scheme, _ := models.PaymentSchemeFromString(
		fmt.Sprintf("CARD_%s", strings.ToUpper(string(details.Card.Brand))),
	)
	return scheme
}

func destinationToPaymentScheme(dest *stripesdk.PayoutDestination) models.PaymentScheme {
	if dest == nil || dest.Card == nil {
		return models.PAYMENT_SCHEME_UNKNOWN
	}

	switch dest.Card.Brand {
	case stripesdk.CardBrandAmericanExpress:
		return models.PAYMENT_SCHEME_CARD_AMEX
	case stripesdk.CardBrandDiscover:
		return models.PAYMENT_SCHEME_CARD_DISCOVER
	case stripesdk.CardBrandDinersClub:
		return models.PAYMENT_SCHEME_CARD_DINERS
	case stripesdk.CardBrandJCB:
		return models.PAYMENT_SCHEME_CARD_JCB
	case stripesdk.CardBrandMasterCard:
		return models.PAYMENT_SCHEME_CARD_MASTERCARD
	case stripesdk.CardBrandUnionPay:
		return models.PAYMENT_SCHEME_CARD_UNION_PAY
	case stripesdk.CardBrandVisa:
		return models.PAYMENT_SCHEME_CARD_VISA
	default:
		return models.PAYMENT_SCHEME_UNKNOWN
	}
}
