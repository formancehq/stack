package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestConfigStore(t *testing.T) {
	newStore := func(t *testing.T) webhooks.ConfigStore {
		pgDB := pgtesting.NewPostgresDatabase(t)
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(pgDB.ConnString())))
		db := bun.NewDB(sqldb, pgdialect.New())
		require.NoError(t, db.Ping())
		require.NoError(t, db.ResetModel(context.TODO(), (*webhooks.Config)(nil)))
		store, err := NewStore(pgDB.ConnString())
		require.NoError(t, err)
		t.Cleanup(func() {
			store.Close(context.Background())
			db.Close()
		})
		return store
	}
	storage.TestConfigStore(t, newStore)
}
