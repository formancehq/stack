package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddErrorOnInstancesSql, downAddErrorOnInstancesSql)
}

func upAddErrorOnInstancesSql(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		alter table "workflow_instances" add column error varchar;
	`); err != nil {
		return err
	}
	return nil
}

func downAddErrorOnInstancesSql(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		alter table "workflow_instances" drop column error;
	`); err != nil {
		return err
	}
	return nil
}
