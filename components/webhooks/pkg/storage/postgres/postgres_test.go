package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestStore(t *testing.T) {
	pgDB := pgtesting.NewPostgresDatabase(t)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(pgDB.ConnString())))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer func() {
		_ = db.Close()
	}()

	require.NoError(t, db.Ping())

	// Cleanup tables
	require.NoError(t, db.ResetModel(context.TODO(), (*webhooks.Config)(nil)))
	require.NoError(t, db.ResetModel(context.TODO(), (*webhooks.Attempt)(nil)))

	store, err := postgres.NewStore(pgDB.ConnString())
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
