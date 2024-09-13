package storage

import (
	"context"
	_ "embed"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
)

// EncryptionKey is set from the migration utility to specify default encryption key to migrate to.
// This can remain empty. Then the config will be removed.
//
//nolint:gochecknoglobals // This is a global variable by design.
var EncryptionKey string

//go:embed migrations/0-init-schema.sql
var initSchema string

func registerMigrations(migrator *migrations.Migrator) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "init schema",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, initSchema)
				return err
			},
		},
	)
}

func getMigrator() *migrations.Migrator {
	migrator := migrations.NewMigrator()
	registerMigrations(migrator)
	return migrator
}

func Migrate(ctx context.Context, db bun.IDB) error {
	return getMigrator().Up(ctx, db)
}
