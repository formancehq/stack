package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upInitSchemaSql, downInitSchemaSql)
}

func upInitSchemaSql(tx *sql.Tx) error {
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
}

func downInitSchemaSql(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		drop table "workflow_instance_stage_statuses";
		drop table "workflow_instances";
		drop table "workflows";
	`); err != nil {
		return err
	}
	return nil
}
