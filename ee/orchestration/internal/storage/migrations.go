package storage

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
)

func Migrate(ctx context.Context, db *bun.DB) error {
	migrator := migrations.NewMigrator()
	registerMigrations(migrator)

	return migrator.Up(ctx, db)
}

func registerMigrations(migrator *migrations.Migrator) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				if _, err := tx.Exec(`
					create table "workflows" (
						config jsonb,
						id varchar not null,
						created_at timestamp default now(),
						updated_at timestamp default now(),
						primary key (id)
					);
					create table "workflow_instances" (
						workflow_id varchar references workflows (id),
						id varchar,
						created_at timestamp default now(),
						updated_at timestamp default now(),
						primary key (id)
					);
					create table "workflow_instance_stage_statuses" (
						instance_id varchar references workflow_instances (id),
						stage int,
						started_at timestamp default now(),
						terminated_at timestamp default null,
						error varchar,
						primary key (instance_id, stage)
					);
				`); err != nil {
					return err
				}
				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				if _, err := tx.Exec(`
					alter table "workflow_instances" add column terminated bool;
					alter table "workflow_instances" add column terminated_at timestamp default null;
				`); err != nil {
					return err
				}
				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				if _, err := tx.Exec(`
					alter table "workflow_instances" add column error varchar;
				`); err != nil {
					return err
				}
				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				if _, err := tx.Exec(`
					alter table "workflows" add column if not exists deleted_at timestamp default null;
				`); err != nil {
					return err
				}
				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				if _, err := tx.Exec(`
					create table triggers (
					    id varchar primary key,
					    workflow_id varchar references workflows(id),
					    filter varchar null,
					    event varchar not null,
					    vars jsonb,
					    created_at timestamp not null default now(),
					    deleted_at timestamp default null
					);
					create table triggers_occurrences (
					    workflow_instance_id varchar references workflow_instances(id),
					    trigger_id varchar references triggers(id),
					    event_id varchar not null,
					    date timestamp not null default now(),
					    event jsonb not null,
					    primary key (trigger_id, event_id)
					);
				`); err != nil {
					return err
				}
				return nil
			},
		},
	)
}
