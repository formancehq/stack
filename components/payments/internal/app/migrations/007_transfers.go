package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	up := func(tx *sql.Tx) error {
		_, err := tx.Exec(`
				CREATE TYPE transfer_status AS ENUM ('PENDING', 'SUCCEEDED', 'FAILED');

				CREATE TABLE payments.transfers (
					id uuid  NOT NULL DEFAULT gen_random_uuid(),
					connector_id uuid  NOT NULL,
					payment_id uuid NULL,
					reference text UNIQUE,
					created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
					amount bigint NOT NULL DEFAULT 0,
					currency text  NOT NULL,
					source text NOT NULL,
					destination text  NOT NULL,
					status transfer_status  NOT NULL DEFAULT 'PENDING',
					error text NULL,
					CONSTRAINT transfer_pk PRIMARY KEY (id)
				);

				ALTER TABLE payments.transfers ADD CONSTRAINT transfer_connector
					FOREIGN KEY (connector_id)
					REFERENCES connectors.connector (id)
					ON DELETE CASCADE
					NOT DEFERRABLE
					INITIALLY IMMEDIATE
				;

				ALTER TABLE payments.transfers ADD CONSTRAINT transfer_payment
					FOREIGN KEY (payment_id)
					REFERENCES payments.payment (id)
					ON DELETE CASCADE
					NOT DEFERRABLE
					INITIALLY IMMEDIATE
				;
		`)
		if err != nil {
			return err
		}

		return nil
	}

	down := func(tx *sql.Tx) error {
		_, err := tx.Exec(`
				DROP TABLE payments.transfers;
				DROP TYPE transfer_status;
		`)
		if err != nil {
			return err
		}

		return nil
	}

	goose.AddMigration(up, down)
}
