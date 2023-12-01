package ledgerstore

import (
	"bytes"
	"context"
	"database/sql"
	"text/template"

	"github.com/formancehq/stack/libs/go-libs/migrations"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
)

type Store struct {
	onDelete func(ctx context.Context) error
	bucket   *Bucket

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

	if err := store.onDelete(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func (store *Store) withTransaction(ctx context.Context, callback func(tx bun.Tx) error) error {
	tx, err := store.bucket.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := callback(tx); err != nil {
		_ = tx.Rollback()
		return sqlutils.PostgresError(err)
	}
	return tx.Commit()
}

func (store *Store) Initialize(ctx context.Context) error {
	tx, err := store.bucket.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return sqlutils.PostgresError(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

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

		create trigger "update_account_{{.Ledger}}" after update on "accounts_{{.Ledger}}"
        for each row execute procedure update_account_metadata_history();

		create trigger "insert_account_{{.Ledger}}" after insert on "accounts_{{.Ledger}}"
        for each row execute procedure insert_account_metadata_history();

		create trigger "update_transaction_{{.Ledger}}" after update on "transactions_{{.Ledger}}"
        for each row execute procedure update_transaction_metadata_history();

		create trigger "insert_transaction_{{.Ledger}}" after insert on "transactions_{{.Ledger}}"
        for each row execute procedure insert_transaction_metadata_history();
`
	buf := bytes.NewBufferString("")
	if err := template.Must(template.New("init-ledger").Parse(sqlScript)).Execute(buf, map[string]any{
		"Ledger": store.name,
	}); err != nil {
		return sqlutils.PostgresError(err)
	}

	_, err = tx.ExecContext(ctx, buf.String())
	if err != nil {
		return sqlutils.PostgresError(err)
	}

	return tx.Commit()
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
	onDelete func(ctx context.Context) error,
) (*Store, error) {
	return &Store{
		bucket:   bucket,
		name:     name,
		onDelete: onDelete,
	}, nil
}
