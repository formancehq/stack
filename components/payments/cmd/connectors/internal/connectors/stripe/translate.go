package stripe

import (
	"encoding/json"
	"log"
	"math/big"
	"runtime/debug"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/internal/models"
	"github.com/stripe/stripe-go/v72"
)

func CreateBatchElement(balanceTransaction *stripe.BalanceTransaction, account string, forward bool) (ingestion.PaymentBatchElement, bool) {
	var payment models.Payment // reference   payments.Referenced
	// paymentData *payments.Data
	// adjustment  *payments.Adjustment

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

	if balanceTransaction.Source.Payout == nil && balanceTransaction.Source.Charge == nil {
		return ingestion.PaymentBatchElement{}, false
	}

	rawData, err := json.Marshal(balanceTransaction)
	if err != nil {
		return ingestion.PaymentBatchElement{}, false
	}

	switch balanceTransaction.Type {
	case stripe.BalanceTransactionTypeCharge:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayIn,
			Status:    models.PaymentStatusSucceeded,
			Amount:    big.NewInt(balanceTransaction.Source.Charge.Amount),
			Asset:     currency.FormatAsset(string(balanceTransaction.Source.Charge.Currency)),
			RawData:   rawData,
			Scheme:    models.PaymentScheme(balanceTransaction.Source.Charge.PaymentMethodDetails.Card.Brand),
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypePayout:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayOut,
			Status:    convertPayoutStatus(balanceTransaction.Source.Payout.Status),
			Amount:    big.NewInt(balanceTransaction.Source.Payout.Amount),
			RawData:   rawData,
			Asset:     currency.FormatAsset(string(balanceTransaction.Source.Payout.Currency)),
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
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypeTransfer:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayOut,
			Status:    models.PaymentStatusSucceeded,
			Amount:    big.NewInt(balanceTransaction.Source.Transfer.Amount),
			RawData:   rawData,
			Asset:     currency.FormatAsset(string(balanceTransaction.Source.Transfer.Currency)),
			Scheme:    models.PaymentSchemeOther,
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypeRefund:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayOut,
			Adjustments: []*models.Adjustment{
				{
					Reference: balanceTransaction.ID,
					Status:    models.PaymentStatusSucceeded,
					Amount:    balanceTransaction.Amount,
					CreatedAt: time.Unix(balanceTransaction.Created, 0),
					RawData:   rawData,
				},
			},
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypePayment:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayIn,
			Status:    models.PaymentStatusSucceeded,
			Amount:    big.NewInt(balanceTransaction.Source.Charge.Amount),
			RawData:   rawData,
			Asset:     currency.FormatAsset(string(balanceTransaction.Source.Charge.Currency)),
			Scheme:    models.PaymentSchemeOther,
			CreatedAt: time.Unix(balanceTransaction.Created, 0),
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypePayoutCancel:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayOut,
			Status:    models.PaymentStatusFailed,
			Adjustments: []*models.Adjustment{
				{
					Reference: balanceTransaction.ID,
					Status:    convertPayoutStatus(balanceTransaction.Source.Payout.Status),
					CreatedAt: time.Unix(balanceTransaction.Created, 0),
					RawData:   rawData,
					Absolute:  true,
				},
			},
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypePayoutFailure:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayIn,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayIn,
			Status:    models.PaymentStatusFailed,
			Adjustments: []*models.Adjustment{
				{
					Reference: balanceTransaction.ID,
					Status:    convertPayoutStatus(balanceTransaction.Source.Payout.Status),
					CreatedAt: time.Unix(balanceTransaction.Created, 0),
					RawData:   rawData,
					Absolute:  true,
				},
			},
		}

		if account != "" {
			payment.DestinationAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypePaymentRefund:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayOut,
			Status:    models.PaymentStatusSucceeded,
			Adjustments: []*models.Adjustment{
				{
					Reference: balanceTransaction.ID,
					Status:    models.PaymentStatusSucceeded,
					Amount:    balanceTransaction.Amount,
					CreatedAt: time.Unix(balanceTransaction.Created, 0),
					RawData:   rawData,
				},
			},
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypeAdjustment:
		payment = models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: balanceTransaction.ID,
					Type:      models.PaymentTypePayOut,
				},
				Provider: models.ConnectorProviderStripe,
			},
			Reference: balanceTransaction.ID,
			Type:      models.PaymentTypePayOut,
			Adjustments: []*models.Adjustment{
				{
					Reference: balanceTransaction.ID,
					Status:    models.PaymentStatusCancelled,
					Amount:    balanceTransaction.Amount,
					CreatedAt: time.Unix(balanceTransaction.Created, 0),
					RawData:   rawData,
				},
			},
		}

		if account != "" {
			payment.SourceAccountID = &models.AccountID{
				Reference: account,
				Provider:  models.ConnectorProviderStripe,
			}
		}

	case stripe.BalanceTransactionTypeStripeFee:
		return ingestion.PaymentBatchElement{}, false
	default:
		return ingestion.PaymentBatchElement{}, false
	}

	return ingestion.PaymentBatchElement{
		Payment: &payment,
		Update:  forward,
	}, true
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
