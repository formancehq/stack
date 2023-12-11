package api

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHoldsGet(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := wallet.NewDebitHold(walletID, wallet.NewLedgerAccountSubject("bank"),
		"USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodGet, "/holds/"+hold.ID, nil)
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
					"USD": big.NewInt(50),
				},
				Volumes: map[string]map[string]*big.Int{
					"USD": {
						"input": big.NewInt(100),
					},
				},
			}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

	ret := wallet.ExpandedDebitHold{}
	readResponse(t, rec, &ret)
	require.EqualValues(t, wallet.ExpandedDebitHold{
		DebitHold:      hold,
		OriginalAmount: big.NewInt(100),
		Remaining:      big.NewInt(50),
	}, ret)
}
