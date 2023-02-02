package send

import (
	"testing"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal/stagestesting"
	"github.com/stretchr/testify/mock"
)

func TestSendSchemaValidation(t *testing.T) {
	stagestesting.TestSchemas(t, "send", []stagestesting.SchemaTestCase{
		{
			Name: "twice destination",
			Data: map[string]any{
				"source": map[string]any{
					"account": map[string]any{
						"id": "bar",
					},
				},
				"destination": map[string]any{
					"wallet": map[string]any{
						"id": "foo",
					},
					"account": map[string]any{
						"id": "foo",
					},
				},
				"amount": map[string]any{
					"amount": float64(100),
					"asset":  "USD",
				},
			},
			ExpectedResolved: Send{
				Source: Source{
					Account: &LedgerAccountSource{
						ID:     "bar",
						Ledger: "default",
					},
				},
				Destination: Destination{
					Wallet: &WalletDestination{
						ID:      "foo",
						Balance: "main",
					},
					Account: &LedgerAccountSource{
						ID:     "foo",
						Ledger: "default",
					},
				},
				Amount: sdk.Monetary{
					Amount: 100,
					Asset:  "USD",
				},
			},
			ExpectedValidationError: true,
		},
		{
			Name: "valid case",
			Data: map[string]any{
				"source": map[string]any{
					"account": map[string]any{
						"id": "bar",
					},
				},
				"destination": map[string]any{
					"wallet": map[string]any{
						"id": "foo",
					},
				},
				"amount": map[string]any{
					"amount": float64(100),
					"asset":  "USD",
				},
			},
			ExpectedResolved: Send{
				Source: Source{
					Account: &LedgerAccountSource{
						ID:     "bar",
						Ledger: "default",
					},
				},
				Destination: Destination{
					Wallet: &WalletSource{
						ID:      "foo",
						Balance: "main",
					},
				},
				Amount: sdk.Monetary{
					Asset:  "USD",
					Amount: 100,
				},
			},
		},
	}...)
}

var paymentToWallet = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithPayment(&PaymentSource{
			ID: "payment1",
		}),
		Destination: NewDestination().WithWallet(&WalletDestination{
			ID:      "wallet1",
			Balance: "main",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.GetPaymentActivity,
			Args: []any{mock.Anything, activities.GetPaymentRequest{
				ID: "payment1",
			}},
			Returns: []any{
				&sdk.PaymentResponse{
					Data: sdk.Payment{
						InitialAmount: 100,
						Asset:         "USD",
						Provider:      sdk.STRIPE,
						Status:        sdk.SUCCEEDED,
					},
				}, nil,
			},
		},
		{
			Activity: activities.GetWalletActivity,
			Args: []any{mock.Anything, activities.GetWalletRequest{
				ID: "wallet1",
			}},
			Returns: []any{
				&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Id:     "wallet1",
						Ledger: "default",
					},
				}, nil,
			},
		},
		{
			Activity: activities.CreditWalletActivity,
			Args: []any{
				mock.Anything, activities.CreditWalletRequest{
					ID: "wallet1",
					Data: sdk.CreditWalletRequest{
						Amount: *sdk.NewMonetary("USD", 100),
						Sources: []sdk.Subject{{
							LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", "world"),
						}},
						Balance: sdk.PtrString("main"),
					},
				},
			},
			Returns: []any{nil},
		},
	},
}

var paymentToAccount = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithPayment(&PaymentSource{
			ID: "payment1",
		}),
		Destination: NewDestination().WithAccount(&LedgerAccountDestination{
			ID:     "foo",
			Ledger: "default",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.GetPaymentActivity,
			Args: []any{mock.Anything, activities.GetPaymentRequest{
				ID: "payment1",
			}},
			Returns: []any{
				&sdk.PaymentResponse{
					Data: sdk.Payment{
						InitialAmount: 100,
						Asset:         "USD",
						Provider:      sdk.STRIPE,
						Status:        sdk.SUCCEEDED,
					},
				}, nil,
			},
		},
		{
			Activity: activities.CreateTransactionActivity,
			Args: []any{
				mock.Anything, activities.CreateTransactionRequest{
					Ledger: "default",
					Data: sdk.PostTransaction{
						Postings: []sdk.Posting{{
							Amount:      100,
							Asset:       "USD",
							Destination: "foo",
							Source:      "world",
						}},
					},
				},
			},
			Returns: []any{&sdk.TransactionsResponse{
				Data: []sdk.Transaction{{}},
			}, nil},
		},
	},
}

var accountToAccount = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithAccount(&LedgerAccountSource{
			ID:     "foo",
			Ledger: "default",
		}),
		Destination: NewDestination().WithAccount(&LedgerAccountDestination{
			ID:     "bar",
			Ledger: "default",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.CreateTransactionActivity,
			Args: []any{
				mock.Anything, activities.CreateTransactionRequest{
					Ledger: "default",
					Data: sdk.PostTransaction{
						Postings: []sdk.Posting{{
							Amount:      100,
							Asset:       "USD",
							Destination: "bar",
							Source:      "foo",
						}},
					},
				},
			},
			Returns: []any{&sdk.TransactionsResponse{
				Data: []sdk.Transaction{{}},
			}, nil},
		},
	},
}

var accountToWallet = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithAccount(&LedgerAccountSource{
			ID:     "foo",
			Ledger: "default",
		}),
		Destination: NewDestination().WithWallet(&WalletDestination{
			ID:      "bar",
			Balance: "main",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.CreditWalletActivity,
			Args: []any{
				mock.Anything, activities.CreditWalletRequest{
					ID: "bar",
					Data: sdk.CreditWalletRequest{
						Amount: *sdk.NewMonetary("USD", 100),
						Sources: []sdk.Subject{{
							LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", "foo"),
						}},
						Balance: sdk.PtrString("main"),
					},
				},
			},
			Returns: []any{nil},
		},
	},
}

