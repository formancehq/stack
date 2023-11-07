package storage

import (
	"context"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/gibson042/canonicaljson-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func registerMigrationsV1(ctx context.Context, migrator *migrations.Migrator) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "",
			Up: func(tx bun.Tx) error {
				if err := migrateConnectors(ctx, tx); err != nil {
					return err
				}

				if err := migrateAccountID(ctx, tx); err != nil {
					return err
				}

				if err := migratePaymentID(ctx, tx); err != nil {
					return err
				}

				if err := migrateTransferInitiationID(ctx, tx); err != nil {
					return err
				}

				return nil
			},
		},
	)
}

type PreviousConnector struct {
	bun.BaseModel `bun:"connectors.connector"`

	ID        uuid.UUID `bun:",pk,nullzero"`
	CreatedAt time.Time `bun:",nullzero"`
	Provider  models.ConnectorProvider

	// EncryptedConfig is a PGP-encrypted JSON string.
	EncryptedConfig string `bun:"config"`

	// Config is a decrypted config. It is not stored in the database.
	Config json.RawMessage `bun:"decrypted_config,scanonly"`
}

type Connector struct {
	bun.BaseModel `bun:"connectors.connector_v2"`

	ID        models.ConnectorID `bun:",pk,nullzero"`
	Name      string
	CreatedAt time.Time `bun:",nullzero"`
	Provider  models.ConnectorProvider

	// EncryptedConfig is a PGP-encrypted JSON string.
	EncryptedConfig string `bun:"config"`
}

