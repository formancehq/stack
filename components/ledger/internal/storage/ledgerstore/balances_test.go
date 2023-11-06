package ledgerstore_test

import (
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"math/big"
	"testing"
	"time"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	internaltesting "github.com/formancehq/ledger/internal/testing"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/stretchr/testify/require"
)

func TestGetBalancesAggregated(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := ledger.Now()
	ctx := logging.TestingContext()

	tx1 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("world", "users:1", "USD", big.NewInt(1)),
		ledger.NewPosting("world", "users:2", "USD", big.NewInt(199)),
	).WithDate(now)

	tx2 := ledger.NewTransaction().WithPostings(
		ledger.NewPosting("world", "users:1", "USD", big.NewInt(1)),
		ledger.NewPosting("world", "users:2", "USD", big.NewInt(199)),
	).WithDate(now.Add(time.Minute)).WithIDUint64(1)

	require.NoError(t, store.InsertLogs(ctx,
		ledger.ChainLogs(
			ledger.NewTransactionLog(tx1, map[string]metadata.Metadata{}).WithDate(tx1.Timestamp),
			ledger.NewTransactionLog(tx2, map[string]metadata.Metadata{}).WithDate(tx2.Timestamp),
			ledger.NewSetMetadataLog(now.Add(time.Minute), ledger.SetMetadataLogPayload{
				TargetType: ledger.MetaTargetTypeAccount,
				TargetID:   "users:1",
				Metadata: metadata.Metadata{
					"category": "premium",
				},
			}),
			ledger.NewSetMetadataLog(now.Add(time.Minute), ledger.SetMetadataLogPayload{
				TargetType: ledger.MetaTargetTypeAccount,
				TargetID:   "users:2",
				Metadata: metadata.Metadata{
					"category": "premium",
				},
			}),
			ledger.NewDeleteMetadataLog(now.Add(2*time.Minute), ledger.DeleteMetadataLogPayload{
				TargetType: ledger.MetaTargetTypeAccount,
				TargetID:   "users:2",
				Key:        "category",
			}),
		)...))

	t.Run("aggregate on all", func(t *testing.T) {
		t.Parallel()
		q := ledgerstore.NewPaginatedQueryOptions(ledgerstore.PITFilter{}).WithPageSize(10)
		cursor, err := store.GetAggregatedBalances(ctx, ledgerstore.NewGetAggregatedBalancesQuery(q))
		require.NoError(t, err)
		internaltesting.RequireEqual(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0),
		}, cursor)
	})
	t.Run("filter on address", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgerstore.NewGetAggregatedBalancesQuery(ledgerstore.NewPaginatedQueryOptions(ledgerstore.PITFilter{}).
			WithQueryBuilder(query.Match("address", "users:")).
			WithPageSize(10),
		))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(400),
		}, ret)
	})
	t.Run("using pit", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgerstore.NewGetAggregatedBalancesQuery(ledgerstore.NewPaginatedQueryOptions(ledgerstore.PITFilter{
			PIT: &now,
		}).
			WithQueryBuilder(query.Match("address", "users:")).
			WithPageSize(10)))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(200),
		}, ret)
	})
	t.Run("using a metadata and pit", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgerstore.NewGetAggregatedBalancesQuery(ledgerstore.NewPaginatedQueryOptions(ledgerstore.PITFilter{
			PIT: pointer.For(now.Add(time.Minute)),
		}).
			WithQueryBuilder(query.Match("metadata[category]", "premium")).
			WithPageSize(10)))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(400),
		}, ret)
	})
	t.Run("using a metadata without pit", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgerstore.NewGetAggregatedBalancesQuery(ledgerstore.NewPaginatedQueryOptions(ledgerstore.PITFilter{}).
			WithQueryBuilder(query.Match("metadata[category]", "premium")).
			WithPageSize(10)))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(2),
		}, ret)
	})
	t.Run("when no matching", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgerstore.NewGetAggregatedBalancesQuery(ledgerstore.NewPaginatedQueryOptions(ledgerstore.PITFilter{}).
			WithQueryBuilder(query.Match("metadata[category]", "guest")).
			WithPageSize(10)))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{}, ret)
	})
}
