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

func TestWalletsPatch(t *testing.T) {
	t.Parallel()

	wallet := core.NewWallet(uuid.NewString(), metadata.Metadata{
		"foo": "bar",
	})
	patchWalletRequest := PatchWalletRequest{
		Metadata: map[string]interface{}{
			"role": "admin",
			"foo":  "baz",
		},
	}

	req := newRequest(t, http.MethodPatch, "/wallets/"+wallet.ID, patchWalletRequest)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetMainAccount(wallet.ID), account)
			return &sdk.AccountWithVolumesAndBalances{
				Address:  account,
				Metadata: wallet.LedgerMetadata(),
			}, nil
		}),
		WithAddMetadataToAccount(func(ctx context.Context, ledger, account string, md metadata.Metadata) error {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetMainAccount(wallet.ID), account)
			require.EqualValues(t, metadata.Metadata{
				core.MetadataKeyWalletID:       wallet.ID,
				core.MetadataKeyWalletName:     wallet.Name,
				core.MetadataKeyWalletSpecType: core.PrimaryWallet,
				core.MetadataKeyWalletCustomData: metadata.Metadata{
					"role": "admin",
					"foo":  "baz",
				},
			}, md)
			return nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}
