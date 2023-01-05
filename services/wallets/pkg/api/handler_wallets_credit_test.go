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

func TestWalletsCredit(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	creditWalletRequest := CreditWalletRequest{
		Amount: core.Monetary{
			Amount: core.NewMonetaryInt(100),
			Asset:  "USD",
		},
		Metadata: map[string]interface{}{
			"foo": "bar",
		},
	}

	req := newRequest(t, http.MethodPost, "/wallets/"+walletID+"/credit", creditWalletRequest)
	rec := httptest.NewRecorder()

	var (
		ledger          string
		transactionData sdk.TransactionData
	)
	testEnv := newTestEnv(
		WithCreateTransaction(func(ctx context.Context, l string, t sdk.TransactionData) error {
			ledger = l
			transactionData = t
			return nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
	require.Equal(t, testEnv.LedgerName(), ledger)
	require.Equal(t, sdk.TransactionData{
		Postings: []sdk.Posting{{
			Amount:      100,
			Asset:       "USD",
			Destination: testEnv.Chart().GetMainAccount(walletID),
			Source:      "world",
		}},
		Metadata: core.WalletTransactionBaseMetadata().Merge(metadata.Metadata{
			core.MetadataKeyWalletCustomData: creditWalletRequest.Metadata,
		}),
	}, transactionData)
}
