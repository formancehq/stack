package storage

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	migrationstorage "github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
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

func newStore(t *testing.T) *Storage {
	t.Helper()

	pgServer := pgtesting.NewPostgresDatabase(t)

	config, err := pgx.ParseConfig(pgServer.ConnString())
	require.NoError(t, err)

	key := make([]byte, 64)
	_, err = rand.Read(key)
	require.NoError(t, err)

	db := bun.NewDB(stdlib.OpenDB(*config), pgdialect.New())
	t.Cleanup(func() {
		_ = db.Close()
	})

	err = migrationstorage.Migrate(context.Background(), db)
	require.NoError(t, err)

	store := NewStorage(
		db,
		string(key),
	)

	return store
}
