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

	Name    string      `bun:"ledger,type:varchar(255),pk"` // Primary key
	AddedAt ledger.Time `bun:"addedat,type:timestamp"`
	Bucket  string      `bun:"bucket,type:varchar(255)"`
}

func (s *Store) ListLedgers(ctx context.Context) ([]string, error) {
	query := s.db.NewSelect().
		Model((*Ledger)(nil)).
		Column("ledger").
		String()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, sqlutils.PostgresError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	res := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, sqlutils.PostgresError(err)
		}
		res = append(res, name)
	}
	return res, nil
}

func (s *Store) DeleteLedger(ctx context.Context, name string) error {
	_, err := s.db.NewDelete().
		Model((*Ledger)(nil)).
		Where("ledger = ?", name).
		Exec(ctx)

	return errors.Wrap(sqlutils.PostgresError(err), "delete ledger from system store")
}

func (s *Store) RegisterLedger(ctx context.Context, l *Ledger) (bool, error) {
	ret, err := s.db.NewInsert().
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

func (s *Store) ExistsLedger(ctx context.Context, ledger string) (bool, error) {
	query := s.db.NewSelect().
		Model((*Ledger)(nil)).
		Column("ledger").
		Where("ledger = ?", ledger).
		String()

	ret := s.db.QueryRowContext(ctx, query)
	if ret.Err() != nil {
		return false, nil
	}

	var t string
	_ = ret.Scan(&t) // Trigger close

	if t == "" {
		return false, nil
	}
	return true, nil
}

func (s *Store) GetLedger(ctx context.Context, name string) (*Ledger, error) {
	ret := &Ledger{}
	if err := s.db.NewSelect().
		Model(ret).
		Column("ledger").
		Where("ledger = ?", name).
		Scan(ctx); err != nil {
		return nil, err
	}

	return ret, nil
}
