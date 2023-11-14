package ledgerstore

import (
	"context"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/require"
	"testing"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

// TODO: remove that
func insertTransactions(ctx context.Context, s *Store, txs ...ledger.Transaction) error {
	var previous *ledger.ChainedLog
	logs := collectionutils.Map(txs, func(from ledger.Transaction) *ledger.ChainedLog {
		previous = ledger.NewTransactionLog(&from, map[string]metadata.Metadata{}).ChainLog(previous)
		return previous
	})
	return s.InsertLogs(ctx, logs...)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctx := logging.TestingContext()

	store := newLedgerStore(t)
	require.NoError(t, store.Delete(ctx))
}
