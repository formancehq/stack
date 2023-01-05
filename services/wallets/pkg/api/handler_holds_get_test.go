package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/wallets/pkg/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHoldsGet(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := core.NewDebitHold(walletID, "bank", "USD", "", metadata.Metadata{})

	req := newRequest(t, http.MethodGet, "/holds/"+hold.ID, nil)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)
			balances := map[string]int32{
				"USD": 50,
			}
			volumes := map[string]map[string]int32{
				"USD": {
					"input": 100,
				},
			}
			return &sdk.AccountWithVolumesAndBalances{
				Address:  testEnv.Chart().GetHoldAccount(hold.ID),
				Metadata: hold.LedgerMetadata(testEnv.Chart()),
				Balances: &balances,
				Volumes:  &volumes,
			}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

	ret := core.ExpandedDebitHold{}
	readResponse(t, rec, &ret)
	require.EqualValues(t, core.ExpandedDebitHold{
		DebitHold:      hold,
		OriginalAmount: *core.NewMonetaryInt(100),
		Remaining:      *core.NewMonetaryInt(50),
	}, ret)
}
