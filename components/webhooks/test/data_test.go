package test_test

import (
	"time"

	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/events/payments"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	EventTypeCommittedTransactions = "committed_transactions"
	EventTypeSavedMetadata         = "saved_metadata"
	EventTypeRevertedTransaction   = "reverted_transaction"
	EventTypeSavedPayments         = "saved_payment"
	EventTypeSavedAccounts         = "saved_account"
	EventTypeConnectorReset        = "connector_reset"
)

var (
	app1   = "app1"
	type1  = "connector_reset"
	event1 = &events.Event{
		CreatedAt: timestamppb.New(time.Now().UTC()),
		App:       app1,
		Event: &events.Event_ResetConnector{
			ResetConnector: &payments.ResetConnector{
				CreatedAt: timestamppb.New(time.Now().UTC()),
				Provider:  "BANKING_CIRCLE",
			},
		},
	}

	app2   = "app2"
	type2  = "saved_payment"
	event2 = &events.Event{
		CreatedAt: timestamppb.New(time.Now().UTC()),
		App:       app2,
		Event: &events.Event_PaymentSaved{
			PaymentSaved: &payments.PaymentSaved{
				Id:        "123456789",
				Reference: "123",
				CreatedAt: timestamppb.New(time.Now().UTC()),
				Provider:  "BANKING_CIRCLE",
				Type:      "PAY_IN",
				Status:    "PENDING",
				Scheme:    "GOOGLE_PAY",
				Asset:     "USD/2",
				Amount:    100,
			},
		},
	}

	app3   = "app3"
	event3 = &events.Event{
		CreatedAt: timestamppb.New(time.Now().UTC()),
		App:       app3,
		Event: &events.Event_AccountSaved{
			AccountSaved: &payments.AccountSaved{
				Accounts: []*payments.Account{
					{
						Id:        "123456789",
						CreatedAt: timestamppb.New(time.Now().UTC()),
						Reference: "123",
						Provider:  "BANKING_CIRCLE",
						Type:      "SOURCE",
					},
				},
			},
		},
	}
)
