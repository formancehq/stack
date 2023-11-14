package ledgerstore

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/formancehq/ledger/internal/storage/sqlutils"
	"github.com/formancehq/ledger/internal/storage/systemstore"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := pgtesting.CreatePostgresServer(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	code := m.Run()
	if err := pgtesting.DestroyPostgresServer(); err != nil {
		logging.Error(err)
	}
	os.Exit(code)
}

func newLedgerStore(t *testing.T) *Store {
	t.Helper()

	name := uuid.NewString()
	ctx := logging.TestingContext()

	pgDatabase := pgtesting.NewPostgresDatabase(t)

	connectionOptions := sqlutils.ConnectionOptions{
		DatabaseSourceName: pgDatabase.ConnString(),
		Debug:              testing.Verbose(),
		MaxIdleConns:       40,
		MaxOpenConns:       40,
		ConnMaxIdleTime:    time.Minute,
	}

	ss, err := systemstore.Connect(ctx, connectionOptions)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = ss.Close()
	})
	require.NoError(t, ss.Migrate(ctx))

	bucket, err := ConnectToBucket(ss, connectionOptions, name)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = bucket.Close()
	})
	require.NoError(t, bucket.Migrate(ctx))

	store, err := bucket.GetLedgerStore(ctx, name)
	require.NoError(t, err)

	return store
}

func appendLog(t *testing.T, store *Store, log *ledger.ChainedLog) *ledger.ChainedLog {
	err := store.InsertLogs(context.Background(), log)
	require.NoError(t, err)
	return log
}
