package send

import (
	"math/big"
	"testing"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal/stagestesting"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/temporal"
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
				Amount: &shared.Monetary{
					Amount: big.NewInt(100),
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
				Amount: &shared.Monetary{
					Asset:  "USD",
					Amount: big.NewInt(100),
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetPaymentActivity,
				Args: []any{mock.Anything, activities.GetPaymentRequest{
					ID: "payment1",
				}},
				Returns: []any{
					&shared.PaymentResponse{
						Data: shared.Payment{
							InitialAmount: big.NewInt(100),
							Asset:         "USD",
							Provider:      shared.ConnectorStripe,
							Status:        shared.PaymentStatusSucceeded,
							Scheme:        shared.PaymentSchemeUnknown,
							Type:          shared.PaymentTypeOther,
						},
					}, nil,
				},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: pointer.For(paymentAccountName("payment1")),
							Metadata:  metadata.Metadata{},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: metadata.Metadata{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
				}, nil},
			},
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet1",
				}},
				Returns: []any{
					&shared.GetWalletResponse{
						Data: shared.WalletWithBalances{
							ID:     "wallet1",
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
						Data: &shared.CreditWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD",
							},
							Sources: []shared.Subject{{
								LedgerAccountSubject: &shared.LedgerAccountSubject{
									Identifier: "world",
									Type:       "ACCOUNT",
								},
								Type: shared.SubjectTypeAccount,
							}},
							Balance: pointer.For("main"),
							Metadata: metadata.Metadata{
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetPaymentActivity,
				Args: []any{mock.Anything, activities.GetPaymentRequest{
					ID: "payment1",
				}},
				Returns: []any{
					&shared.PaymentResponse{
						Data: shared.Payment{
							InitialAmount: big.NewInt(100),
							Asset:         "USD",
							Provider:      shared.ConnectorStripe,
							Status:        shared.PaymentStatusSucceeded,
							Scheme:        shared.PaymentSchemeUnknown,
							Type:          shared.PaymentTypeOther,
						},
					}, nil,
				},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: pointer.For(paymentAccountName("payment1")),
							Metadata:  metadata.Metadata{},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: metadata.Metadata{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "foo",
								Source:      "world",
							}},
							Metadata: metadata.Metadata{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetPaymentActivity,
				Args: []any{mock.Anything, activities.GetPaymentRequest{
					ID: "payment1",
				}},
				Returns: []any{
					&shared.PaymentResponse{
						Data: shared.Payment{
							InitialAmount: big.NewInt(100),
							Asset:         "USD",
							Provider:      shared.ConnectorStripe,
							Status:        shared.PaymentStatusSucceeded,
							Scheme:        shared.PaymentSchemeUnknown,
							Type:          shared.PaymentTypeOther,
						},
					}, nil,
				},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: pointer.For(paymentAccountName("payment1")),
							Metadata:  metadata.Metadata{},
						},
					},
				},
				Returns: []any{nil, temporal.NewApplicationError("", "CONFLICT", "")},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: metadata.Metadata{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "foo",
								Source:      "world",
							}},
							Metadata: metadata.Metadata{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "bar",
								Source:      "foo",
							}},
							Metadata: metadata.Metadata{},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger1",
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      "account1",
							}},
							Metadata: metadata.Metadata{
								moveToLedgerMetadata: "ledger2",
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger2",
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "account2",
								Source:      "world",
							}},
							Metadata: metadata.Metadata{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "bar",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						Ledger: "default",
					},
				}, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "bar",
						Data: &shared.CreditWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD",
							},
							Sources: []shared.Subject{{
								LedgerAccountSubject: &shared.LedgerAccountSubject{
									Identifier: "foo",
									Type:       "ACCOUNT",
								},
								Type: shared.SubjectTypeAccount,
							}},
							Balance: pointer.For("main"),
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						Ledger: "ledger2",
					},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger1",
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      "account1",
							}},
							Metadata: metadata.Metadata{
								moveToLedgerMetadata: "ledger2",
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "wallet",
						Data: &shared.CreditWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD",
							},
							Sources: []shared.Subject{{
								LedgerAccountSubject: &shared.LedgerAccountSubject{
									Identifier: "world",
									Type:       "ACCOUNT",
								},
								Type: shared.SubjectTypeAccount,
							}},
							Balance: pointer.For("main"),
							Metadata: metadata.Metadata{
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetAccountActivity,
				Args: []any{mock.Anything, activities.GetAccountRequest{
					Ledger: "default",
					ID:     "foo",
				}},
				Returns: []any{&shared.AccountResponse{
					Data: shared.AccountWithVolumesAndBalances{
						Address: "foo",
						Metadata: metadata.Metadata{
							"stripeConnectID": "abcd",
						},
					},
				}, nil},
			},
			{
				Activity: activities.StripeTransferActivity,
				Args: []any{
					mock.Anything, shared.StripeTransferRequest{
						Amount:      big.NewInt(100),
						Asset:       pointer.For("USD"),
						Destination: pointer.For("abcd"),
					},
				},
				Returns: []any{nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      "foo",
							}},
							Metadata: metadata.Metadata{},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "foo",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						ID: "foo",
						Metadata: metadata.Metadata{
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
						Data: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Destination: &shared.Subject{
								LedgerAccountSubject: &shared.LedgerAccountSubject{
									Identifier: "bar",
									Type:       "ACCOUNT",
								},
								Type: shared.SubjectTypeAccount,
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						ID:     "wallet",
						Ledger: "ledger1",
					},
				}, nil},
			},
			{
				Activity: activities.DebitWalletActivity,
				Args: []any{
					mock.Anything, activities.DebitWalletRequest{
						ID: "wallet",
						Data: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Balances: []string{"main"},
							Metadata: metadata.Metadata{
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
						Data: shared.PostTransaction{
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "account",
								Source:      "world",
							}},
							Metadata: metadata.Metadata{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{&shared.CreateTransactionResponse{
					Data: shared.Transaction{},
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "foo",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						Ledger: "default",
					},
				}, nil},
			},
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "bar",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						Ledger: "default",
					},
				}, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "bar",
						Data: &shared.CreditWalletRequest{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Sources: []shared.Subject{{
								WalletSubject: &shared.WalletSubject{
									Type:       "WALLET",
									Identifier: "foo",
									Balance:    pointer.For("main"),
								},
								Type: shared.SubjectTypeWallet,
							}},
							Balance:  pointer.For("main"),
							Metadata: map[string]string{},
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet1",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						Ledger: "ledger1",
					},
				}, nil},
			},
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "wallet2",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						Ledger: "ledger2",
					},
				}, nil},
			},
			{
				Activity: activities.DebitWalletActivity,
				Args: []any{
					mock.Anything, activities.DebitWalletRequest{
						ID: "wallet1",
						Data: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Metadata: metadata.Metadata{
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
						Data: &shared.CreditWalletRequest{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Balance: pointer.For("main"),
							Metadata: metadata.Metadata{
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
			Amount: &shared.Monetary{
				Amount: big.NewInt(100),
				Asset:  "USD",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.GetWalletActivity,
				Args: []any{mock.Anything, activities.GetWalletRequest{
					ID: "foo",
				}},
				Returns: []any{&shared.GetWalletResponse{
					Data: shared.WalletWithBalances{
						Metadata: metadata.Metadata{
							"stripeConnectID": "abcd",
						},
					},
				}, nil},
			},
			{
				Activity: activities.StripeTransferActivity,
				Args: []any{
					mock.Anything, shared.StripeTransferRequest{
						Amount:      big.NewInt(100),
						Asset:       pointer.For("USD"),
						Destination: pointer.For("abcd"),
					},
				},
				Returns: []any{nil},
			},
			{
				Activity: activities.DebitWalletActivity,
				Args: []any{
					mock.Anything, activities.DebitWalletRequest{
						ID: "foo",
						Data: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
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
