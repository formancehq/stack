package stripe

import (
	"encoding/json"
	"log"
	"math/big"
	"runtime/debug"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/internal/models"
	"github.com/stripe/stripe-go/v72"
)

func createBatchElement(
	connectorID models.ConnectorID,
	balanceTransaction *stripe.BalanceTransaction,
	account string,
	forward bool,
) (ingestion.PaymentBatchElement, bool) {
	var payment *models.Payment
	var adjustment *models.PaymentAdjustment

	defer func() {
		// DEBUG
		if e := recover(); e != nil {
			log.Println("Error translating transaction")
			debug.PrintStack()
			spew.Dump(balanceTransaction)
			panic(e)
		}
	}()

	if balanceTransaction.Source == nil {
		return ingestion.PaymentBatchElement{}, false
	}

	rawData, err := json.Marshal(balanceTransaction)
	if err != nil {
		return ingestion.PaymentBatchElement{}, false
	}

	switch balanceTransaction.Type {
	case stripe.BalanceTransactionTypeCharge:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}

		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayIn,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Charge.Amount - balanceTransaction.Source.Charge.AmountRefunded),
			InitialAmount: big.NewInt(balanceTransaction.Source.Charge.Amount),
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			RawData:       rawData,
			Scheme:        models.PaymentScheme(balanceTransaction.Source.Charge.PaymentMethodDetails.Card.Brand),
			CreatedAt:     time.Unix(balanceTransaction.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

	case stripe.BalanceTransactionTypeRefund:
		// Refund a charge
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}
		// Created when a credit card charge refund is initiated.
		// If you authorize and capture separately and the capture amount is
		// less than the initial authorization, you see a balance transaction
		// of type charge for the full authorization amount and another balance
		// transaction of type refund for the uncaptured portion.
		// cf https://stripe.com/docs/reports/balance-transaction-types
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayIn,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Refund.Charge.Amount - balanceTransaction.Source.Refund.Charge.AmountRefunded),
			InitialAmount: big.NewInt(balanceTransaction.Source.Refund.Charge.Amount),
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			RawData:       rawData,
			Scheme:        models.PaymentScheme(balanceTransaction.Source.Refund.Charge.PaymentMethodDetails.Card.Brand),
			CreatedAt:     time.Unix(balanceTransaction.Source.Refund.Charge.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Amount:    big.NewInt(balanceTransaction.Source.Refund.Amount),
			Status:    models.PaymentStatusRefunded,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypeRefundFailure:
		// Refund a charge
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}
		// Created when a credit card charge refund is initiated.
		// If you authorize and capture separately and the capture amount is
		// less than the initial authorization, you see a balance transaction
		// of type charge for the full authorization amount and another balance
		// transaction of type refund for the uncaptured portion.
		// cf https://stripe.com/docs/reports/balance-transaction-types
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayIn,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Refund.Charge.Amount - balanceTransaction.Source.Refund.Charge.AmountRefunded),
			InitialAmount: big.NewInt(balanceTransaction.Source.Refund.Charge.Amount),
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			RawData:       rawData,
			Scheme:        models.PaymentScheme(balanceTransaction.Source.Refund.Charge.PaymentMethodDetails.Card.Brand),
			CreatedAt:     time.Unix(balanceTransaction.Source.Refund.Charge.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Amount:    big.NewInt(balanceTransaction.Source.Refund.Amount),
			Status:    models.PaymentStatusRefundedFailure,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypePayment:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}

		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayIn,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Charge.Amount - balanceTransaction.Source.Charge.AmountRefunded),
			InitialAmount: big.NewInt(balanceTransaction.Source.Charge.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:        models.PaymentSchemeOther,
			CreatedAt:     time.Unix(balanceTransaction.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

	case stripe.BalanceTransactionTypePaymentRefund:
		// Refund a charge
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}
		// Created when a credit card charge refund is initiated.
		// If you authorize and capture separately and the capture amount is
		// less than the initial authorization, you see a balance transaction
		// of type charge for the full authorization amount and another balance
		// transaction of type refund for the uncaptured portion.
		// cf https://stripe.com/docs/reports/balance-transaction-types
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayIn,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Refund.Charge.Amount - balanceTransaction.Source.Refund.Charge.AmountRefunded),
			InitialAmount: big.NewInt(balanceTransaction.Source.Refund.Charge.Amount),
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:        models.PaymentSchemeOther,
			RawData:       rawData,
			CreatedAt:     time.Unix(balanceTransaction.Source.Refund.Charge.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Amount:    big.NewInt(balanceTransaction.Source.Refund.Amount),
			Status:    models.PaymentStatusRefunded,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypePaymentFailureRefund:
		// Refund a charge
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Refund.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}
		// Created when a credit card charge refund is initiated.
		// If you authorize and capture separately and the capture amount is
		// less than the initial authorization, you see a balance transaction
		// of type charge for the full authorization amount and another balance
		// transaction of type refund for the uncaptured portion.
		// cf https://stripe.com/docs/reports/balance-transaction-types
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayIn,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Refund.Charge.Amount - balanceTransaction.Source.Refund.Charge.AmountRefunded),
			InitialAmount: big.NewInt(balanceTransaction.Source.Refund.Charge.Amount),
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			Scheme:        models.PaymentSchemeOther,
			RawData:       rawData,
			CreatedAt:     time.Unix(balanceTransaction.Source.Refund.Charge.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Refund.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Amount:    big.NewInt(balanceTransaction.Source.Refund.Amount),
			Status:    models.PaymentStatusRefundedFailure,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypePayout:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Payout.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}

		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayOut,
			Status:        convertPayoutStatus(balanceTransaction.Source.Payout.Status),
			Amount:        big.NewInt(balanceTransaction.Source.Payout.Amount),
			InitialAmount: big.NewInt(balanceTransaction.Source.Payout.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balanceTransaction.Source.Payout.Currency)),
			Scheme: func() models.PaymentScheme {
				switch balanceTransaction.Source.Payout.Type {
				case stripe.PayoutTypeBank:
					return models.PaymentSchemeSepaCredit
				case stripe.PayoutTypeCard:
					return models.PaymentScheme(balanceTransaction.Source.Payout.Card.Brand)
				}

				return models.PaymentSchemeUnknown
			}(),
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

	case stripe.BalanceTransactionTypePayoutFailure:
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Payout.BalanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Payout.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayOut,
			Status:        models.PaymentStatusFailed,
			Amount:        big.NewInt(balanceTransaction.Source.Payout.Amount),
			InitialAmount: big.NewInt(balanceTransaction.Source.Payout.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balanceTransaction.Source.Payout.Currency)),
			Scheme: func() models.PaymentScheme {
				switch balanceTransaction.Source.Payout.Type {
				case stripe.PayoutTypeBank:
					return models.PaymentSchemeSepaCredit
				case stripe.PayoutTypeCard:
					return models.PaymentScheme(balanceTransaction.Source.Payout.Card.Brand)
				}

				return models.PaymentSchemeUnknown
			}(),
			CreatedAt: time.Unix(balanceTransaction.Source.Payout.Created, 0),
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Payout.BalanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Status:    models.PaymentStatusFailed,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypePayoutCancel:
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Payout.BalanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Payout.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayOut,
			Status:        models.PaymentStatusCancelled,
			Amount:        big.NewInt(balanceTransaction.Source.Payout.Amount),
			InitialAmount: big.NewInt(balanceTransaction.Source.Payout.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balanceTransaction.Source.Payout.Currency)),
			Scheme: func() models.PaymentScheme {
				switch balanceTransaction.Source.Payout.Type {
				case stripe.PayoutTypeBank:
					return models.PaymentSchemeSepaCredit
				case stripe.PayoutTypeCard:
					return models.PaymentScheme(balanceTransaction.Source.Payout.Card.Brand)
				}

				return models.PaymentSchemeUnknown
			}(),
			CreatedAt: time.Unix(balanceTransaction.Source.Payout.Created, 0),
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Payout.BalanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Status:    models.PaymentStatusCancelled,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypeTransfer:
		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Transfer.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}

		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypeTransfer,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypeTransfer,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Transfer.Amount - balanceTransaction.Source.Transfer.AmountReversed),
			InitialAmount: big.NewInt(balanceTransaction.Source.Transfer.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balanceTransaction.Source.Transfer.Currency)),
			Scheme:        models.PaymentSchemeOther,
			CreatedAt:     time.Unix(balanceTransaction.Created, 0),
		}

		if balanceTransaction.Source.Transfer.Destination != nil {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   balanceTransaction.Source.Transfer.Destination.ID,
				ConnectorID: connectorID,
			}
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

	case stripe.BalanceTransactionTypeTransferRefund:
		// Two things to insert here: the balance transaction at the origin
		// of the refund and the balance transaction of the refund, which is an
		// adjustment of the origin.
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Transfer.BalanceTransaction.ID,
					Type:      models.PaymentTypeTransfer,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Transfer.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypeTransfer,
			Status:        models.PaymentStatusSucceeded,
			Amount:        big.NewInt(balanceTransaction.Source.Transfer.Amount - balanceTransaction.Source.Transfer.AmountReversed),
			InitialAmount: big.NewInt(balanceTransaction.Source.Transfer.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balanceTransaction.Source.Transfer.Currency)),
			Scheme:        models.PaymentSchemeOther,
			CreatedAt:     time.Unix(balanceTransaction.Source.Transfer.Created, 0),
		}

		if balanceTransaction.Source.Transfer.Destination != nil {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   balanceTransaction.Source.Transfer.Destination.ID,
				ConnectorID: connectorID,
			}
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Transfer.BalanceTransaction.ID,
					Type:      models.PaymentTypeTransfer,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Amount:    big.NewInt(balanceTransaction.Amount),
			Status:    models.PaymentStatusRefunded,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypeTransferCancel:
		// Two things to insert here: the balance transaction at the origin
		// of the refund and the balance transaction of the refund, which is an
		// adjustment of the origin.
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Transfer.BalanceTransaction.ID,
					Type:      models.PaymentTypeTransfer,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Transfer.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypeTransfer,
			Status:        models.PaymentStatusCancelled,
			Amount:        big.NewInt(balanceTransaction.Source.Transfer.Amount - balanceTransaction.Source.Transfer.AmountReversed),
			InitialAmount: big.NewInt(balanceTransaction.Source.Transfer.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balanceTransaction.Source.Transfer.Currency)),
			Scheme:        models.PaymentSchemeOther,
			CreatedAt:     time.Unix(balanceTransaction.Source.Transfer.Created, 0),
		}

		if balanceTransaction.Source.Transfer.Destination != nil {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   balanceTransaction.Source.Transfer.Destination.ID,
				ConnectorID: connectorID,
			}
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Transfer.BalanceTransaction.ID,
					Type:      models.PaymentTypeTransfer,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Amount:    big.NewInt(balanceTransaction.Amount),
			Status:    models.PaymentStatusCancelled,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypeTransferFailure:
		// Two things to insert here: the balance transaction at the origin
		// of the refund and the balance transaction of the refund, which is an
		// adjustment of the origin.
		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Transfer.BalanceTransaction.ID,
					Type:      models.PaymentTypeTransfer,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Transfer.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypeTransfer,
			Status:        models.PaymentStatusFailed,
			Amount:        big.NewInt(balanceTransaction.Source.Transfer.Amount - balanceTransaction.Source.Transfer.AmountReversed),
			InitialAmount: big.NewInt(balanceTransaction.Source.Transfer.Amount),
			RawData:       rawData,
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, string(balanceTransaction.Source.Transfer.Currency)),
			Scheme:        models.PaymentSchemeOther,
			CreatedAt:     time.Unix(balanceTransaction.Source.Transfer.Created, 0),
		}

		if balanceTransaction.Source.Transfer.Destination != nil {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   balanceTransaction.Source.Transfer.Destination.ID,
				ConnectorID: connectorID,
			}
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Transfer.BalanceTransaction.ID,
					Type:      models.PaymentTypeTransfer,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Amount:    big.NewInt(balanceTransaction.Amount),
			Status:    models.PaymentStatusFailed,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypeAdjustment:
		if balanceTransaction.Source.Dispute == nil {
			// We are only handle dispute adjustments
			return ingestion.PaymentBatchElement{}, false
		}

		transactionCurrency := strings.ToUpper(string(balanceTransaction.Source.Dispute.Charge.Currency))
		_, ok := supportedCurrenciesWithDecimal[transactionCurrency]
		if !ok {
			return ingestion.PaymentBatchElement{}, false
		}

		disputeStatus := convertDisputeStatus(balanceTransaction.Source.Dispute.Status)
		paymentStatus := models.PaymentStatusPending
		switch disputeStatus {
		case models.PaymentStatusDisputeWon:
			paymentStatus = models.PaymentStatusSucceeded
		case models.PaymentStatusDisputeLost:
			paymentStatus = models.PaymentStatusFailed
		default:
			paymentStatus = models.PaymentStatusPending
		}

		payment = &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Dispute.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			Reference:     balanceTransaction.Source.Dispute.Charge.BalanceTransaction.ID,
			ConnectorID:   connectorID,
			Type:          models.PaymentTypePayIn,
			Status:        paymentStatus, // Dispute is occuring, we don't know the outcome yet
			Amount:        big.NewInt(balanceTransaction.Source.Dispute.Charge.Amount - balanceTransaction.Source.Dispute.Charge.AmountRefunded),
			InitialAmount: big.NewInt(balanceTransaction.Source.Dispute.Charge.Amount),
			Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transactionCurrency),
			RawData:       rawData,
			Scheme:        models.PaymentScheme(balanceTransaction.Source.Dispute.Charge.PaymentMethodDetails.Card.Brand),
			CreatedAt:     time.Unix(balanceTransaction.Source.Dispute.Charge.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference:   account,
				ConnectorID: connectorID,
			}
		}

		adjustment = &models.PaymentAdjustment{
			PaymentID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.Source.Dispute.Charge.BalanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				ConnectorID: connectorID,
			},
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
			Reference: balanceTransaction.ID,
			Status:    disputeStatus,
			RawData:   rawData,
		}

	case stripe.BalanceTransactionTypeStripeFee:
		return ingestion.PaymentBatchElement{}, false
	default:
		return ingestion.PaymentBatchElement{}, false
	}

	return ingestion.PaymentBatchElement{
		Payment:    payment,
		Adjustment: adjustment,
	}, true
}

func convertDisputeStatus(status stripe.DisputeStatus) models.PaymentStatus {
	switch status {
	case stripe.DisputeStatusNeedsResponse, stripe.DisputeStatusUnderReview:
		return models.PaymentStatusDispute
	case stripe.DisputeStatusLost:
		return models.PaymentStatusDisputeLost
	case stripe.DisputeStatusWon:
		return models.PaymentStatusDisputeWon
	default:
		return models.PaymentStatusDispute
	}
}

func convertPayoutStatus(status stripe.PayoutStatus) models.PaymentStatus {
	switch status {
	case stripe.PayoutStatusCanceled:
		return models.PaymentStatusCancelled
	case stripe.PayoutStatusFailed:
		return models.PaymentStatusFailed
	case stripe.PayoutStatusInTransit, stripe.PayoutStatusPending:
		return models.PaymentStatusPending
	case stripe.PayoutStatusPaid:
		return models.PaymentStatusSucceeded
	}

	return models.PaymentStatusOther
}
