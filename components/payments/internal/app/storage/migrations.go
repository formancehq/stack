package storage

import (
	"context"
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

// EncryptionKey is set from the migration utility to specify default encryption key to migrate to.
// This can remain empty. Then the config will be removed.
//
//nolint:gochecknoglobals // This is a global variable by design.
var EncryptionKey string

func Migrate(ctx context.Context, db *bun.DB) error {
	migrator := migrations.NewMigrator()
	registerMigrations(migrator)

	return migrator.Up(ctx, db)
}

func registerMigrations(migrator *migrations.Migrator) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					CREATE SCHEMA IF NOT EXISTS connectors;
					CREATE SCHEMA IF NOT EXISTS tasks;
					CREATE SCHEMA IF NOT EXISTS accounts;
					CREATE SCHEMA IF NOT EXISTS payments;
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					CREATE TYPE "public".connector_provider AS ENUM ('BANKING-CIRCLE', 'CURRENCY-CLOUD', 'DUMMY-PAY', 'MODULR', 'STRIPE', 'WISE');;
					CREATE TABLE connectors.connector (
					   id uuid  NOT NULL DEFAULT gen_random_uuid(),
					   created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
					   provider connector_provider  NOT NULL UNIQUE,
					   enabled boolean  NOT NULL DEFAULT false,
					   config json NULL,
					   CONSTRAINT connector_pk PRIMARY KEY (id)
					);
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					CREATE TYPE "public".task_status AS ENUM ('STOPPED', 'PENDING', 'ACTIVE', 'TERMINATED', 'FAILED');;
					CREATE TABLE tasks.task (
						id uuid  NOT NULL DEFAULT gen_random_uuid(),
						connector_id uuid  NOT NULL,
						created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
						updated_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=updated_at),
						name text  NOT NULL,
						descriptor json  NULL,
						status task_status  NOT NULL,
						error text  NULL,
						state json  NULL,
						CONSTRAINT task_pk PRIMARY KEY (id)
					);
					ALTER TABLE tasks.task ADD CONSTRAINT task_connector
						FOREIGN KEY (connector_id)
						REFERENCES connectors.connector (id)
						ON DELETE CASCADE
						NOT DEFERRABLE
						INITIALLY IMMEDIATE
					;
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					CREATE TYPE "public".account_type AS ENUM('SOURCE', 'TARGET', 'UNKNOWN');;

					CREATE TABLE accounts.account (
						id uuid  NOT NULL DEFAULT gen_random_uuid(),
						created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
						reference text  NOT NULL UNIQUE,
						provider text  NOT NULL,
						type account_type  NOT NULL,
						CONSTRAINT account_pk PRIMARY KEY (id)
					);
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					CREATE TYPE "public".payment_type AS ENUM ('PAY-IN', 'PAYOUT', 'TRANSFER', 'OTHER');
					CREATE TYPE "public".payment_status AS ENUM ('SUCCEEDED', 'CANCELLED', 'FAILED', 'PENDING', 'OTHER');;

					CREATE TABLE payments.adjustment (
						id uuid  NOT NULL DEFAULT gen_random_uuid(),
						payment_id uuid  NOT NULL,
						created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
						amount bigint NOT NULL DEFAULT 0,
						reference text  NOT NULL UNIQUE,
						status payment_status  NOT NULL,
						absolute boolean  NOT NULL DEFAULT FALSE,
						raw_data json NULL,
						CONSTRAINT adjustment_pk PRIMARY KEY (id)
					);

					CREATE TABLE payments.metadata (
						payment_id uuid  NOT NULL,
						created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
						key text  NOT NULL,
						value text  NOT NULL,
						changelog jsonb NOT NULL,
						CONSTRAINT metadata_pk PRIMARY KEY (payment_id,key)
					);

					CREATE TABLE payments.payment (
						id uuid  NOT NULL DEFAULT gen_random_uuid(),
						connector_id uuid  NOT NULL,
						account_id uuid DEFAULT NULL,
						created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
						reference text  NOT NULL UNIQUE,
						type payment_type  NOT NULL,
						status payment_status  NOT NULL,
						amount bigint NOT NULL DEFAULT 0,
						raw_data json  NULL,
						scheme text  NOT NULL,
						asset text  NOT NULL,
						CONSTRAINT payment_pk PRIMARY KEY (id)
					);

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

					ALTER TABLE payments.payment ADD CONSTRAINT payment_account
						FOREIGN KEY (account_id)
						REFERENCES accounts.account (id)
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
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		//nolint:varnamelen
		migrations.Migration{
			Up: func(tx bun.Tx) error {
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
					UPDATE connectors.connector SET config = pgp_sym_encrypt(config_unencrypted::TEXT, ?, 'compress-algo=1, cipher-algo=aes256');
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
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					CREATE TYPE "public".transfer_status AS ENUM ('PENDING', 'SUCCEEDED', 'FAILED');

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
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					ALTER TABLE payments.payment ALTER COLUMN id DROP DEFAULT;

					ALTER TABLE payments.adjustment drop constraint IF EXISTS adjustment_payment;
					ALTER TABLE payments.metadata drop constraint IF EXISTS metadata_payment;
					ALTER TABLE payments.transfers drop constraint IF EXISTS transfer_payment;
					ALTER TABLE payments.payment ALTER COLUMN id TYPE CHARACTER VARYING;
					ALTER TABLE payments.adjustment ALTER COLUMN payment_id TYPE CHARACTER VARYING;
					ALTER TABLE payments.metadata ALTER COLUMN payment_id TYPE CHARACTER VARYING;
					ALTER TABLE payments.transfers ALTER COLUMN payment_id TYPE CHARACTER VARYING;

					ALTER TABLE payments.metadata ADD CONSTRAINT metadata_payment
						FOREIGN KEY (payment_id)
						REFERENCES payments.payment (id)
						ON DELETE CASCADE
						NOT DEFERRABLE
						INITIALLY IMMEDIATE
					;

					ALTER TABLE payments.adjustment ADD CONSTRAINT adjustment_payment
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
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					ALTER TYPE connector_provider ADD VALUE IF NOT EXISTS 'MANGOPAY';
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					ALTER TYPE connector_provider ADD VALUE IF NOT EXISTS 'MONEYCORP';
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					ALTER TABLE tasks.task ADD COLUMN IF NOT EXISTS "scheduler_options" json;
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					ALTER TABLE accounts.account DROP COLUMN IF EXISTS "type";
					DROP TYPE IF EXISTS "public".account_type;
					CREATE TYPE "public".account_type AS ENUM('INTERNAL', 'EXTERNAL', 'UNKNOWN');;
					ALTER TABLE accounts.account ADD COLUMN IF NOT EXISTS "type" "public".account_type;

					ALTER TABLE accounts.account drop constraint IF EXISTS account_reference_key;
					ALTER TABLE accounts.account ADD COLUMN IF NOT EXISTS "raw_data" json;
					ALTER TABLE accounts.account ADD COLUMN IF NOT EXISTS "default_currency" text NOT NULL DEFAULT '';
					ALTER TABLE accounts.account ADD COLUMN IF NOT EXISTS "account_name" text NOT NULL DEFAULT '';

					ALTER TABLE accounts.account ALTER COLUMN id DROP DEFAULT;
					ALTER TABLE payments.payment DROP CONSTRAINT IF EXISTS payment_account;
					ALTER TABLE payments.payment DROP COLUMN IF EXISTS "account_id";
					ALTER TABLE payments.payment ADD COLUMN IF NOT EXISTS "source_account_id" CHARACTER VARYING;
					ALTER TABLE payments.payment ADD COLUMN IF NOT EXISTS "destination_account_id" CHARACTER VARYING;
					ALTER TABLE accounts.account ALTER COLUMN id TYPE CHARACTER VARYING;

					ALTER TABLE payments.payment DROP CONSTRAINT IF EXISTS payment_source_account;
					ALTER TABLE payments.payment ADD CONSTRAINT payment_source_account
						FOREIGN KEY (source_account_id)
						REFERENCES accounts.account (id)
						ON DELETE CASCADE
						NOT DEFERRABLE
						INITIALLY IMMEDIATE
					;

					ALTER TABLE payments.payment DROP CONSTRAINT IF EXISTS payment_destination_account;
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
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				// Since only one connector is inserting accounts,
				// let's just delete the table, since connectors will be
				// resetted. Delete it cascade, or we will have an error
				_, err := tx.Exec(`
					DELETE FROM accounts.account CASCADE;
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				// Since only one connector is inserting accounts,
				// let's just delete the table, since connectors will be
				// resetted. Delete it cascade, or we will have an error
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS accounts.balances (
						created_at timestamp with time zone  NOT NULL,
						account_id CHARACTER VARYING NOT NULL,
						currency text  NOT NULL,
						balance numeric NOT NULL DEFAULT 0,
						last_updated_at timestamp with time zone  NOT NULL,
						PRIMARY KEY (account_id, created_at, currency)
					);

					DROP INDEX IF EXISTS accounts.idx_created_at_account_id_currency;
					CREATE INDEX idx_created_at_account_id_currency ON accounts.balances(account_id, last_updated_at desc, currency);

					ALTER TABLE accounts.balances DROP CONSTRAINT IF EXISTS balances_account;
					ALTER TABLE accounts.balances ADD CONSTRAINT balances_account
						FOREIGN KEY (account_id)
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
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
					ALTER TABLE payments.adjustment ALTER COLUMN amount TYPE numeric;
					ALTER TABLE payments.payment ALTER COLUMN amount TYPE numeric;
					ALTER TABLE payments.transfers ALTER COLUMN amount TYPE numeric;
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				// In this migration, we have to delete the accounts table since
				// we wanna reset the connector, but the connector_id was not
				// added, hence the table will not be cleaned up when resetting.
				_, err := tx.Exec(`
				DELETE FROM accounts.account CASCADE;
				ALTER TABLE accounts.account ADD COLUMN IF NOT EXISTS "connector_id" uuid;

				ALTER TABLE accounts.account ADD CONSTRAINT accounts_connector
				FOREIGN KEY (connector_id)
				REFERENCES connectors.connector (id)
				ON DELETE CASCADE
				NOT DEFERRABLE
				INITIALLY IMMEDIATE
			;
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
				ALTER TABLE accounts.account ADD COLUMN IF NOT EXISTS metadata jsonb;

				CREATE SCHEMA IF NOT EXISTS transfers;

				CREATE TABLE IF NOT EXISTS transfers.transfer_initiation (
					id character varying  NOT NULL,
					reference text,
					created_at timestamp with time zone  NOT NULL,
					updated_at timestamp with time zone  NOT NULL,
					description text NOT NULL,
					type int NOT NULL,
					source_account_id character varying  NOT NULL,
					destination_account_id character varying  NOT NULL,
					provider connector_provider  NOT NULL,
					amount numeric NOT NULL,
					asset text  NOT NULL,
					status int NOT NULL,
					error text,
					PRIMARY KEY (id)
				);

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
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.Exec(`
				ALTER TABLE transfers.transfer_initiation ALTER COLUMN source_account_id DROP NOT NULL;
				ALTER TABLE transfers.transfer_initiation RENAME COLUMN reference TO payment_id;
				ALTER TYPE "public".account_type ADD VALUE IF NOT EXISTS 'EXTERNAL_FORMANCE';

				CREATE TABLE accounts.bank_account (
					id uuid  NOT NULL DEFAULT gen_random_uuid(),
					created_at timestamp with time zone  NOT NULL DEFAULT NOW() CHECK (created_at<=NOW()),
					name text  NOT NULL,
					provider connector_provider  NOT NULL,
					account_number bytea,
					iban bytea,
					swift_bic_code bytea,
					country text,
					CONSTRAINT bank_account_pk PRIMARY KEY (id)
				);
				`)
				if err != nil {
					return err
				}

				return nil
			},
		},
	)
}
