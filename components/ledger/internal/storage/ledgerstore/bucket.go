package ledgerstore

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/formancehq/ledger/internal/storage/sqlutils"
	"github.com/formancehq/ledger/internal/storage/systemstore"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

//go:embed migrations/0-init-schema.sql
var initSchema string

type Bucket struct {
	name        string
	systemStore *systemstore.Store
	db          *bun.DB
}

func (b *Bucket) getMigrator() *migrations.Migrator {
	migrator := migrations.NewMigrator(migrations.WithSchema(b.name, true))
	registerMigrations(migrator, b.name)
	return migrator
}

func (b *Bucket) Migrate(ctx context.Context) error {
	return sqlutils.PostgresError(b.getMigrator().Up(ctx, b.db))
}

func (b *Bucket) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	return b.getMigrator().GetMigrations(ctx, b.db)
}

func (b *Bucket) IsUpToDate(ctx context.Context) (bool, error) {
	return b.getMigrator().IsUpToDate(ctx, b.db)
}

func (b *Bucket) Close() error {
	return b.db.Close()
}

func (b *Bucket) newLedgerStore(name string) (*Store, error) {
	return New(b.db, name, func(ctx context.Context) error {
		return b.systemStore.DeleteLedger(ctx, name)
	})
}

func (b *Bucket) createLedgerStore(ctx context.Context, name string) (*Store, error) {

	ledgerExists, err := b.systemStore.ExistsLedger(ctx, name)
	if err != nil {
		return nil, err
	}
	if ledgerExists {
		return nil, sqlutils.ErrStoreAlreadyExists
	}

	_, err = b.systemStore.RegisterLedger(ctx, name)
	if err != nil {
		return nil, err
	}

	store, err := b.newLedgerStore(name)
	if err != nil {
		return nil, err
	}

	err = store.Initialize(ctx)
	if err != nil {
		return nil, err
	}

	return store, err
}

func (b *Bucket) CreateLedgerStore(ctx context.Context, name string) (*Store, error) {
	return b.createLedgerStore(ctx, name)
}

func (b *Bucket) GetLedgerStore(ctx context.Context, name string) (*Store, error) {
	exists, err := b.systemStore.ExistsLedger(ctx, name)
	if err != nil {
		return nil, err
	}

	var store *Store
	if !exists {
		store, err = b.createLedgerStore(ctx, name)
	} else {
		store, err = b.newLedgerStore(name)
	}
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (b *Bucket) DB() *bun.DB {
	return b.db
}

func registerMigrations(migrator *migrations.Migrator, name string) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {

				needV1Upgrade := false
				row := tx.QueryRowContext(ctx, `select exists (
					select from pg_tables
					where schemaname = ? and tablename  = 'log'
				)`, name)
				if row.Err() != nil {
					return row.Err()
				}
				var ret string
				if err := row.Scan(&ret); err != nil {
					panic(err)
				}
				needV1Upgrade = ret != "false"

				oldSchemaRenamed := fmt.Sprintf(name + oldSchemaRenameSuffix)
				if needV1Upgrade {
					_, err := tx.ExecContext(ctx, fmt.Sprintf(`alter schema "%s" rename to "%s"`, name, oldSchemaRenamed))
					if err != nil {
						return errors.Wrap(err, "renaming old schema")
					}
					_, err = tx.ExecContext(ctx, fmt.Sprintf(`create schema if not exists "%s"`, name))
					if err != nil {
						return errors.Wrap(err, "creating new schema")
					}
				}

				_, err := tx.ExecContext(ctx, initSchema)
				if err != nil {
					return errors.Wrap(err, "initializing new schema")
				}

				if needV1Upgrade {
					if err := migrateLogs(ctx, oldSchemaRenamed, name, tx); err != nil {
						return errors.Wrap(err, "migrating logs")
					}

					_, err = tx.ExecContext(ctx, fmt.Sprintf(`create table goose_db_version as table "%s".goose_db_version with no data`, oldSchemaRenamed))
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
	)
}

func ConnectToBucket(systemStore *systemstore.Store, connectionOptions sqlutils.ConnectionOptions, name string) (*Bucket, error) {
	db, err := sqlutils.OpenDBWithSchema(connectionOptions, name)
	if err != nil {
		return nil, sqlutils.PostgresError(err)
	}

	return &Bucket{
		db:          db,
		name:        name,
		systemStore: systemStore,
	}, nil
}
