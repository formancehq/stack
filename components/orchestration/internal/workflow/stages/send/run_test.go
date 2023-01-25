package send

import (
	"fmt"
	"testing"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/orchestration/internal/schema"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestSendSchemaValidation(t *testing.T) {
	type testCase struct {
		data          map[string]any
		expectedError bool
	}
	testCases := []testCase{
		{
			data: map[string]any{
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
			expectedError: true,
		},
		{
			data: map[string]any{
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
		},
	}
	for _, testCase := range testCases {
		s, err := schema.Resolve(schema.Context{
			Variables: map[string]string{},
		}, testCase.data, "send")
		require.NoError(t, err, "resolving schema")
		err = schema.ValidateRequirements(s)
		if testCase.expectedError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

type mockedActivity struct {
	activity any
	args     []any
	returns  []any
}

type testCase struct {
	source         Source
	destination    Destination
	amount         sdk.Monetary
	mockedActivity []mockedActivity
	name           string
}

var paymentToWallet = testCase{
	source: NewSource().WithPayment(&PaymentSource{
		ID: "payment1",
	}),
	destination: NewDestination().WithWallet(&WalletDestination{
		ID:      "wallet1",
		Balance: "main",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.GetPaymentActivity,
			args: []any{mock.Anything, activities.GetPaymentRequest{
				ID: "payment1",
			}},
			returns: []any{
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
			activity: activities.GetWalletActivity,
			args: []any{mock.Anything, activities.GetWalletRequest{
				ID: "wallet1",
			}},
			returns: []any{
				&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Id:     "wallet1",
						Ledger: "default",
					},
				}, nil,
			},
		},
		{
			activity: activities.CreditWalletActivity,
			args: []any{
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
			returns: []any{nil},
		},
	},
}

var paymentToAccount = testCase{
	source: NewSource().WithPayment(&PaymentSource{
		ID: "payment1",
	}),
	destination: NewDestination().WithAccount(&LedgerAccountDestination{
		ID:     "foo",
		Ledger: "default",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.GetPaymentActivity,
			args: []any{mock.Anything, activities.GetPaymentRequest{
				ID: "payment1",
			}},
			returns: []any{
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
			activity: activities.CreateTransactionActivity,
			args: []any{
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
			returns: []any{&sdk.TransactionsResponse{
				Data: []sdk.Transaction{{}},
			}, nil},
		},
	},
}

var accountToAccount = testCase{
	source: NewSource().WithAccount(&LedgerAccountSource{
		ID:     "foo",
		Ledger: "default",
	}),
	destination: NewDestination().WithAccount(&LedgerAccountDestination{
		ID:     "bar",
		Ledger: "default",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.CreateTransactionActivity,
			args: []any{
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
			returns: []any{&sdk.TransactionsResponse{
				Data: []sdk.Transaction{{}},
			}, nil},
		},
	},
}

var accountToWallet = testCase{
	source: NewSource().WithAccount(&LedgerAccountSource{
		ID:     "foo",
		Ledger: "default",
	}),
	destination: NewDestination().WithWallet(&WalletDestination{
		ID:      "bar",
		Balance: "main",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.CreditWalletActivity,
			args: []any{
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
			returns: []any{nil},
		},
	},
}

var accountToPayment = testCase{
	source: NewSource().WithAccount(&LedgerAccountSource{
		ID:     "foo",
		Ledger: "default",
	}),
	destination: NewDestination().WithPayment(&PaymentDestination{
		PSP: "stripe",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.GetAccountActivity,
			args: []any{mock.Anything, activities.GetAccountRequest{
				Ledger: "default",
				ID:     "foo",
			}},
			returns: []any{&sdk.AccountResponse{
				Data: sdk.AccountWithVolumesAndBalances{
					Address: "foo",
					Metadata: map[string]interface{}{
						"stripeConnectID": "abcd",
					},
				},
			}, nil},
		},
		{
			activity: activities.StripeTransferActivity,
			args: []any{
				mock.Anything, sdk.StripeTransferRequest{
					Amount:      sdk.PtrInt64(100),
					Asset:       sdk.PtrString("USD"),
					Destination: sdk.PtrString("abcd"),
				},
			},
			returns: []any{nil},
		},
		{
			activity: activities.CreateTransactionActivity,
			args: []any{
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
			returns: []any{&sdk.TransactionsResponse{
				Data: []sdk.Transaction{{}},
			}, nil},
		},
	},
}

var walletToAccount = testCase{
	source: NewSource().WithWallet(&WalletSource{
		ID:      "foo",
		Balance: "main",
	}),
	destination: NewDestination().WithAccount(&LedgerAccountDestination{
		ID:     "bar",
		Ledger: "default",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.GetWalletActivity,
			args: []any{mock.Anything, activities.GetWalletRequest{
				ID: "foo",
			}},
			returns: []any{&sdk.GetWalletResponse{
				Data: sdk.WalletWithBalances{
					Id: "foo",
					Metadata: map[string]interface{}{
						"stripeConnectID": "abcd",
					},
				},
			}, nil},
		},
		{
			activity: activities.DebitWalletActivity,
			args: []any{
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
			returns: []any{nil, nil},
		},
	},
}

var walletToWallet = testCase{
	source: NewSource().WithWallet(&WalletSource{
		ID:      "foo",
		Balance: "main",
	}),
	destination: NewDestination().WithWallet(&WalletDestination{
		ID:      "bar",
		Balance: "main",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.CreditWalletActivity,
			args: []any{
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
			returns: []any{nil},
		},
	},
}

var walletToPayment = testCase{
	source: NewSource().WithWallet(&WalletSource{
		ID:      "foo",
		Balance: "main",
	}),
	destination: NewDestination().WithPayment(&PaymentDestination{
		PSP: "stripe",
	}),
	amount: *sdk.NewMonetary("USD", 100),
	mockedActivity: []mockedActivity{
		{
			activity: activities.GetWalletActivity,
			args: []any{mock.Anything, activities.GetWalletRequest{
				ID: "foo",
			}},
			returns: []any{&sdk.GetWalletResponse{
				Data: sdk.WalletWithBalances{
					Metadata: map[string]interface{}{
						"stripeConnectID": "abcd",
					},
				},
			}, nil},
		},
		{
			activity: activities.StripeTransferActivity,
			args: []any{
				mock.Anything, sdk.StripeTransferRequest{
					Amount:      sdk.PtrInt64(100),
					Asset:       sdk.PtrString("USD"),
					Destination: sdk.PtrString("abcd"),
				},
			},
			returns: []any{nil},
		},
		{
			activity: activities.DebitWalletActivity,
			args: []any{
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
			returns: []any{nil, nil},
		},
	},
}

var testCases = []testCase{
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

	for _, tc := range testCases {
		var (
			from, to string
		)
		switch {
		case tc.source.Wallet != nil:
			from = "wallet"
		case tc.source.Payment != nil:
			from = "payment"
		case tc.source.Account != nil:
			from = "account"
		}
		switch {
		case tc.destination.Wallet != nil:
			to = "wallet"
		case tc.destination.Payment != nil:
			to = "payment"
		case tc.destination.Account != nil:
			to = "account"
		}
		tc := tc
		testName := fmt.Sprintf("%s->%s", from, to)
		if tc.name != "" {
			testName += "/" + tc.name
		}
		t.Run(testName, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}

			env := testSuite.NewTestWorkflowEnvironment()
			for _, ma := range tc.mockedActivity {
				env.OnActivity(ma.activity, ma.args...).Return(ma.returns...)
			}

			send := Send{
				Source:      tc.source,
				Destination: tc.destination,
				Amount:      tc.amount,
			}

			env.ExecuteWorkflow(RunSend, send)
			require.True(t, env.IsWorkflowCompleted())
			require.NoError(t, env.GetWorkflowError())
		})
	}
}
