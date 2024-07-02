package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWalletsPatch(t *testing.T) {
	t.Parallel()

	patchWalletRequest := wallet.PatchRequest{
		Metadata: metadata.Metadata{
			"role": "admin",
			"foo":  "baz",
		},
	}
	w := wallet.NewWallet(uuid.NewString(), "default", metadata.Metadata{
		"foo": "bar",
	})

	req := newRequest(t, http.MethodPatch, "/wallets/"+w.ID, patchWalletRequest)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetMainBalanceAccount(w.ID), account)

			return &wallet.AccountWithVolumesAndBalances{
				Account: wallet.Account{
					Address:  account,
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(w.LedgerMetadata()),
				},
			}, nil
		}),
		WithAddMetadataToAccount(func(ctx context.Context, ledger, account, ik string, md map[string]string) error {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetMainBalanceAccount(w.ID), account)
			compareJSON(t, metadata.Metadata{
				wallet.MetadataKeyWalletID:       w.ID,
				wallet.MetadataKeyWalletName:     w.Name,
				wallet.MetadataKeyWalletSpecType: wallet.PrimaryWallet,
				wallet.MetadataKeyBalanceName:    wallet.MainBalance,
				wallet.MetadataKeyWalletBalance:  wallet.TrueValue,
				wallet.MetadataKeyCreatedAt:      w.CreatedAt.UTC().Format(time.RFC3339Nano),
			}.Merge(wallet.EncodeCustomMetadata(map[string]string{
				"role": "admin",
				"foo":  "baz",
			})), md)
			return nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}
