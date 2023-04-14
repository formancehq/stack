package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name                    string
	request                 wallet.DebitRequest
	postTransactionError    *apiErrorMock
	expectedPostTransaction func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction
	expectedStatusCode      int
	expectedErrorCode       string
}

type apiErrorMock struct {
	ErrorCode    sdk.ErrorsEnum `json:"errorCode,omitempty"`
	ErrorMessage string         `json:"errorMessage,omitempty"`
	Details      *string        `json:"details,omitempty"`
}

func (a *apiErrorMock) Model() any {
	if a == nil {
		return nil
	}
	return sdk.ErrorResponse{
		ErrorCode:    a.ErrorCode,
		ErrorMessage: a.ErrorMessage,
		Details:      a.Details,
	}
}

func (a *apiErrorMock) Error() string {
	if a == nil {
		return ""
	}
	by, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(by)
}

var walletDebitTestCases = []testCase{
	{
		name: "nominal",
		request: wallet.DebitRequest{
			Amount: wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": wallet.DefaultDebitDest.Identifier,
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
		name: "with custom destination as ledger account",
		request: wallet.DebitRequest{
			Amount:      wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
			Destination: wallet.Ptr(wallet.NewLedgerAccountSubject("account1")),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": "account1",
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
		name: "with custom destination as wallet",
		request: wallet.DebitRequest{
			Amount:      wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
			Destination: wallet.Ptr(wallet.NewWalletSubject("wallet1", "")),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": testEnv.Chart().GetMainBalanceAccount("wallet1"),
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
		name: "with insufficient funds",
		request: wallet.DebitRequest{
			Amount: wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
		},
		postTransactionError: &apiErrorMock{
			ErrorCode: sdk.INSUFFICIENT_FUND,
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  string(sdk.INSUFFICIENT_FUND),
	},
	{
		name: "with debit hold",
		request: wallet.DebitRequest{
			Amount:  wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
			Pending: true,
			Metadata: map[string]string{
				"foo": "bar",
			},
			Description: "a first tx",
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": testEnv.Chart().GetHoldAccount(h.ID),
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
		expectedStatusCode: http.StatusCreated,
	},
	{
		name: "with custom balance as source",
		request: wallet.DebitRequest{
			Amount:   wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
			Balances: []string{"secondary"},
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetBalanceAccount(walletID, "secondary")),
					Vars: map[string]interface{}{
						"destination": "world",
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
		name: "with wildcard balance as source",
		request: wallet.DebitRequest{
			Amount:   wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
			Balances: []string{"*"},
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetBalanceAccount(walletID, "secondary")),
					Vars: map[string]interface{}{
						"destination": "world",
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
		name: "with wildcard plus another source",
		request: wallet.DebitRequest{
			Amount:   wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
			Balances: []string{"*", "secondary"},
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetBalanceAccount(walletID, "secondary")),
					Vars: map[string]interface{}{
						"destination": "world",
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: wallet.TransactionMetadata(nil),
			}
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  string(sdk.VALIDATION),
	},
	{
		name: "with custom balance as destination",
		request: wallet.DebitRequest{
			Amount:      wallet.NewMonetary(wallet.NewMonetaryInt(100), "USD"),
			Destination: wallet.Ptr(wallet.NewWalletSubject("wallet1", "secondary")),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) sdk.PostTransaction {
			return sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": testEnv.Chart().GetBalanceAccount("wallet1", "secondary"),
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
}

func TestWalletsDebit(t *testing.T) {
	t.Parallel()
	for _, testCase := range walletDebitTestCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			walletID := uuid.NewString()

			req := newRequest(t, http.MethodPost, "/wallets/"+walletID+"/debit", testCase.request)
			rec := httptest.NewRecorder()

			var (
				testEnv             *testEnv
				postTransaction     sdk.PostTransaction
				holdAccount         string
				holdAccountMetadata metadata.Metadata
			)
			testEnv = newTestEnv(
				WithAddMetadataToAccount(func(ctx context.Context, ledger, account string, m metadata.Metadata) error {
					require.Equal(t, testEnv.LedgerName(), ledger)
					holdAccount = account
					holdAccountMetadata = m
					return nil
				}),
				WithListAccounts(func(ctx context.Context, ledger string, query wallet.ListAccountsQuery) (*sdk.AccountsCursorResponseCursor, error) {
					require.Equal(t, testEnv.LedgerName(), ledger)
					require.Equal(t, query.Metadata, wallet.BalancesMetadataFilter(walletID))
					return &sdk.AccountsCursorResponseCursor{
						Data: []sdk.Account{{
							Address: testEnv.Chart().GetBalanceAccount(walletID, "secondary"),
							Metadata: wallet.Balance{
								Name: "secondary",
							}.LedgerMetadata(walletID),
						}},
					}, nil
				}),
				WithCreateTransaction(func(ctx context.Context, ledger string, p sdk.PostTransaction) (*sdk.TransactionResponse, error) {
					require.Equal(t, testEnv.LedgerName(), ledger)
					postTransaction = p
					if testCase.postTransactionError != nil {
						return nil, testCase.postTransactionError
					}
					//nolint:nilnil
					return nil, nil
				}),
			)
			testEnv.Router().ServeHTTP(rec, req)

			expectedStatusCode := testCase.expectedStatusCode
			if expectedStatusCode == 0 {
				expectedStatusCode = http.StatusNoContent
			}
			require.Equal(t, expectedStatusCode, rec.Result().StatusCode)

			hold := &wallet.DebitHold{}
			switch expectedStatusCode {
			case http.StatusCreated:
				readResponse(t, rec, hold)
			case http.StatusNoContent:
			default:
				errorResponse := readErrorResponse(t, rec)
				require.Equal(t, testCase.expectedErrorCode, errorResponse.ErrorCode)
				return
			}

			if testCase.expectedPostTransaction != nil {
				expectedPostTransaction := testCase.expectedPostTransaction(testEnv, walletID, hold)
				require.Equal(t, expectedPostTransaction, postTransaction)
			}

			if testCase.request.Pending {
				require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), holdAccount)
				require.Equal(t, walletID, hold.WalletID)
				require.Equal(t, testCase.request.Amount.Asset, hold.Asset)
				require.Equal(t, hold.LedgerMetadata(testEnv.Chart()), holdAccountMetadata)
			}
		})
	}
}