var accountToPayment = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithAccount(&LedgerAccountSource{
			ID:     "foo",
			Ledger: "default",
		}),
		Destination: NewDestination().WithPayment(&PaymentDestination{
			PSP: "stripe",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.GetAccountActivity,
			Args: []any{mock.Anything, activities.GetAccountRequest{
				Ledger: "default",
				ID:     "foo",
			}},
			Returns: []any{&sdk.AccountResponse{
				Data: sdk.AccountWithVolumesAndBalances{
					Address: "foo",
					Metadata: map[string]interface{}{
						"stripeConnectID": "abcd",
					},
				},
			}, nil},
		},
		{
			Activity: activities.StripeTransferActivity,
			Args: []any{
				mock.Anything, sdk.StripeTransferRequest{
					Amount:      sdk.PtrInt64(100),
					Asset:       sdk.PtrString("USD"),
					Destination: sdk.PtrString("abcd"),
				},
			},
			Returns: []any{nil},
		},
		{
			Activity: activities.CreateTransactionActivity,
			Args: []any{
				mock.Anything, activities.CreateTransactionRequest{
					Ledger: "default",
					Data: sdk.PostTransaction{
						Postings: []sdk.Posting{{
							Amount:      100,
							Asset:       "USD",
							Destination: "world",
							Source:      "foo",
						}},
					},
				},
			},
			Returns: []any{&sdk.TransactionsResponse{
				Data: []sdk.Transaction{{}},
			}, nil},
		},
	},
}

var walletToAccount = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithWallet(&WalletSource{
			ID:      "foo",
			Balance: "main",
		}),
		Destination: NewDestination().WithAccount(&LedgerAccountDestination{
			ID:     "bar",
			Ledger: "default",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.GetWalletActivity,
			Args: []any{mock.Anything, activities.GetWalletRequest{
				ID: "foo",
			}},
			Returns: []any{&sdk.GetWalletResponse{
				Data: sdk.WalletWithBalances{
					Id: "foo",
					Metadata: map[string]interface{}{
						"stripeConnectID": "abcd",
					},
				},
			}, nil},
		},
		{
			Activity: activities.DebitWalletActivity,
			Args: []any{
				mock.Anything, activities.DebitWalletRequest{
					ID: "foo",
					Data: sdk.DebitWalletRequest{
						Amount: sdk.Monetary{
							Asset:  "USD",
							Amount: 100,
						},
						Destination: &sdk.Subject{
							LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", "bar"),
						},
						Balances: []string{"main"},
					},
				},
			},
			Returns: []any{nil, nil},
		},
	},
}

var walletToWallet = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithWallet(&WalletSource{
			ID:      "foo",
			Balance: "main",
		}),
		Destination: NewDestination().WithWallet(&WalletDestination{
			ID:      "bar",
			Balance: "main",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.CreditWalletActivity,
			Args: []any{
				mock.Anything, activities.CreditWalletRequest{
					ID: "bar",
					Data: sdk.CreditWalletRequest{
						Amount: *sdk.NewMonetary("USD", 100),
						Sources: []sdk.Subject{{
							WalletSubject: &sdk.WalletSubject{
								Type:       "WALLET",
								Identifier: "foo",
								Balance:    sdk.PtrString("main"),
							},
						}},
						Balance: sdk.PtrString("main"),
					},
				},
			},
			Returns: []any{nil},
		},
	},
}

var walletToPayment = stagestesting.WorkflowTestCase[Send]{
	Stage: Send{
		Source: NewSource().WithWallet(&WalletSource{
			ID:      "foo",
			Balance: "main",
		}),
		Destination: NewDestination().WithPayment(&PaymentDestination{
			PSP: "stripe",
		}),
		Amount: *sdk.NewMonetary("USD", 100),
	},
	MockedActivities: []stagestesting.MockedActivity{
		{
			Activity: activities.GetWalletActivity,
			Args: []any{mock.Anything, activities.GetWalletRequest{
				ID: "foo",
			}},
			Returns: []any{&sdk.GetWalletResponse{
				Data: sdk.WalletWithBalances{
					Metadata: map[string]interface{}{
						"stripeConnectID": "abcd",
					},
				},
			}, nil},
		},
		{
			Activity: activities.StripeTransferActivity,
			Args: []any{
				mock.Anything, sdk.StripeTransferRequest{
					Amount:      sdk.PtrInt64(100),
					Asset:       sdk.PtrString("USD"),
					Destination: sdk.PtrString("abcd"),
				},
			},
			Returns: []any{nil},
		},
		{
			Activity: activities.DebitWalletActivity,
			Args: []any{
				mock.Anything, activities.DebitWalletRequest{
					ID: "foo",
					Data: sdk.DebitWalletRequest{
						Amount: sdk.Monetary{
							Asset:  "USD",
							Amount: 100,
						},
						Balances: []string{"main"},
					},
				},
			},
			Returns: []any{nil, nil},
		},
	},
}

var testCases = []stagestesting.WorkflowTestCase[Send]{
	paymentToWallet,
	paymentToAccount,
	accountToAccount,
	accountToWallet,
	accountToPayment,
	walletToAccount,
	walletToWallet,
	walletToPayment,
}

func TestSend(t *testing.T) {
	stagestesting.RunWorkflows(t, testCases...)
}
