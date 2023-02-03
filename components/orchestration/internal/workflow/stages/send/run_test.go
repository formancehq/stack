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

var (
	paymentToWallet = stagestesting.WorkflowTestCase[Send]{
		Name: "payment to wallet",
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
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: sdk.PtrString(paymentAccountName("payment1")),
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: map[string]interface{}{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
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
							Metadata: map[string]interface{}{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{nil},
			},
		},
	}
	paymentToAccount = stagestesting.WorkflowTestCase[Send]{
		Name: "payment to account",
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
						Ledger: internalLedger,
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: sdk.PtrString(paymentAccountName("payment1")),
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: map[string]interface{}{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
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
							Metadata: map[string]interface{}{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
		},
	}
	paymentToAccountWithAlreadyUsedPayment = stagestesting.WorkflowTestCase[Send]{
		Name: "payment to account with already used payment",
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
						Ledger: internalLedger,
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: sdk.PtrString(paymentAccountName("payment1")),
						},
					},
				},
				Returns: []any{nil, activities.ErrTransactionReferenceConflict},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: map[string]interface{}{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
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
							Metadata: map[string]interface{}{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
		},
	}
	accountToAccount = stagestesting.WorkflowTestCase[Send]{
		Name: "account to account",
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
	accountToAccountMixedLedger = stagestesting.WorkflowTestCase[Send]{
		Name: "account to account mixed ledger",
		Stage: Send{
			Source: NewSource().WithAccount(&LedgerAccountSource{
				ID:     "account1",
				Ledger: "ledger1",
			}),
			Destination: NewDestination().WithAccount(&LedgerAccountDestination{
				ID:     "account2",
				Ledger: "ledger2",
			}),
			Amount: *sdk.NewMonetary("USD", 100),
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger1",
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: "world",
								Source:      "account1",
							}},
							Metadata: map[string]interface{}{
								moveToLedgerMetadata: "ledger2",
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger2",
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: "account2",
								Source:      "world",
							}},
							Metadata: map[string]interface{}{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
		},
	}
	accountToWallet = stagestesting.WorkflowTestCase[Send]{
		Name: "account to wallet",
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
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "bar",
				}},
				Returns: []any{&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Ledger: "default",
					},
				}, nil},
			},
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
	accountToWalletMixedLedger = stagestesting.WorkflowTestCase[Send]{
		Name: "account to wallet mixed ledger",
		Stage: Send{
			Source: NewSource().WithAccount(&LedgerAccountSource{
				ID:     "account1",
				Ledger: "ledger1",
			}),
			Destination: NewDestination().WithWallet(&WalletDestination{
				ID:      "wallet",
				Balance: "main",
			}),
			Amount: *sdk.NewMonetary("USD", 100),
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet",
				}},
				Returns: []any{&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Ledger: "ledger2",
					},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger1",
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: "world",
								Source:      "account1",
							}},
							Metadata: map[string]interface{}{
								moveToLedgerMetadata: "ledger2",
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "wallet",
						Data: sdk.CreditWalletRequest{
							Amount: *sdk.NewMonetary("USD", 100),
							Sources: []sdk.Subject{{
								LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", "world"),
							}},
							Balance: sdk.PtrString("main"),
							Metadata: map[string]interface{}{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{nil},
			},
		},
	}
	accountToPayment = stagestesting.WorkflowTestCase[Send]{
		Name: "account to payment",
		Stage: Send{
			Source: NewSource().WithAccount(&LedgerAccountSource{
				ID:     "foo",
				Ledger: "default",
			}),
			Destination: NewDestination().WithPayment(&PaymentDestination{
				PSP:      "stripe",
				Metadata: "stripeConnectID",
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
	walletToAccount = stagestesting.WorkflowTestCase[Send]{
		Name: "wallet to account",
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
						Ledger: "default",
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
	walletToAccountMixedLedger = stagestesting.WorkflowTestCase[Send]{
		Name: "wallet to account mixed ledger",
		Stage: Send{
			Source: NewSource().WithWallet(&WalletSource{
				ID:      "wallet",
				Balance: "main",
			}),
			Destination: NewDestination().WithAccount(&LedgerAccountDestination{
				ID:     "account",
				Ledger: "ledger2",
			}),
			Amount: *sdk.NewMonetary("USD", 100),
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet",
				}},
				Returns: []any{&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Id:     "wallet",
						Ledger: "ledger1",
					},
				}, nil},
			},
			{
				Activity: activities.DebitWalletActivity,
				Args: []any{
					mock.Anything, activities.DebitWalletRequest{
						ID: "wallet",
						Data: sdk.DebitWalletRequest{
							Amount: sdk.Monetary{
								Asset:  "USD",
								Amount: 100,
							},
							Balances: []string{"main"},
							Metadata: map[string]interface{}{
								moveToLedgerMetadata: "ledger2",
							},
						},
					},
				},
				Returns: []any{nil, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger2",
						Data: sdk.PostTransaction{
							Postings: []sdk.Posting{{
								Amount:      100,
								Asset:       "USD",
								Destination: "account",
								Source:      "world",
							}},
							Metadata: map[string]interface{}{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{&sdk.TransactionsResponse{
					Data: []sdk.Transaction{{}},
				}, nil},
			},
		},
	}
	walletToWallet = stagestesting.WorkflowTestCase[Send]{
		Name: "wallet to wallet",
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
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "foo",
				}},
				Returns: []any{&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Ledger: "default",
					},
				}, nil},
			},
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "bar",
				}},
				Returns: []any{&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Ledger: "default",
					},
				}, nil},
			},
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
	walletToWalletMixedLedger = stagestesting.WorkflowTestCase[Send]{
		Name: "wallet to wallet mixed ledger",
		Stage: Send{
			Source: NewSource().WithWallet(&WalletSource{
				ID:      "wallet1",
				Balance: "main",
			}),
			Destination: NewDestination().WithWallet(&WalletDestination{
				ID:      "wallet2",
				Balance: "main",
			}),
			Amount: *sdk.NewMonetary("USD", 100),
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet1",
				}},
				Returns: []any{&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Ledger: "ledger1",
					},
				}, nil},
			},
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet2",
				}},
				Returns: []any{&sdk.GetWalletResponse{
					Data: sdk.WalletWithBalances{
						Ledger: "ledger2",
					},
				}, nil},
			},
			{
				Activity: activities.DebitWalletActivity,
				Args: []any{
					mock.Anything, activities.DebitWalletRequest{
						ID: "wallet1",
						Data: sdk.DebitWalletRequest{
							Amount: *sdk.NewMonetary("USD", 100),
							Metadata: map[string]interface{}{
								moveToLedgerMetadata: "ledger2",
							},
							Balances: []string{"main"},
						},
					},
				},
				Returns: []any{nil, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "wallet2",
						Data: sdk.CreditWalletRequest{
							Amount:  *sdk.NewMonetary("USD", 100),
							Balance: sdk.PtrString("main"),
							Metadata: map[string]interface{}{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{nil},
			},
		},
	}
	walletToPayment = stagestesting.WorkflowTestCase[Send]{
		Name: "wallet to payment",
		Stage: Send{
			Source: NewSource().WithWallet(&WalletSource{
				ID:      "foo",
				Balance: "main",
			}),
			Destination: NewDestination().WithPayment(&PaymentDestination{
				PSP:      "stripe",
				Metadata: "stripeConnectID",
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
)

var testCases = []stagestesting.WorkflowTestCase[Send]{
	paymentToWallet,
	paymentToAccount,
	paymentToAccountWithAlreadyUsedPayment,
	accountToAccount,
	accountToAccountMixedLedger,
	accountToWallet,
	accountToWalletMixedLedger,
	accountToPayment,
	walletToAccount,
	walletToAccountMixedLedger,
	walletToWallet,
	walletToWalletMixedLedger,
	walletToPayment,
}

func TestSend(t *testing.T) {
	stagestesting.RunWorkflows(t, testCases...)
}
