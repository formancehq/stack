package api

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func compareJSON(t *testing.T, expected, actual any) {
	data, err := json.Marshal(expected)
	require.NoError(t, err)

	expectedAsMap := make(map[string]any)
	require.NoError(t, json.Unmarshal(data, &expectedAsMap))

	data, err = json.Marshal(actual)
	require.NoError(t, err)

	actualAsMap := make(map[string]any)
	require.NoError(t, json.Unmarshal(data, &actualAsMap))

	require.Equal(t, expectedAsMap, actualAsMap)
}

type testCase struct {
	name                    string
	request                 wallet.DebitRequest
	postTransactionError    *apiErrorMock
	expectedPostTransaction func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction
	expectedStatusCode      int
	expectedErrorCode       string
}

type apiErrorMock struct {
	ErrorCode    shared.ErrorsEnum `json:"errorCode,omitempty"`
	ErrorMessage string            `json:"errorMessage,omitempty"`
	Details      *string           `json:"details,omitempty"`
}

func (a *apiErrorMock) Model() any {
	if a == nil {
		return nil
	}
	return shared.ErrorResponse{
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
			Amount: wallet.NewMonetary(big.NewInt(100), "USD"),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": wallet.DefaultDebitDest.Identifier,
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}
		},
	},
	{
		name: "with custom destination as ledger account",
		request: wallet.DebitRequest{
			Amount:      wallet.NewMonetary(big.NewInt(100), "USD"),
			Destination: wallet.Ptr(wallet.NewLedgerAccountSubject("account1")),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": "account1",
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}
		},
	},
	{
		name: "with custom destination as wallet",
		request: wallet.DebitRequest{
			Amount:      wallet.NewMonetary(big.NewInt(100), "USD"),
			Destination: wallet.Ptr(wallet.NewWalletSubject("wallet1", "")),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": testEnv.Chart().GetMainBalanceAccount("wallet1"),
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}
		},
	},
	{
		name: "with insufficient funds",
		request: wallet.DebitRequest{
			Amount: wallet.NewMonetary(big.NewInt(100), "USD"),
		},
		postTransactionError: &apiErrorMock{
			ErrorCode: shared.ErrorsEnumInsufficientFund,
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  string(shared.ErrorsEnumInsufficientFund),
	},
	{
		name: "with debit hold",
		request: wallet.DebitRequest{
			Amount:  wallet.NewMonetary(big.NewInt(100), "USD"),
			Pending: true,
			Metadata: map[string]string{
				"foo": "bar",
			},
			Description: "a first tx",
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": testEnv.Chart().GetHoldAccount(h.ID),
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(metadata.Metadata{
					"foo": "bar",
				})),
			}
		},
		expectedStatusCode: http.StatusCreated,
	},
	{
		name: "with custom balance as source",
		request: wallet.DebitRequest{
			Amount:   wallet.NewMonetary(big.NewInt(100), "USD"),
			Balances: []string{"secondary"},
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetBalanceAccount(walletID, "secondary")),
					Vars: map[string]interface{}{
						"destination": "world",
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}
		},
	},
	{
		name: "with wildcard balance as source",
		request: wallet.DebitRequest{
			Amount:   wallet.NewMonetary(big.NewInt(100), "USD"),
			Balances: []string{"*"},
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(
						testEnv.Chart().GetBalanceAccount(walletID, "coupon1"),
						testEnv.Chart().GetBalanceAccount(walletID, "coupon4"),
						testEnv.Chart().GetBalanceAccount(walletID, "coupon2"),
						testEnv.Chart().GetBalanceAccount(walletID, "main"),
					),
					Vars: map[string]interface{}{
						"destination": "world",
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}
		},
	},
	{
		name: "with wildcard plus another source",
		request: wallet.DebitRequest{
			Amount:   wallet.NewMonetary(big.NewInt(100), "USD"),
			Balances: []string{"*", "secondary"},
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetBalanceAccount(walletID, "secondary")),
					Vars: map[string]interface{}{
						"destination": "world",
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  string(shared.ErrorErrorCodeValidation),
	},
	{
		name: "with custom balance as destination",
		request: wallet.DebitRequest{
			Amount:      wallet.NewMonetary(big.NewInt(100), "USD"),
			Destination: wallet.Ptr(wallet.NewWalletSubject("wallet1", "secondary")),
		},
		expectedPostTransaction: func(testEnv *testEnv, walletID string, h *wallet.DebitHold) wallet.PostTransaction {
			return wallet.PostTransaction{
				Script: &shared.PostTransactionScript{
					Plain: wallet.BuildDebitWalletScript(testEnv.Chart().GetMainBalanceAccount(walletID)),
					Vars: map[string]interface{}{
						"destination": testEnv.Chart().GetBalanceAccount("wallet1", "secondary"),
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
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
				postTransaction     wallet.PostTransaction
				holdAccount         string
				holdAccountMetadata map[string]string
			)
			testEnv = newTestEnv(
				WithAddMetadataToAccount(func(ctx context.Context, ledger, account string, m map[string]string) error {
					require.Equal(t, testEnv.LedgerName(), ledger)
					holdAccount = account
					holdAccountMetadata = m
					return nil
				}),
				WithListAccounts(func(ctx context.Context, ledger string, query wallet.ListAccountsQuery) (*wallet.AccountsCursorResponseCursor, error) {
					require.Equal(t, testEnv.LedgerName(), ledger)
					require.Equal(t, query.Metadata, wallet.BalancesMetadataFilter(walletID))

					return &wallet.AccountsCursorResponseCursor{
						Data: []wallet.Account{
							{
								Address: testEnv.Chart().GetBalanceAccount(walletID, "coupon2"),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.Balance{
									Name:     "coupon2",
									Priority: 10,
								}.LedgerMetadata(walletID)),
							},
							{
								Address: testEnv.Chart().GetBalanceAccount(walletID, "coupon1"),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.Balance{
									Name:      "coupon1",
									ExpiresAt: ptr(time.Now().Add(5 * time.Second)),
								}.LedgerMetadata(walletID)),
							},
							{
								Address: testEnv.Chart().GetBalanceAccount(walletID, "coupon3"),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.Balance{
									Name:      "coupon3",
									ExpiresAt: ptr(time.Now().Add(-time.Minute)),
								}.LedgerMetadata(walletID)),
							},
							{
								Address: testEnv.Chart().GetBalanceAccount(walletID, "coupon4"),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.Balance{
									Name: "coupon4",
								}.LedgerMetadata(walletID)),
							},
							{
								Address: testEnv.Chart().GetBalanceAccount(walletID, "main"),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.Balance{
									Name: "main",
								}.LedgerMetadata(walletID)),
							},
						},
					}, nil
				}),
				WithCreateTransaction(func(ctx context.Context, ledger string, p wallet.PostTransaction) (*shared.Transaction, error) {
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
				compareJSON(t, expectedPostTransaction, postTransaction)
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
