package storage

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/query"
	"github.com/formancehq/reconciliation/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (s *Storage) CreateReconciation(ctx context.Context, reco *models.Reconciliation) error {
	_, err := s.db.NewInsert().
		Model(reco).
		Exec(ctx)
	if err != nil {
		return e("failed to create reconciliation", err)
	}

	return nil
}

func (s *Storage) GetReconciliation(ctx context.Context, id uuid.UUID) (*models.Reconciliation, error) {
	var reco models.Reconciliation
	err := s.db.NewSelect().
		Model(&reco).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get reconciliation", err)
	}

	return &reco, nil
}

func (s *Storage) buildReconciliationListQuery(selectQuery *bun.SelectQuery, q GetReconciliationsQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = selectQuery.
		Order("created_at DESC")

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Storage) ListReconciliations(ctx context.Context, q GetReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error) {
	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = s.reconciliationQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffset[PaginatedQueryOptions[ReconciliationsFilters], models.Reconciliation](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[ReconciliationsFilters]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildReconciliationListQuery(query, q, where, args)
		},
	)
}

func (s *Storage) reconciliationQueryContext(qb query.Builder, q GetReconciliationsQuery) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "policyID":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrInvalidQuery, "'policyID' column can only be used with $match")
			}
			switch pID := value.(type) {
			case string:
				return "policy_id = ?", []any{pID}, nil
			default:
				return "", nil, errors.Wrap(ErrInvalidQuery, "'policyID' column can only be used with string")
			}
		default:
			return "", nil, errors.Wrapf(ErrInvalidQuery, "unknown key '%s' when building query", key)
		}
	}))
}

type ReconciliationsFilters struct{}

type GetReconciliationsQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[ReconciliationsFilters]]

func NewGetReconciliationsQuery(opts PaginatedQueryOptions[ReconciliationsFilters]) GetReconciliationsQuery {
	return GetReconciliationsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}
