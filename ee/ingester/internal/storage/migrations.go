package storage

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
)

func NewMigrator() *migrations.Migrator {
	migrator := migrations.NewMigrator()
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "init-database",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					create table connectors (
					    id varchar,
					    driver varchar,
					    config varchar,
					    created_at timestamp,
					    
					    primary key(id)   
					);

					create table pipelines (
					    id varchar,
					    module varchar,
					    connector_id varchar references connectors (id),
					    created_at timestamp,
					    state jsonb,
					    disabled bool,
					    
					    primary key(id)
					);
					create unique index on pipelines (module, connector_id);
				`)
				return err
			},
		},
	)
	return migrator
}
