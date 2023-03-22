package ledger

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage"
	sqlerrors "github.com/formancehq/ledger/pkg/storage/sqlstorage/errors"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/migrations"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/schema"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/worker"
	"github.com/formancehq/stack/libs/go-libs/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

const (
	SQLCustomFuncMetaCompare = "meta_compare"
)

type Store struct {
	schema   schema.Schema
	onClose  func(ctx context.Context) error
	onDelete func(ctx context.Context) error

	logsBatchWorker *worker.Worker[core.Log]
}

func (s *Store) error(err error) error {
	if err == nil {
		return nil
	}
	return sqlerrors.PostgresError(err)
}

func (s *Store) Schema() schema.Schema {
	return s.schema
}

func (s *Store) Name() string {
	return s.schema.Name()
}

func (s *Store) Delete(ctx context.Context) error {
	if err := s.schema.Delete(ctx); err != nil {
		return err
	}
	return errors.Wrap(s.onDelete(ctx), "deleting ledger store")
}

func (s *Store) Initialize(ctx context.Context) (bool, error) {
	logging.FromContext(ctx).Debug("Initialize store")

	ms, err := migrations.CollectMigrationFiles(MigrationsFS)
	if err != nil {
		return false, err
	}

	return migrations.Migrate(ctx, s.schema, ms...)
}

func (s *Store) Close(ctx context.Context) error {
	if err := s.logsBatchWorker.Stop(ctx); err != nil {
		return err
	}
	return s.onClose(ctx)
}

func NewStore(
	ctx context.Context,
	schema schema.Schema,
	onClose, onDelete func(ctx context.Context) error,
) *Store {
	s := &Store{
		schema:   schema,
		onClose:  onClose,
		onDelete: onDelete,
	}

	logsBatchWorker := worker.NewWorker(s.batchLogs)
	s.logsBatchWorker = logsBatchWorker

	go logsBatchWorker.Run(ctx)

	return s
}

var _ storage.LedgerStore = &Store{}
