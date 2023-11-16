package ledgerstore

import (
	"bytes"
	"context"
	"database/sql"
	"text/template"

	"github.com/pkg/errors"

	"github.com/formancehq/stack/libs/go-libs/migrations"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
)

type Store struct {
	bucket *Bucket

	name string
}

func (store *Store) Name() string {
	return store.name
}

func (store *Store) GetDB() *bun.DB {
	return store.bucket.db
}

func (store *Store) Delete(ctx context.Context) error {

	tx, err := store.bucket.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const sqlScript = `
		drop table "transactions_{{.Ledger}}" cascade;
		drop table "transactions_metadata_{{.Ledger}}" cascade;
		drop table "accounts_{{.Ledger}}" cascade;
		drop table "accounts_metadata_{{.Ledger}}" cascade;
		drop table "moves_{{.Ledger}}" cascade;
		drop table "logs_{{.Ledger}}" cascade;
`
	buf := bytes.NewBufferString("")
	if err := template.Must(template.New("delete-ledger").Parse(sqlScript)).Execute(buf, map[string]any{
		"Ledger": store.name,
	}); err != nil {
		return sqlutils.PostgresError(err)
	}

	_, err = tx.ExecContext(ctx, buf.String())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (store *Store) withTransaction(ctx context.Context, callback func(tx bun.Tx) error) error {
	return store.bucket.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		return callback(tx)
	})
}

func (store *Store) IsUpToDate(ctx context.Context) (bool, error) {
	return store.bucket.IsUpToDate(ctx)
}

func (store *Store) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	return store.bucket.GetMigrationsInfo(ctx)
}

func New(
	bucket *Bucket,
	name string,
) (*Store, error) {
	return &Store{
		bucket: bucket,
		name:   name,
	}, nil
}

func InitializeLedgerStore(ctx context.Context, db bun.IDB, name string) error {
	return db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		const sqlScript = `
		create table "transactions_{{.Ledger}}" partition of transactions for values in ('{{.Ledger}}');
		create table "transactions_metadata_{{.Ledger}}" partition of transactions_metadata for values in ('{{.Ledger}}');;
		create table "accounts_{{.Ledger}}" partition of accounts for values in ('{{.Ledger}}');
		create table "accounts_metadata_{{.Ledger}}" partition of accounts_metadata for values in ('{{.Ledger}}');;
		create table "moves_{{.Ledger}}" partition of moves for values in ('{{.Ledger}}');;
		create table "logs_{{.Ledger}}" partition of logs for values in ('{{.Ledger}}');;

		/** Define the trigger which populate table in response to new logs **/
		create trigger "insert_log_{{.Ledger}}" after insert on "logs_{{.Ledger}}"
		for each row execute procedure handle_log();
`
		buf := bytes.NewBufferString("")
		if err := template.Must(template.New("init-ledger").Parse(sqlScript)).Execute(buf, map[string]any{
			"Ledger": name,
		}); err != nil {
			return errors.Wrap(err, "templating init ledger script")
		}

		_, err := tx.ExecContext(ctx, buf.String())
		return errors.Wrap(sqlutils.PostgresError(err), "executing store sql script")
	})
}
