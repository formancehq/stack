package migrations

import (
	"database/sql"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	require.NoError(t, pgtesting.CreatePostgresServer())
	defer func() {
		require.NoError(t, pgtesting.DestroyPostgresServer())
	}()

	migrator := NewMigrator()
	migrator.RegisterMigrations(
		Migration{
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`CREATE TABLE "foo" (id varchar)`)
				return err
			},
		},
	)

	db := pgtesting.NewPostgresDatabase(t)
	sqlDB, err := sql.Open("pgx", db.ConnString())
	require.NoError(t, err)

	require.NoError(t, migrator.Up(sqlDB))
	version, err := migrator.GetDBVersion(sqlDB)
	require.NoError(t, err)
	require.EqualValues(t, 1, version)
}
