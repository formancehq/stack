package ledger_test

import (
	"context"
	"os"
	"testing"

	"github.com/formancehq/ledger/pkg/cache"
	"github.com/formancehq/ledger/pkg/ledger"
	"github.com/formancehq/ledger/pkg/ledgertesting"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.M) {
	if err := pgtesting.CreatePostgresServer(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	code := t.Run()
	if err := pgtesting.DestroyPostgresServer(); err != nil {
		logging.Error(err)
	}
	os.Exit(code)
}

func runOnLedger(t interface {
	pgtesting.TestingT
	Parallel()
}, f func(l *ledger.Ledger), ledgerOptions ...ledger.Option) {

	t.Parallel()

	storageDriver, close, err := ledgertesting.StorageDriver(t)
	require.NoError(t, err)
	defer close()

	require.NoError(t, storageDriver.Initialize(context.Background()))

	name := uuid.New()
	store, _, err := storageDriver.GetLedgerStore(context.Background(), name, true)
	require.NoError(t, err)

	_, err = store.Initialize(context.Background())
	require.NoError(t, err)

	dbCache := cache.NewCache(storageDriver)

	// 100 000 000 is 100MB
	l, err := ledger.NewLedger(store,
		dbCache.ForLedger(name),
		ledgerOptions...)
	require.NoError(t, err)
	defer l.Close(context.Background())

	f(l)
}
