package ledger

import (
	"context"

	sqlerrors "github.com/formancehq/ledger/pkg/storage/sqlstorage/errors"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/schema"
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

	return false, nil
	// TODO(polo): handle migration
	// migrations, err := CollectMigrationFiles(MigrationsFS)
	// if err != nil {
	// 	return false, err
	// }

	// return Migrate(ctx, s.schema, migrations...)
}

func (s *Store) Close(ctx context.Context) error {
	return s.onClose(ctx)
}

func NewStore(
	schema schema.Schema,
	onClose, onDelete func(ctx context.Context) error,
) *Store {
	return &Store{
		schema:   schema,
		onClose:  onClose,
		onDelete: onDelete,
	}
}

//------------------------------------------------------------------------------

// TODO(polo): Reinstate this when we have a proper ledger store defined in
// storage package.
// var _ ledger.Store = &Store{}
