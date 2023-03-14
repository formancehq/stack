package cache

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledgertesting"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	require.NoError(t, pgtesting.CreatePostgresServer())
	defer func() {
		require.NoError(t, pgtesting.DestroyPostgresServer())
	}()

	driver, close, err := ledgertesting.StorageDriver(t)
	require.NoError(t, err)
	defer close()

	require.NoError(t, driver.Initialize(context.Background()))

	dbCache := NewCache(driver)
	ledger := uuid.NewString()
	cache := dbCache.ForLedger(ledger)

	lastInsertedTransactionID, err := cache.GetLastTransaction(context.Background())
	require.NoError(t, err)
	require.Nil(t, lastInsertedTransactionID)

	account, err := cache.GetAccountWithVolumes(context.Background(), "world")
	require.NoError(t, err)
	require.Nil(t, account)

	ledgerStore, _, err := driver.GetLedgerStore(context.Background(), ledger, true)
	require.NoError(t, err)

	_, err = ledgerStore.Initialize(context.Background())
	require.NoError(t, err)

	require.NoError(t, ledgerStore.EnsureAccountExists(context.Background(), "world"))
	require.NoError(t, ledgerStore.UpdateVolumes(context.Background(), core.AccountsAssetsVolumes{
		"world": {
			"USD/2": {
				Input:  core.NewMonetaryInt(100),
				Output: core.NewMonetaryInt(0),
			},
		},
	}))

	account, err = cache.GetAccountWithVolumes(context.Background(), "world")
	require.NoError(t, err)
	require.NotNil(t, account)
	require.Equal(t, core.AccountWithVolumes{
		Account: core.Account{
			Address:  "world",
			Metadata: core.Metadata{},
		},
		Volumes: map[string]core.Volumes{
			"USD/2": {
				Input:  core.NewMonetaryInt(100),
				Output: core.NewMonetaryInt(0),
			},
		},
		Balances: map[string]*core.MonetaryInt{
			"USD/2": core.NewMonetaryInt(100),
		},
	}, *account)

	volumes := account.Volumes["USD/2"]
	volumes.Output = account.Volumes["USD/2"].Output.Add(core.NewMonetaryInt(10))
	account.Volumes["USD/2"] = volumes
	account.Balances["USD/2"] = core.NewMonetaryInt(90)

	errChan := ledgerStore.AppendLogs(
		context.Background(),
		core.NewTransactionLog(core.Transaction{
			TransactionData: core.TransactionData{
				Postings: []core.Posting{{
					Source:      "world",
					Destination: "bank",
					Amount:      core.NewMonetaryInt(10),
					Asset:       "USD/2",
				}},
			},
		}, nil),
		core.NewSetMetadataLog(core.Now(), core.SetMetadataLogPayload{
			TargetType: core.MetaTargetTypeAccount,
			TargetID:   "bank",
			Metadata: core.Metadata{
				"category": "gold",
			},
		}),
	)
	select {
	case err := <-errChan:
		require.NoError(t, err)
	case <-time.After(time.Second):
		require.Fail(t, "timeout waiting for log insertion")
	}

	account, err = cache.GetAccountWithVolumes(context.Background(), "bank")
	require.NoError(t, err)
	require.NotNil(t, account)

	require.Equal(t, core.AccountWithVolumes{
		Account: core.Account{
			Address: "bank",
			Metadata: core.Metadata{
				"category": "gold",
			},
		},
		Volumes: map[string]core.Volumes{
			"USD/2": {
				Input:  core.NewMonetaryInt(10),
				Output: core.NewMonetaryInt(0),
			},
		},
		Balances: map[string]*core.MonetaryInt{
			"USD/2": core.NewMonetaryInt(10),
		},
	}, *account)

}
