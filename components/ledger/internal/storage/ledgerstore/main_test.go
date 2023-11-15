package ledgerstore

import (
	"context"
	"github.com/uptrace/bun"
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

type T interface {
	require.TestingT
	Helper()
	Cleanup(func())
}

func newBucket(t T, hooks ...bun.QueryHook) *Bucket {
	name := uuid.NewString()
	ctx := logging.TestingContext()

	pgDatabase := pgtesting.NewPostgresDatabase(t)

	//li := pq.NewListener(pgDatabase.ConnString(), 10*time.Second, time.Minute, func(event pq.ListenerEventType, err error) {
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//})
	//
	//if err := li.Listen("raise_notice"); err != nil {
	//	panic(err)
	//}
	//
	//go func() {
	//	for {
	//		select {
	//		case n := <-li.Notify:
	//			// n.Extra contains the payload from the notification
	//			fmt.Println("notification:", n.Extra)
	//		case <-time.After(5 * time.Minute):
	//			if err := li.Ping(); err != nil {
	//				panic(err)
	//			}
	//		}
	//	}
	//}()
	//
	//defer func() {
	//	<-time.After(time.Second)
	//}()
	//
	//<-time.After(time.Second)

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

	bucket, err := ConnectToBucket(ss, connectionOptions, name, hooks...)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = bucket.Close()
	})

	require.NoError(t, bucket.Migrate(ctx))

	return bucket
}

func newLedgerStore(t *testing.T) *Store {
	t.Helper()

	ledgerName := uuid.NewString()
	ctx := logging.TestingContext()

	bucket := newBucket(t)
	store, err := bucket.GetLedgerStore(ctx, ledgerName)
	require.NoError(t, err)

	return store
}

func appendLog(t *testing.T, store *Store, log *ledger.ChainedLog) *ledger.ChainedLog {
	err := store.InsertLogs(context.Background(), log)
	require.NoError(t, err)
	return log
}
