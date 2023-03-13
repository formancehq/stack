package api

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	balance := wallet.NewBalance(uuid.NewString(), nil)
	assets := map[string]*big.Int{
		"USD": big.NewInt(50),
	}

	req := newRequest(t, http.MethodGet, "/wallets/"+walletID+"/balances/"+balance.Name, nil)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetBalanceAccount(walletID, balance.Name), account)

			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address:  account,
					Metadata: balance.LedgerMetadata(walletID),
				},
				Balances: assets,
			}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

	ret := wallet.ExpandedBalance{}
	readResponse(t, rec, &ret)
	require.EqualValues(t, wallet.ExpandedBalance{
		Balance: balance,
		Assets:  assets,
	}, ret)
}

func TestGetBalanceNotFound(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()

	req := newRequest(t, http.MethodGet, "/wallets/"+walletID+"/balances/xxx", nil)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address: account,
				},
			}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}
