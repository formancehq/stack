package systemstore

import (
	"context"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
)

func (s *Store) getMigrator() *migrations.Migrator {
	migrator := migrations.NewMigrator(migrations.WithSchema(Schema, true))
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.NewCreateTable().
					Model((*Ledger)(nil)).
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return sqlutils.PostgresError(err)
				}

				_, err = tx.NewCreateTable().
					Model((*Buckets)(nil)).
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return sqlutils.PostgresError(err)
				}

				_, err = s.db.NewCreateTable().
					Model((*configuration)(nil)).
					IfNotExists().
					Exec(ctx)
				return sqlutils.PostgresError(err)
			},
		},
	)
	return migrator
}