func migrateConnectors(ctx context.Context, tx bun.Tx) error {
	_, err := tx.Exec(`
	ALTER TABLE accounts.account ALTER COLUMN connector_id SET NOT NULL;
	ALTER TABLE accounts.bank_account ALTER COLUMN connector_id SET NOT NULL;

	CREATE TABLE connectors.connector_v2 (
		id CHARACTER VARYING  NOT NULL,
		name text NOT NULL UNIQUE,
		created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
		provider connector_provider  NOT NULL,
		config bytea NULL,
		CONSTRAINT connector_v2_pk PRIMARY KEY (id)
	);
	`)
	if err != nil {
		return err
	}

	var oldConnectors []*PreviousConnector
	err = tx.NewSelect().
		Model(&oldConnectors).
		Scan(ctx)
	if err != nil {
		return err
	}

	newConnectors := make([]*Connector, 0, len(oldConnectors))
	for _, oldConnector := range oldConnectors {
		newConnectors = append(newConnectors, &Connector{
			ID: models.ConnectorID{
				Reference: uuid.New(),
				Provider:  oldConnector.Provider,
			},
			Name:            oldConnector.Provider.String(),
			CreatedAt:       oldConnector.CreatedAt,
			Provider:        oldConnector.Provider,
			EncryptedConfig: oldConnector.EncryptedConfig,
		})
	}

	if len(newConnectors) > 0 {
		_, err = tx.NewInsert().
			Model(&newConnectors).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
	ALTER TABLE tasks.task ADD COLUMN IF NOT EXISTS provider connector_provider;
	UPDATE tasks.task SET provider = (SELECT provider FROM connectors.connector WHERE id = task.connector_id);
	ALTER TABLE payments.payment ADD COLUMN IF NOT EXISTS provider connector_provider;
	UPDATE payments.payment SET provider = (SELECT provider FROM connectors.connector WHERE id = payment.connector_id);
	ALTER TABLE tasks.task DROP CONSTRAINT IF EXISTS task_connector;
	ALTER TABLE tasks.task ALTER COLUMN connector_id TYPE CHARACTER VARYING;
	ALTER TABLE payments.payment DROP CONSTRAINT IF EXISTS payment_connector;
	ALTER TABLE payments.payment ALTER COLUMN connector_id TYPE CHARACTER VARYING;
	ALTER TABLE payments.transfers DROP CONSTRAINT IF EXISTS transfer_connector;
	ALTER TABLE payments.transfers ALTER COLUMN connector_id TYPE CHARACTER VARYING;
	ALTER TABLE accounts.account DROP CONSTRAINT IF EXISTS accounts_connector;
	ALTER TABLE accounts.account ALTER COLUMN connector_id TYPE CHARACTER VARYING;
	ALTER TABLE accounts.bank_account DROP CONSTRAINT IF EXISTS bank_accounts_connector;
	ALTER TABLE accounts.bank_account ALTER COLUMN connector_id TYPE CHARACTER VARYING;
	ALTER TABLE transfers.transfer_initiation DROP CONSTRAINT IF EXISTS transfer_initiation_connector_id;

	DROP TABLE connectors.connector;
	ALTER TABLE connectors.connector_v2 RENAME TO connector;

	UPDATE tasks.task SET connector_id = (SELECT id FROM connectors.connector WHERE provider = task.provider);
	UPDATE payments.payment SET connector_id = (SELECT id FROM connectors.connector WHERE provider = payment.provider);
	UPDATE accounts.account SET connector_id = (SELECT id FROM connectors.connector WHERE provider::text = account.provider);
	UPDATE accounts.bank_account SET connector_id = (SELECT id FROM connectors.connector WHERE provider = bank_account.provider);

	ALTER TABLE tasks.task DROP COLUMN IF EXISTS provider;
	ALTER TABLE accounts.account DROP COLUMN IF EXISTS provider;
	ALTER TABLE accounts.bank_account DROP COLUMN IF EXISTS provider;
	ALTER TABLE payments.payment DROP COLUMN IF EXISTS provider;

	ALTER TABLE tasks.task ADD CONSTRAINT task_connector
	FOREIGN KEY (connector_id)
	REFERENCES connectors.connector (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;
	
	ALTER TABLE payments.payment ADD CONSTRAINT payment_connector
	FOREIGN KEY (connector_id)
	REFERENCES connectors.connector (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE accounts.account ADD CONSTRAINT accounts_connector
	FOREIGN KEY (connector_id)
	REFERENCES connectors.connector (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE accounts.bank_account ADD CONSTRAINT bank_accounts_connector
	FOREIGN KEY (connector_id)
	REFERENCES connectors.connector (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE transfers.transfer_initiation ADD COLUMN IF NOT EXISTS connector_id CHARACTER VARYING  NOT NULL;
	ALTER TABLE transfers.transfer_initiation ADD CONSTRAINT transfer_initiation_connector_id
	FOREIGN KEY (connector_id)
	REFERENCES connectors.connector (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	UPDATE transfers.transfer_initiation SET connector_id = (SELECT id FROM connectors.connector WHERE provider = transfer_initiation.provider);
	`)
	if err != nil {
		return err
	}

	return nil
}

type PreviousAccountID struct {
	Reference string
	Provider  models.ConnectorProvider
}

func PreviousAccountIDFromString(value string) (*PreviousAccountID, error) {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	ret := PreviousAccountID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (aid *PreviousAccountID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("account id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := PreviousAccountIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse account id %s: %v", v, err)
			}

			*aid = *id
			return nil
		}
	}

	return fmt.Errorf("failed to scan account id: %v", value)
}

func (aid PreviousAccountID) Value() (driver.Value, error) {
	return aid.String(), nil
}

func (aid *PreviousAccountID) String() string {
	if aid == nil || aid.Reference == "" {
		return ""
	}

	data, err := canonicaljson.Marshal(aid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(data)
}

func migrateAccountID(ctx context.Context, tx bun.Tx) error {
	_, err := tx.Exec(`
	ALTER TABLE accounts.balances DROP CONSTRAINT IF EXISTS balances_account;
	ALTER TABLE accounts.bank_account DROP CONSTRAINT IF EXISTS bank_account_account_id;
	ALTER TABLE transfers.transfer_initiation DROP CONSTRAINT IF EXISTS destination_account;
	ALTER TABLE transfers.transfer_initiation DROP CONSTRAINT IF EXISTS source_account;
	ALTER TABLE payments.payment DROP CONSTRAINT IF EXISTS payment_destination_account;
	ALTER TABLE payments.payment DROP CONSTRAINT IF EXISTS payment_source_account;
	`)
	if err != nil {
		return err
	}

	var previousIDs []PreviousAccountID
	var connectorIDs []models.ConnectorID
	err = tx.NewSelect().Model((*models.Account)(nil)).Column("id", "connector_id").Scan(ctx, &previousIDs, &connectorIDs)
	if err != nil {
		return err
	}

	if len(previousIDs) != len(connectorIDs) {
		return fmt.Errorf("migrateAccountID: previousIDs and connectorIDs have different length")
	}

	type AccoutIDMigration struct {
		PreviousAccountID PreviousAccountID
		NewAccountID      models.AccountID
	}
	migrations := make([]AccoutIDMigration, 0, len(previousIDs))
	for i, previousID := range previousIDs {
		migrations = append(migrations, AccoutIDMigration{
			PreviousAccountID: previousID,
			NewAccountID: models.AccountID{
				Reference:   previousID.Reference,
				ConnectorID: connectorIDs[i],
			},
		})
	}

	for _, migration := range migrations {
		_, err := tx.NewUpdate().
			Model((*models.Account)(nil)).
			Set("id = ?", migration.NewAccountID).
			Where("id = ?", migration.PreviousAccountID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.Balance)(nil)).
			Set("account_id = ?", migration.NewAccountID).
			Where("account_id = ?", migration.PreviousAccountID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.BankAccount)(nil)).
			Set("account_id = ?", migration.NewAccountID).
			Where("account_id = ?", migration.PreviousAccountID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.TransferInitiation)(nil)).
			Set("source_account = ?", migration.NewAccountID).
			Where("source_account = ?", migration.PreviousAccountID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.TransferInitiation)(nil)).
			Set("destination_account = ?", migration.NewAccountID).
			Where("destination_account = ?", migration.PreviousAccountID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.Payment)(nil)).
			Set("source_account = ?", migration.NewAccountID).
			Where("source_account = ?", migration.PreviousAccountID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.Payment)(nil)).
			Set("destination_account = ?", migration.NewAccountID).
			Where("destination_account = ?", migration.PreviousAccountID).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
	ALTER TABLE transfers.transfer_initiation ADD CONSTRAINT source_account
	FOREIGN KEY (source_account_id)
	REFERENCES accounts.account (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE transfers.transfer_initiation ADD CONSTRAINT destination_account
	FOREIGN KEY (destination_account_id)
	REFERENCES accounts.account (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE accounts.balances ADD CONSTRAINT balances_account
	FOREIGN KEY (account_id)
	REFERENCES accounts.account (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE accounts.bank_account ADD CONSTRAINT bank_account_account_id
	FOREIGN KEY (account_id)
	REFERENCES accounts.account (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE payments.payment ADD CONSTRAINT payment_source_account
	FOREIGN KEY (source_account_id)
	REFERENCES accounts.account (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE payments.payment ADD CONSTRAINT payment_destination_account
	FOREIGN KEY (destination_account_id)
	REFERENCES accounts.account (id)
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

type PreviousPaymentID struct {
	models.PaymentReference
	Provider models.ConnectorProvider
}

func (pid PreviousPaymentID) String() string {
	data, err := canonicaljson.Marshal(pid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(data)
}

func PaymentIDFromString(value string) (*PreviousPaymentID, error) {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	ret := PreviousPaymentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (pid PreviousPaymentID) Value() (driver.Value, error) {
	return pid.String(), nil
}

func (pid *PreviousPaymentID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := PaymentIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse paymentid %s: %v", v, err)
			}

			*pid = *id
			return nil
		}
	}

	return fmt.Errorf("failed to scan paymentid: %v", value)
}

func migratePaymentID(ctx context.Context, tx bun.Tx) error {
	_, err := tx.Exec(`
	ALTER TABLE payments.adjustment DROP CONSTRAINT IF EXISTS adjustment_payment;
	ALTER TABLE payments.metadata DROP CONSTRAINT IF EXISTS metadata_payment;
	`)
	if err != nil {
		return err
	}

	var previousIDs []PreviousPaymentID
	var connectorIDs []models.ConnectorID
	err = tx.NewSelect().Model((*models.Payment)(nil)).Column("id", "connector_id").Scan(ctx, &previousIDs, &connectorIDs)
	if err != nil {
		return err
	}

	if len(previousIDs) != len(connectorIDs) {
		return fmt.Errorf("migrateAccountID: previousIDs and connectorIDs have different length")
	}

	type PaymentIDMigration struct {
		PreviousPaymentID PreviousPaymentID
		NewPaymentID      models.PaymentID
	}
	migrations := make([]PaymentIDMigration, 0, len(previousIDs))
	for i, previousID := range previousIDs {
		migrations = append(migrations, PaymentIDMigration{
			PreviousPaymentID: previousID,
			NewPaymentID: models.PaymentID{
				PaymentReference: previousID.PaymentReference,
				ConnectorID:      connectorIDs[i],
			},
		})
	}

	for _, migration := range migrations {
		_, err := tx.NewUpdate().
			Model((*models.Payment)(nil)).
			Set("id = ?", migration.NewPaymentID).
			Where("id = ?", migration.PreviousPaymentID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.Adjustment)(nil)).
			Set("payment_id = ?", migration.NewPaymentID).
			Where("payment_id = ?", migration.PreviousPaymentID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.Metadata)(nil)).
			Set("payment_id = ?", migration.NewPaymentID).
			Where("payment_id = ?", migration.PreviousPaymentID).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
	ALTER TABLE payments.adjustment ADD CONSTRAINT adjustment_payment
	FOREIGN KEY (payment_id)
	REFERENCES payments.payment (id)
	ON DELETE CASCADE
	NOT DEFERRABLE
	INITIALLY IMMEDIATE
	;

	ALTER TABLE payments.metadata ADD CONSTRAINT metadata_payment
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

type PreviousTransferInitiationID struct {
	Reference string
	Provider  models.ConnectorProvider
}

func (tid PreviousTransferInitiationID) String() string {
	data, err := canonicaljson.Marshal(tid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(data)
}

func TransferInitiationIDFromString(value string) (PreviousTransferInitiationID, error) {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return PreviousTransferInitiationID{}, err
	}
	ret := PreviousTransferInitiationID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return PreviousTransferInitiationID{}, err
	}

	return ret, nil
}

func (tid PreviousTransferInitiationID) Value() (driver.Value, error) {
	return tid.String(), nil
}

func (tid *PreviousTransferInitiationID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := TransferInitiationIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse paymentid %s: %v", v, err)
			}

			*tid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan paymentid: %v", value)
}

func migrateTransferInitiationID(ctx context.Context, tx bun.Tx) error {
	_, err := tx.Exec(`
	ALTER TABLE transfers.transfer_initiation_payments DROP CONSTRAINT IF EXISTS transfer_initiation_id_constraint;
	`)
	if err != nil {
		return err
	}

	var previousIDs []PreviousTransferInitiationID
	var connectorIDs []models.ConnectorID
	err = tx.NewSelect().Model((*models.TransferInitiation)(nil)).Column("id", "connector_id").Scan(ctx, &previousIDs, &connectorIDs)
	if err != nil {
		return err
	}

	if len(previousIDs) != len(connectorIDs) {
		return fmt.Errorf("migrateAccountID: previousIDs and connectorIDs have different length")
	}

	type TransferInitiationIDMigration struct {
		PreviousTransferInitiationID PreviousTransferInitiationID
		NewTransferInitiationID      models.TransferInitiationID
	}

	migrations := make([]TransferInitiationIDMigration, 0, len(previousIDs))
	for i, previousID := range previousIDs {
		migrations = append(migrations, TransferInitiationIDMigration{
			PreviousTransferInitiationID: previousID,
			NewTransferInitiationID: models.TransferInitiationID{
				Reference:   previousID.Reference,
				ConnectorID: connectorIDs[i],
			},
		})
	}

	for _, migration := range migrations {
		_, err := tx.NewUpdate().
			Model((*models.TransferInitiation)(nil)).
			Set("id = ?", migration.NewTransferInitiationID).
			Where("id = ?", migration.PreviousTransferInitiationID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*models.TransferInitiationPayments)(nil)).
			Set("transfer_initiation_id = ?", migration.NewTransferInitiationID).
			Where("transfer_initiation_id = ?", migration.PreviousTransferInitiationID).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
	ALTER TABLE transfers.transfer_initiation_payments ADD CONSTRAINT transfer_initiation_id_constraint
	FOREIGN KEY (transfer_initiation_id)
	REFERENCES transfers.transfer_initiation (id)
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
