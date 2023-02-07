package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pressly/goose/v3"
)

// EncryptionKey is set from the migration utility to specify default encryption key to migrate to.
// This can remain empty. Then the config will be removed.
//
//nolint:gochecknoglobals // This is a global variable by design.
var EncryptionKey string

func init() {
	up := func(tx *sql.Tx) error {
		var exists bool

		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM connectors.connector)").Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if connectors table exists: %w", err)
		}

		if exists && EncryptionKey == "" {
			return errors.New("encryption key is not set")
		}

		_, err = tx.Exec(`
				CREATE EXTENSION IF NOT EXISTS pgcrypto;
				ALTER TABLE connectors.connector RENAME COLUMN config TO config_unencrypted;
				ALTER TABLE connectors.connector ADD COLUMN config bytea NULL;
		`)
		if err != nil {
			return fmt.Errorf("failed to create config column: %w", err)
		}

		_, err = tx.Exec(`
			UPDATE connectors.connector SET config = pgp_sym_encrypt(config_unencrypted::TEXT, $1, 'compress-algo=1, cipher-algo=aes256');
		`, EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt config: %w", err)
		}

		_, err = tx.Exec(`
			ALTER TABLE connectors.connector DROP COLUMN config_unencrypted;
		`)
		if err != nil {
			return fmt.Errorf("failed to drop config_unencrypted column: %w", err)
		}

		return nil
	}

	down := func(tx *sql.Tx) error {
		var exists bool

		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM connectors.connector)").Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if connectors table exists: %w", err)
		}

		if exists && EncryptionKey == "" {
			return errors.New("encryption key is not set")
		}

		_, err = tx.Exec(`
				ALTER TABLE connectors.connector RENAME COLUMN config TO config_encrypted;
				ALTER TABLE connectors.connector ADD COLUMN config JSON NULL;
		`)
		if err != nil {
			return fmt.Errorf("failed to create config column: %w", err)
		}

		_, err = tx.Exec(`
				UPDATE connectors.connector SET config = pgp_sym_decrypt(config_encrypted, $1, 'compress-algo=1, cipher-algo=aes256')::JSON;
		`, EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt config: %w", err)
		}

		_, err = tx.Exec(`
				ALTER TABLE connectors.connector DROP COLUMN config_encrypted;
		`)
		if err != nil {
			return fmt.Errorf("failed to drop config_encrypted column: %w", err)
		}

		return nil
	}

	goose.AddMigration(up, down)
}
