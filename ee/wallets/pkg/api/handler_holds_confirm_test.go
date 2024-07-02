package api

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHoldsConfirm(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := wallet.NewDebitHold(walletID, wallet.NewLedgerAccountSubject("bank"), "USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodPost, "/holds/"+hold.ID+"/confirm", nil)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)
			balances := map[string]*big.Int{
				"USD": big.NewInt(100),
			}
			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address:  testEnv.Chart().GetHoldAccount(hold.ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold.LedgerMetadata(testEnv.Chart())),
				},
				Balances: balances,
			}, nil
		}),
		WithCreateTransaction(func(ctx context.Context, name, ik string, postTransaction wallet.PostTransaction) (*shared.V2Transaction, error) {
			compareJSON(t, wallet.PostTransaction{
				Script: &shared.V2PostTransactionScript{
					Plain: wallet.BuildConfirmHoldScript(false, "USD"),
					Vars: map[string]interface{}{
						"hold": testEnv.Chart().GetHoldAccount(hold.ID),
						"amount": map[string]any{
							"amount": uint64(100),
							"asset":  "USD",
						},
						"dest": "bank",
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}, postTransaction)
			return &shared.V2Transaction{}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}

func TestHoldsPartialConfirm(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := wallet.NewDebitHold(walletID, wallet.NewLedgerAccountSubject("bank"), "USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodPost, "/holds/"+hold.ID+"/confirm", ConfirmHoldRequest{
		Amount: 50,
	})
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)

			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address:  testEnv.Chart().GetHoldAccount(hold.ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold.LedgerMetadata(testEnv.Chart())),
				},
				Balances: map[string]*big.Int{
					"USD": big.NewInt(100),
				},
				Volumes: map[string]shared.V2Volume{
					"USD": {
						Input: big.NewInt(100),
					},
				},
			}, nil
		}),
		WithCreateTransaction(func(ctx context.Context, name, ik string, postTransaction wallet.PostTransaction) (*shared.V2Transaction, error) {
			compareJSON(t, wallet.PostTransaction{
				Script: &shared.V2PostTransactionScript{
					Plain: wallet.BuildConfirmHoldScript(false, "USD"),
					Vars: map[string]interface{}{
						"hold": testEnv.Chart().GetHoldAccount(hold.ID),
						"amount": map[string]any{
							"amount": uint64(50),
							"asset":  "USD",
						},
						"dest": "bank",
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}, postTransaction)
			return &shared.V2Transaction{}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}

func TestHoldsConfirmWithTooHighAmount(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := wallet.NewDebitHold(walletID, wallet.NewLedgerAccountSubject("bank"), "USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodPost, "/holds/"+hold.ID+"/confirm", ConfirmHoldRequest{
		Amount: 500,
	})
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)

			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address:  testEnv.Chart().GetHoldAccount(hold.ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold.LedgerMetadata(testEnv.Chart())),
				},
				Balances: map[string]*big.Int{
					"USD": big.NewInt(100),
				},
				Volumes: map[string]shared.V2Volume{
					"USD": {
						Input: big.NewInt(100),
					},
				},
			}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	errorResponse := readErrorResponse(t, rec)
	require.Equal(t, ErrorCodeInsufficientFund, errorResponse.ErrorCode)
}

func TestHoldsConfirmWithClosedHold(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := wallet.NewDebitHold(walletID, wallet.NewLedgerAccountSubject("bank"), "USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodPost, "/holds/"+hold.ID+"/confirm", ConfirmHoldRequest{})
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)

			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address:  testEnv.Chart().GetHoldAccount(hold.ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold.LedgerMetadata(testEnv.Chart())),
				},
				Balances: map[string]*big.Int{
					"USD": big.NewInt(0),
				},
				Volumes: map[string]shared.V2Volume{
					"USD": {
						Input: big.NewInt(100),
					},
				},
			}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	errorResponse := readErrorResponse(t, rec)
	require.Equal(t, ErrorCodeClosedHold, errorResponse.ErrorCode)
}

func TestHoldsPartialConfirmWithFinal(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := wallet.NewDebitHold(walletID, wallet.NewLedgerAccountSubject("bank"),
		"USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodPost, "/holds/"+hold.ID+"/confirm", ConfirmHoldRequest{
		Amount: 50,
		Final:  true,
	})
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)

			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address:  testEnv.Chart().GetHoldAccount(hold.ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold.LedgerMetadata(testEnv.Chart())),
				},
				Balances: map[string]*big.Int{
					"USD": big.NewInt(100),
				},
				Volumes: map[string]shared.V2Volume{
					"USD": {
						Input: big.NewInt(100),
					},
				},
			}, nil
		}),
		WithCreateTransaction(func(ctx context.Context, name, ik string, script wallet.PostTransaction) (*shared.V2Transaction, error) {
			compareJSON(t, wallet.PostTransaction{
				Script: &shared.V2PostTransactionScript{
					Plain: wallet.BuildConfirmHoldScript(true, "USD"),
					Vars: map[string]interface{}{
						"hold": testEnv.Chart().GetHoldAccount(hold.ID),
						"amount": map[string]any{
							"amount": uint64(50),
							"asset":  "USD",
						},
						"dest":             "bank",
						"void_destination": testEnv.Chart().GetMainBalanceAccount(hold.WalletID),
					},
				},
				Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallet.TransactionMetadata(nil)),
			}, script)
			return &shared.V2Transaction{}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}
