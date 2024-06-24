package send

import (
	"math/big"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal/stagestesting"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/temporal"
)

func TestSendSchemaValidation(t *testing.T) {
	now := time.Now().Round(time.Millisecond).UTC()
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
						WalletReference: WalletReference{
							ID: "foo",
						},
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
			Name: "invalid wallet reference",
			Data: map[string]any{
				"source": map[string]any{
					"account": map[string]any{
						"id": "bar",
					},
				},
				"destination": map[string]any{
					"wallet": map[string]any{},
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
						Balance: "main",
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
			Name: "valid case with variables",
			Data: map[string]any{
				"source": map[string]any{
					"payment": map[string]any{
						"id":  "${paymentID}",
						"psp": "test",
					},
				},
				"destination": map[string]any{
					"account": map[string]any{
						"id":     "test:${var1}:test:${var2}:test:${var3}:test",
						"ledger": "default",
					},
				},
				"amount": map[string]any{
					"amount": "${amount}",
					"asset":  "EUR/2",
				},
			},
			Variables: map[string]string{
				"paymentID": "test",
				"var1":      "001",
				"var2":      "003",
				"var3":      "f5649",
				"amount":    "2819",
			},
			ExpectedResolved: Send{
				Source: Source{
					Payment: &PaymentSource{
						ID: "test",
					},
				},
				Destination: Destination{
					Account: &LedgerAccountDestination{
						ID:     "test:001:test:003:test:f5649:test",
						Ledger: "default",
					},
				},
				Amount: &shared.Monetary{
					Amount: big.NewInt(2819),
					Asset:  "EUR/2",
				},
			},
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
						WalletReference: WalletReference{
							ID: "foo",
						},
						Balance: "main",
					},
				},
				Amount: &shared.Monetary{
					Asset:  "USD",
					Amount: big.NewInt(100),
				},
			},
		},
		{
			Name: "use metadata",
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
				"metadata": map[string]string{
					"foo": "bar",
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
						WalletReference: WalletReference{
							ID: "foo",
						},
						Balance: "main",
					},
				},
				Amount: &shared.Monetary{
					Asset:  "USD",
					Amount: big.NewInt(100),
				},
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			Name: "use timestamp",
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
				"timestamp": now.Format(time.RFC3339Nano),
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
						WalletReference: WalletReference{
							ID: "foo",
						},
						Balance: "main",
					},
				},
				Amount: &shared.Monetary{
					Asset:  "USD",
					Amount: big.NewInt(100),
				},
				Timestamp: &now,
			},
		},
		{
			Name: "use timestamp as variable",
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
				"timestamp": "${timestamp}",
			},
			Variables: map[string]string{
				"timestamp": now.Format(time.RFC3339Nano),
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
						WalletReference: WalletReference{
							ID: "foo",
						},
						Balance: "main",
					},
				},
				Amount: &shared.Monetary{
					Asset:  "USD",
					Amount: big.NewInt(100),
				},
				Timestamp: &now,
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
				WalletReference: WalletReference{
					ID: "wallet1",
				},
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
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: pointer.For(paymentAccountName("payment1")),
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: map[string]string{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
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
						Data: &activities.CreditWalletRequestPayload{
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
							Metadata: map[string]string{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{nil},
			},
		},
	}
	paymentToWalletByName = stagestesting.WorkflowTestCase[Send]{
		Name: "payment to wallet by name",
		Stage: Send{
			Source: NewSource().WithPayment(&PaymentSource{
				ID: "payment1",
			}),
			Destination: NewDestination().WithWallet(&WalletDestination{
				WalletReference: WalletReference{
					Name: "user:1",
				},
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
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: pointer.For(paymentAccountName("payment1")),
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: map[string]string{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.ListWalletsActivity,
				Args: []any{mock.Anything, activities.ListWalletsRequest{
					Name: "user:1",
				}},
				Returns: []any{
					&shared.ListWalletsResponse{
						Cursor: shared.ListWalletsResponseCursor{
							Data: []shared.Wallet{{
								ID:     "wallet1",
								Ledger: "default",
							}},
						},
					}, nil,
				},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "wallet1",
						Data: &activities.CreditWalletRequestPayload{
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
							Metadata: map[string]string{
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
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: pointer.For(paymentAccountName("payment1")),
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: internalLedger,
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: map[string]string{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "foo",
								Source:      "world",
							}},
							Metadata: map[string]string{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
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
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: paymentAccountName("payment1"),
								Source:      "world",
							}},
							Reference: pointer.For(paymentAccountName("payment1")),
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
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      paymentAccountName("payment1"),
							}},
							Metadata: map[string]string{
								moveToLedgerMetadata: "default",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "foo",
								Source:      "world",
							}},
							Metadata: map[string]string{
								moveFromLedgerMetadata: internalLedger,
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
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
			Metadata: map[string]string{
				"foo": "bar",
			},
		},
		MockedActivities: []stagestesting.MockedActivity{
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "bar",
								Source:      "foo",
							}},
							Metadata: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
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
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      "account1",
							}},
							Metadata: map[string]string{
								moveToLedgerMetadata: "ledger2",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger2",
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "account2",
								Source:      "world",
							}},
							Metadata: map[string]string{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
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
				WalletReference: WalletReference{
					ID: "bar",
				},
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
						ID:     "bar",
						Ledger: "default",
					},
				}, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "bar",
						Data: &activities.CreditWalletRequestPayload{
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
							Balance:  pointer.For("main"),
							Metadata: map[string]string{},
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
				WalletReference: WalletReference{
					ID: "wallet",
				},
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
						ID:     "wallet",
					},
				}, nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "ledger1",
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      "account1",
							}},
							Metadata: map[string]string{
								moveToLedgerMetadata: "ledger2",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "wallet",
						Data: &activities.CreditWalletRequestPayload{
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
							Metadata: map[string]string{
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
				PSP:         "stripe",
				Metadata:    "stripeConnectID",
				ConnectorID: nil,
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
						Metadata: map[string]any{
							"stripeConnectID": "abcd",
						},
					},
				}, nil},
			},
			{
				Activity: activities.StripeTransferActivity,
				Args: []any{
					mock.Anything, activities.StripeTransferRequest{
						Amount:            big.NewInt(100),
						Asset:             pointer.For("USD"),
						Destination:       pointer.For("abcd"),
						ConnectorID:       nil,
						WaitingValidation: pointer.For(false),
						Metadata:          map[string]string{},
					},
				},
				Returns: []any{nil},
			},
			{
				Activity: activities.CreateTransactionActivity,
				Args: []any{
					mock.Anything, activities.CreateTransactionRequest{
						Ledger: "default",
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "world",
								Source:      "foo",
							}},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
		},
	}
	walletToAccount = stagestesting.WorkflowTestCase[Send]{
		Name: "wallet to account",
		Stage: Send{
			Source: NewSource().WithWallet(&WalletSource{
				WalletReference: WalletReference{
					ID: "foo",
				},
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
						Metadata: map[string]string{
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
						Data: &activities.DebitWalletRequestPayload{
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
				WalletReference: WalletReference{
					ID: "wallet",
				},
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
						Data: &activities.DebitWalletRequestPayload{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Balances: []string{"main"},
							Metadata: map[string]string{
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
						Data: activities.PostTransaction{
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Destination: "account",
								Source:      "world",
							}},
							Metadata: map[string]string{
								moveFromLedgerMetadata: "ledger1",
							},
						},
					},
				},
				Returns: []any{&shared.V2CreateTransactionResponse{
					Data: shared.V2Transaction{},
				}, nil},
			},
		},
	}
	walletToWallet = stagestesting.WorkflowTestCase[Send]{
		Name: "wallet to wallet",
		Stage: Send{
			Source: NewSource().WithWallet(&WalletSource{
				WalletReference: WalletReference{
					ID: "foo",
				},
				Balance: "main",
			}),
			Destination: NewDestination().WithWallet(&WalletDestination{
				WalletReference: WalletReference{
					ID: "bar",
				},
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
						ID:     "foo",
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
						ID:     "bar",
					},
				}, nil},
			},
			{
				Activity: activities.CreditWalletActivity,
				Args: []any{
					mock.Anything, activities.CreditWalletRequest{
						ID: "bar",
						Data: &activities.CreditWalletRequestPayload{
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
				WalletReference: WalletReference{
					ID: "wallet1",
				},
				Balance: "main",
			}),
			Destination: NewDestination().WithWallet(&WalletDestination{
				WalletReference: WalletReference{
					ID: "wallet2",
				},
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
						ID:     "wallet1",
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
						ID:     "wallet2",
					},
				}, nil},
			},
			{
				Activity: activities.DebitWalletActivity,
				Args: []any{
					mock.Anything, activities.DebitWalletRequest{
						ID: "wallet1",
						Data: &activities.DebitWalletRequestPayload{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Metadata: map[string]string{
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
						Data: &activities.CreditWalletRequestPayload{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Balance: pointer.For("main"),
							Metadata: map[string]string{
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
				WalletReference: WalletReference{
					ID: "foo",
				},
				Balance: "main",
			}),
			Destination: NewDestination().WithPayment(&PaymentDestination{
				PSP:         "stripe",
				Metadata:    "stripeConnectID",
				ConnectorID: nil,
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
						Metadata: map[string]string{
							"stripeConnectID": "abcd",
						},
					},
				}, nil},
			},
			{
				Activity: activities.StripeTransferActivity,
				Args: []any{
					mock.Anything, activities.StripeTransferRequest{
						Amount:            big.NewInt(100),
						Asset:             pointer.For("USD"),
						Destination:       pointer.For("abcd"),
						ConnectorID:       nil,
						WaitingValidation: pointer.For(false),
						Metadata:          map[string]string{},
					},
				},
				Returns: []any{nil},
			},
			{
				Activity: activities.DebitWalletActivity,
				Args: []any{
					mock.Anything, activities.DebitWalletRequest{
						ID: "foo",
						Data: &activities.DebitWalletRequestPayload{
							Amount: shared.Monetary{
								Asset:  "USD",
								Amount: big.NewInt(100),
							},
							Balances: []string{"main"},
							Metadata: map[string]string{},
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
	paymentToWalletByName,
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
