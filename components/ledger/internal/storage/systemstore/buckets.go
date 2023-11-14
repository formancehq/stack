package systemstore

import (
	"context"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/sqlutils"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Buckets struct {
	bun.BaseModel `bun:"_system.buckets,alias:buckets"`

	Name    string      `bun:"name,type:varchar(255),pk"` // Primary key
	AddedAt ledger.Time `bun:"addedat,type:timestamp"`
}

func (s *Store) ListBuckets(ctx context.Context) ([]string, error) {
	query := s.db.NewSelect().
		Model((*Buckets)(nil)).
		Column("name").
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

func (s *Store) DeleteBucket(ctx context.Context, name string) error {
	_, err := s.db.NewDelete().
		Model((*Buckets)(nil)).
		Where("name = ?", name).
		Exec(ctx)

	return errors.Wrap(sqlutils.PostgresError(err), "delete bucket from system store")
}

func (s *Store) RegisterBucket(ctx context.Context, bucketName string) (bool, error) {
	l := &Buckets{
		Name:    bucketName,
		AddedAt: ledger.Now(),
	}

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

func (s *Store) ExistsBucket(ctx context.Context, bucket string) (bool, error) {
	query := s.db.NewSelect().
		Model((*Buckets)(nil)).
		Column("name").
		Where("name = ?", bucket).
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
