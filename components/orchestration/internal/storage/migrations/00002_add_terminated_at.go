package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddTerminatedAtSql, downAddTerminatedAtSql)
}

func upAddTerminatedAtSql(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		alter table "workflow_instances" add column terminated bool;
		alter table "workflow_instances" add column terminated_at timestamp default null;
	`); err != nil {
		return err
	}
	return nil
}

func downAddTerminatedAtSql(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		alter table "workflow_instances" drop column terminated;
		alter table "workflow_instances" drop column terminated_at;
	`); err != nil {
		return err
	}
	return nil
}
