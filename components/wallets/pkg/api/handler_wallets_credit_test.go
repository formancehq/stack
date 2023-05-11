package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWalletsCredit(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                    string
		request                 wallet.CreditRequest
		postTransactionResult   sdk.TransactionResponse
		expectedPostTransaction func(testEnv *testEnv, walletID string) sdk.PostTransaction
		expectedStatusCode      int
		expectedErrorCode       string
	}
	testCases := []testCase{
		{
			name: "nominal",
			request: wallet.CreditRequest{
				Amount: wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
				Metadata: metadata.Metadata{
					"foo": "bar",
				},
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) sdk.PostTransaction {
				return sdk.PostTransaction{
					Script: &sdk.PostTransactionScript{
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
				Amount: wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
				Sources: []wallet.Subject{
					wallet.NewLedgerAccountSubject("emitter1"),
					wallet.NewWalletSubject("wallet1", ""),
				},
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) sdk.PostTransaction {
				return sdk.PostTransaction{
					Script: &sdk.PostTransactionScript{
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
				Amount: wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
				Sources: []wallet.Subject{
					wallet.NewWalletSubject("emitter1", "secondary"),
				},
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) sdk.PostTransaction {
				return sdk.PostTransaction{
					Script: &sdk.PostTransactionScript{
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
				Amount:  wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
				Balance: "secondary",
			},
			expectedPostTransaction: func(testEnv *testEnv, walletID string) sdk.PostTransaction {
				return sdk.PostTransaction{
					Script: &sdk.PostTransactionScript{
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
				Amount:  wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
				Balance: "not-existing",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrorCodeValidation,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			walletID := uuid.NewString()
			secondaryBalance := wallet.NewBalance("secondary")

			req := newRequest(t, http.MethodPost, "/wallets/"+walletID+"/credit", testCase.request)
			rec := httptest.NewRecorder()

			var (
				testEnv         *testEnv
				postTransaction sdk.PostTransaction
			)
			testEnv = newTestEnv(
				WithCreateTransaction(func(ctx context.Context, ledger string, p sdk.PostTransaction) (*sdk.TransactionResponse, error) {
					require.Equal(t, testEnv.LedgerName(), ledger)
					postTransaction = p
					return &testCase.postTransactionResult, nil
				}),
				WithGetAccount(func(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
					if testEnv.Chart().GetBalanceAccount(walletID, secondaryBalance.Name) == account {
						return &sdk.AccountWithVolumesAndBalances{
							Address:  account,
							Metadata: secondaryBalance.LedgerMetadata(walletID),
						}, nil
					}
					return &sdk.AccountWithVolumesAndBalances{}, nil
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
