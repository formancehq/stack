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

func TestHoldsVoid(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := wallet.NewDebitHold(walletID, wallet.NewLedgerAccountSubject("bank"), "USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodPost, "/holds/"+hold.ID+"/void", nil)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)

			return &sdk.AccountWithVolumesAndBalances{
				Address:  testEnv.Chart().GetHoldAccount(hold.ID),
				Metadata: hold.LedgerMetadata(testEnv.Chart()),
				Balances: map[string]int64{
					"USD": 100,
				},
				Volumes: map[string]map[string]int64{
					"USD": {
						"input": 100,
					},
				},
			}, nil
		}),
		WithCreateTransaction(func(ctx context.Context, name string, script sdk.PostTransaction) (*sdk.TransactionResponse, error) {
			require.Equal(t, sdk.PostTransaction{
				Script: &sdk.PostTransactionScript{
					Plain: wallet.BuildCancelHoldScript("USD"),
					Vars: map[string]interface{}{
						"hold": testEnv.Chart().GetHoldAccount(hold.ID),
					},
				},
				Metadata: wallet.TransactionMetadata(nil),
			}, script)
			return &sdk.TransactionResponse{}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}
