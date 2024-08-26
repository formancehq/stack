package postgres_test

import (
	"context"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"
	"github.com/uptrace/bun"

	"github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}

	pgDB := srv.NewDatabase()
	db, err := bunconnect.OpenSQLDB(logging.TestingContext(), bunconnect.ConnectionOptions{
		DatabaseSourceName: pgDB.ConnString(),
	}, hooks...)
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	require.NoError(t, db.Ping())
	require.NoError(t, storage.Migrate(context.Background(), db))

	// Cleanup tables
	require.NoError(t, db.ResetModel(context.TODO(), (*webhooks.Config)(nil)))
	require.NoError(t, db.ResetModel(context.TODO(), (*webhooks.Attempt)(nil)))

	store, err := postgres.NewStore(db)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = store.Close(context.Background())
	})

	cfgs, err := store.FindManyConfigs(context.Background(), map[string]any{})
	require.NoError(t, err)
	require.Equal(t, 0, len(cfgs))

	ids, err := store.FindWebhookIDsToRetry(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(ids))

	atts, err := store.FindAttemptsToRetryByWebhookID(context.Background(), "")
	require.NoError(t, err)
	require.Equal(t, 0, len(atts))
}
