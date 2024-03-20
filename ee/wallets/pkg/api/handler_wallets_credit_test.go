package api

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWalletsCredit(t *testing.T) {
	t.Parallel()
	now := time.Now()

	type testCase struct {
		name                    string
		request                 wallet.CreditRequest
		postTransactionResult   shared.Transaction
		expectedPostTransaction func(testEnv *testEnv, walletID string) wallet.PostTransaction
		expectedStatusCode      int
		expectedErrorCode       string
	}
	testCases := []testCase{
		{
			name: "nominal",
			request: wallet.CreditRequest{
				Amount: wallet.NewMonetary(big.NewInt(100), "USD"),
				Metadata: metadata.Metadata{
					"foo": "bar",
				},
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) wallet.PostTransaction {
				return wallet.PostTransaction{
					Script: &shared.PostTransactionScript{
						Plain: wallet.BuildCreditWalletScript("world"),
						Vars: map[string]interface{}{
							"destination": testEnv.chart.GetMainBalanceAccount(walletID),
							"amount": map[string]any{
								"amount": uint64(100),
								"asset":  "USD",
							},
						},
					},
					Metadata: wallet.TransactionMetadata(metadata.Metadata{
						"foo": "bar",
					}),
				}
			},
		},
		{
			name: "with source list",
			request: wallet.CreditRequest{
				Amount: wallet.NewMonetary(big.NewInt(100), "USD"),
				Sources: []wallet.Subject{
					wallet.NewLedgerAccountSubject("emitter1"),
					wallet.NewWalletSubject("wallet1", ""),
				},
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) wallet.PostTransaction {
				return wallet.PostTransaction{
					Script: &shared.PostTransactionScript{
						Plain: wallet.BuildCreditWalletScript(
							"emitter1",
							testEnv.Chart().GetMainBalanceAccount("wallet1"),
						),
						Vars: map[string]interface{}{
							"destination": testEnv.chart.GetMainBalanceAccount(walletID),
							"amount": map[string]any{
								"amount": uint64(100),
								"asset":  "USD",
							},
						},
					},
					Metadata: wallet.TransactionMetadata(nil),
				}
			},
		},
		{
			name: "with secondary balance from source",
			request: wallet.CreditRequest{
				Amount: wallet.NewMonetary(big.NewInt(100), "USD"),
				Sources: []wallet.Subject{
					wallet.NewWalletSubject("emitter1", "secondary"),
				},
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) wallet.PostTransaction {
				return wallet.PostTransaction{
					Script: &shared.PostTransactionScript{
						Plain: wallet.BuildCreditWalletScript(
							testEnv.Chart().GetBalanceAccount("emitter1", "secondary"),
						),
						Vars: map[string]interface{}{
							"destination": testEnv.Chart().GetMainBalanceAccount(walletID),
							"amount": map[string]any{
								"amount": uint64(100),
								"asset":  "USD",
							},
						},
					},
					Metadata: wallet.TransactionMetadata(nil),
				}
			},
		},
		{
			name: "with secondary balance as destination",
			request: wallet.CreditRequest{
				Amount:  wallet.NewMonetary(big.NewInt(100), "USD"),
				Balance: "secondary",
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) wallet.PostTransaction {
				return wallet.PostTransaction{
					Script: &shared.PostTransactionScript{
						Plain: wallet.BuildCreditWalletScript("world"),
						Vars: map[string]interface{}{
							"destination": testEnv.Chart().GetBalanceAccount(walletID, "secondary"),
							"amount": map[string]any{
								"amount": uint64(100),
								"asset":  "USD",
							},
						},
					},
					Metadata: wallet.TransactionMetadata(nil),
				}
			},
		},
		{
			name: "with not existing secondary balance as destination",
			request: wallet.CreditRequest{
				Amount:  wallet.NewMonetary(big.NewInt(100), "USD"),
				Balance: "not-existing",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrorCodeValidation,
		},
		{
			name: "with specified timestamp",
			request: wallet.CreditRequest{
				Amount:    wallet.NewMonetary(big.NewInt(100), "USD"),
				Timestamp: &now,
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) wallet.PostTransaction {
				return wallet.PostTransaction{
					Script: &shared.PostTransactionScript{
						Plain: wallet.BuildCreditWalletScript("world"),
						Vars: map[string]interface{}{
							"destination": testEnv.chart.GetMainBalanceAccount(walletID),
							"amount": map[string]any{
								"amount": uint64(100),
								"asset":  "USD",
							},
						},
					},
					Metadata:  wallet.TransactionMetadata(nil),
					Timestamp: &now,
				}
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			walletID := uuid.NewString()
			secondaryBalance := wallet.NewBalance("secondary", nil)

			req := newRequest(t, http.MethodPost, "/wallets/"+walletID+"/credit", testCase.request)
			rec := httptest.NewRecorder()

			var (
				testEnv         *testEnv
				postTransaction wallet.PostTransaction
			)
			testEnv = newTestEnv(
				WithCreateTransaction(func(ctx context.Context, ledger string, p wallet.PostTransaction) (*shared.Transaction, error) {
					require.Equal(t, testEnv.LedgerName(), ledger)
					postTransaction = p
					return &testCase.postTransactionResult, nil
				}),
				WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
					if testEnv.Chart().GetBalanceAccount(walletID, secondaryBalance.Name) == account {
						return &wallet.AccountWithVolumesAndBalances{
							Account: wallet.Account{
								Address:  account,
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(secondaryBalance.LedgerMetadata(walletID)),
							},
						}, nil
					}
					return &wallet.AccountWithVolumesAndBalances{
						Account: wallet.Account{
							Metadata: map[string]string{},
						},
					}, nil
				}),
			)
			testEnv.Router().ServeHTTP(rec, req)

			expectedStatusCode := testCase.expectedStatusCode
			if expectedStatusCode == 0 {
				expectedStatusCode = http.StatusNoContent
			}

			require.Equal(t, expectedStatusCode, rec.Result().StatusCode)
			if expectedStatusCode == http.StatusNoContent {
				if testCase.expectedPostTransaction != nil {
					expectedScript := testCase.expectedPostTransaction(testEnv, walletID)
					require.Equal(t, expectedScript, postTransaction)
				}
			} else {
				errorResponse := readErrorResponse(t, rec)
				require.Equal(t, ErrorCodeValidation, errorResponse.ErrorCode)
			}
		})
	}
}
