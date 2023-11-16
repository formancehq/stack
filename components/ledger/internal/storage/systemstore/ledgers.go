package systemstore

import (
	"context"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Ledger struct {
	bun.BaseModel `bun:"_system.ledgers,alias:ledgers"`

	Name    string      `bun:"ledger,type:varchar(255),pk" json:"name"` // Primary key
	AddedAt ledger.Time `bun:"addedat,type:timestamp" json:"addedAt"`
	Bucket  string      `bun:"bucket,type:varchar(255)" json:"bucket"`
}

func (s *Store) ListLedgers(ctx context.Context) ([]Ledger, error) {
	ret := make([]Ledger, 0)
	if err := s.db.NewSelect().
		Model(&ret).
		Column("ledger", "bucket", "addedat").
		Scan(ctx); err != nil {
		return nil, sqlutils.PostgresError(err)
	}

	return ret, nil
}

func (s *Store) DeleteLedger(ctx context.Context, name string) error {
	_, err := s.db.NewDelete().
		Model((*Ledger)(nil)).
		Where("ledger = ?", name).
		Exec(ctx)

	return errors.Wrap(sqlutils.PostgresError(err), "delete ledger from system store")
}

func (s *Store) RegisterLedger(ctx context.Context, l *Ledger) (bool, error) {
	return RegisterLedger(ctx, s.db, l)
}

func (s *Store) GetLedger(ctx context.Context, name string) (*Ledger, error) {
	ret := &Ledger{}
	if err := s.db.NewSelect().
		Model(ret).
		Column("ledger", "bucket", "addedat").
		Where("ledger = ?", name).
		Scan(ctx); err != nil {
		return nil, sqlutils.PostgresError(err)
	}

	return ret, nil
}

func RegisterLedger(ctx context.Context, db bun.IDB, l *Ledger) (bool, error) {
	ret, err := db.NewInsert().
		Model(l).
		Ignore().
		Exec(ctx)
	if err != nil {
		return false, sqlutils.PostgresError(err)
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return false, sqlutils.PostgresError(err)
	}

	return affected > 0, nil
}
