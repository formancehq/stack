package ledgerstore

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (store *Store) getMigrator() *migrations.Migrator {
	migrator := migrations.NewMigrator(migrations.WithSchema(store.Name(), true))
	registerMigrations(migrator, store.name)
	return migrator
}

func (store *Store) Migrate(ctx context.Context) (bool, error) {
	migrator := store.getMigrator()

	if err := migrator.Up(ctx, store.db); err != nil {
		return false, err
	}

	// TODO: Update migrations package to return modifications
	return false, nil
}

func (store *Store) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	return store.getMigrator().GetMigrations(ctx, store.db)
}

//go:embed migrations/0-init-schema.sql
var initSchema string

func registerMigrations(migrator *migrations.Migrator, name string) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			Up: func(tx bun.Tx) error {

				v1SchemaExists := false
				row := tx.QueryRow(`select exists (
					select from pg_tables
					where schemaname = ? and tablename  = 'log'
				)`, name)
				if row.Err() != nil {
					return row.Err()
				}
				var ret string
				if err := row.Scan(&ret); err != nil {
					panic(err)
				}
				v1SchemaExists = ret != "false"

				if v1SchemaExists {
					_, err := tx.Exec(`alter schema rename ? to ?`, name, fmt.Sprintf(name+oldSchemaRenameSuffix))
					if err != nil {
						return errors.Wrap(err, "renaming old schema")
					}
					_, err = tx.Exec(`create schema if not exists ?`, name)
					if err != nil {
						return errors.Wrap(err, "creating new schema")
					}
				}

				_, err := tx.Exec(initSchema)
				if err != nil {
					return errors.Wrap(err, "initializing new schema")
				}

				if v1SchemaExists {
					if err := migrateLogs(context.Background(), fmt.Sprintf(name+oldSchemaRenameSuffix), name, tx); err != nil {
						return errors.Wrap(err, "migrating logs")
					}
				}

				return nil
			},
		},
	)
}
