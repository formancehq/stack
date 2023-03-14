package ledger_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/formancehq/ledger/pkg/cache"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger"
	"github.com/stretchr/testify/require"
)

func TestAccountMetadata(t *testing.T) {
	runOnLedger(t, func(l *ledger.Ledger) {

		logs, err := l.SaveMeta(context.Background(), core.MetaTargetTypeAccount, "users:001", core.Metadata{
			"a random metadata": "old value",
		})
		require.NoError(t, err)
		require.NoError(t, logs.Wait(context.Background()))

		logs, err = l.SaveMeta(context.Background(), core.MetaTargetTypeAccount, "users:001", core.Metadata{
			"a random metadata": "new value",
		})
		require.NoError(t, err)
		require.NoError(t, logs.Wait(context.Background()))

		acc, err := l.GetDBCache().GetAccountWithVolumes(context.Background(), "users:001")
		require.NoError(t, err)

		meta, ok := acc.Metadata["a random metadata"]
		require.True(t, ok)

		require.Equalf(t, meta, "new value",
			"metadata entry did not match in get: expected \"new value\", got %v", meta)

		// We have to create at least one transaction to retrieve an account from GetAccounts store method
		_, logs, err = l.CreateTransaction(context.Background(), false, core.TxToScriptData(core.TransactionData{
			Postings: core.Postings{
				{
					Source:      "world",
					Amount:      core.NewMonetaryInt(100),
					Asset:       "USD",
					Destination: "users:001",
				},
			},
		}))
		require.NoError(t, err)
		require.NoError(t, logs.Wait(context.Background()))

		acc, err = l.GetDBCache().GetAccountWithVolumes(context.Background(), "users:001")
		require.NoError(t, err)
		require.NotNil(t, acc)

		meta, ok = acc.Metadata["a random metadata"]
		require.True(t, ok)
		require.Equalf(t, meta, "new value",
			"metadata entry did not match in find: expected \"new value\", got %v", meta)
	})
}

func TestTransactionMetadata(t *testing.T) {
	runOnLedger(t, func(l *ledger.Ledger) {
		logs, err := l.SaveMeta(context.Background(), core.MetaTargetTypeTransaction, uint64(0), core.Metadata{
			"a random metadata": "old value",
		})
		require.NoError(t, err)
		require.NoError(t, logs.Wait(context.Background()))
	})
}

func TestRevertTransaction(t *testing.T) {
	runOnLedger(t, func(l *ledger.Ledger) {
		tx := core.Transaction{
			TransactionData: core.TransactionData{
				Reference: "foo",
				Postings: []core.Posting{
					core.NewPosting("world", "payments:001", "COIN", core.NewMonetaryInt(100)),
				},
			},
		}
		expandedTx := core.ExpandedTransaction{
			Transaction: tx,
			PreCommitVolumes: map[string]core.AssetsVolumes{
				"world": {
					"COIN": core.NewEmptyVolumes().WithOutput(core.NewMonetaryInt(10)),
				},
				"payments:001": {
					"COIN": core.NewEmptyVolumes(),
				},
			},
			PostCommitVolumes: map[string]core.AssetsVolumes{
				"world": {
					"COIN": core.NewEmptyVolumes().WithOutput(core.NewMonetaryInt(110)),
				},
				"payments:001": {
					"COIN": core.NewEmptyVolumes().WithInput(core.NewMonetaryInt(100)),
				},
			},
		}

		require.NoError(t, l.GetLedgerStore().InsertTransactions(context.Background(), expandedTx))
		require.NoError(t, l.GetLedgerStore().EnsureAccountExists(context.Background(), "payments:001"))
		require.NoError(t, l.GetLedgerStore().UpdateVolumes(context.Background(), core.AccountsAssetsVolumes{
			"payments:001": {
				"COIN": core.NewEmptyVolumes().WithInput(core.NewMonetaryInt(110)),
			},
			"world": {
				"COIN": core.NewEmptyVolumes().WithOutput(core.NewMonetaryInt(110)),
			},
		}))

		revertTx, logs, err := l.RevertTransaction(context.Background(), tx.ID)
		require.NoError(t, err)
		require.NoError(t, logs.Wait(context.Background()))

		require.Equal(t, core.Postings{
			{
				Source:      "payments:001",
				Destination: "world",
				Amount:      core.NewMonetaryInt(100),
				Asset:       "COIN",
			},
		}, revertTx.TransactionData.Postings)

		require.EqualValues(t, fmt.Sprintf("%d", tx.ID), revertTx.Metadata[core.RevertMetadataSpecKey()])

		lastTXInfo, err := l.GetDBCache().GetLastTransaction(context.Background())
		require.NoError(t, err)
		require.NotNil(t, lastTXInfo)
		require.Equal(t, cache.TxInfo{
			Date: revertTx.Timestamp,
			ID:   tx.ID + 1,
		}, *lastTXInfo)

		account, err := l.GetDBCache().GetAccountWithVolumes(context.Background(), "payments:001")
		require.NoError(t, err)
		require.Equal(t, core.AccountWithVolumes{
			Account: core.Account{
				Address:  "payments:001",
				Metadata: core.Metadata{},
			},
			Volumes: core.AssetsVolumes{
				"COIN": core.NewEmptyVolumes().
					WithInput(core.NewMonetaryInt(110)).
					WithOutput(tx.Postings[0].Amount),
			},
			Balances: map[string]*core.MonetaryInt{
				"COIN": core.NewMonetaryInt(10),
			},
		}, *account)

		rawLogs, err := l.GetLedgerStore().ReadLogsStartingFromID(context.Background(), 0)
		require.NoError(t, err)
		require.Len(t, rawLogs, 1)
		require.Equal(t, core.NewRevertedTransactionLog(revertTx.Timestamp, tx.ID, revertTx.Transaction).
			ComputeHash(nil), rawLogs[0])
	})
}

func TestVeryBigTransaction(t *testing.T) {
	runOnLedger(t, func(l *ledger.Ledger) {
		amount, err := core.ParseMonetaryInt(
			"199999999999999999992919191919192929292939847477171818284637291884661818183647392936472918836161728274766266161728493736383838")
		require.NoError(t, err)

		_, _, err = l.CreateTransactionAndWait(context.Background(), false,
			core.TxToScriptData(core.TransactionData{
				Postings: []core.Posting{{
					Source:      "world",
					Destination: "bank",
					Asset:       "ETH/18",
					Amount:      amount,
				}},
			}))
		require.NoError(t, err)
	})
}
