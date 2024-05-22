package storage

import (
	"context"

	storageerrors "github.com/formancehq/reconciliation/internal/utils/storage/errors"
	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (s *Storage) CreateReconciliation(ctx context.Context, reconciliation *models.Reconciliation) error {
	_, err := s.db.NewInsert().
		Model(reconciliation).
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to create reconciliation", err)
	}
	return nil
}

func (s *Storage) GetReconciliation(ctx context.Context, id uuid.UUID) (*models.Reconciliation, error) {
	var reconciliation models.Reconciliation
	err := s.db.NewSelect().
		Model(&reconciliation).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.E("failed to get reconciliation", err)
	}
	return &reconciliation, nil
}

func (s *Storage) buildReconciliationListQuery(selectQuery *bun.SelectQuery, q ListReconciliationsQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = selectQuery.
		Order("created_at DESC")

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Storage) ListReconciliations(ctx context.Context, q ListReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error) {
	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = s.reconciliationsQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffset[bunpaginate.PaginatedQueryOptions[ReconciliationFilters], models.Reconciliation](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[ReconciliationFilters]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildReconciliationListQuery(query, q, where, args)
		},
	)
}

func (s *Storage) reconciliationsQueryContext(qb query.Builder, q ListReconciliationsQuery) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "name":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'name' column can only be used with $match")
			}
			switch name := value.(type) {
			case string:
				return "name = ?", []any{name}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'name' column can only be used with string")
			}

		case key == "policyID":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyID' column can only be used with $match")
			}
			switch policyID := value.(type) {
			case string:
				return "policy_id = ?", []any{policyID}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyID' column can only be used with string")
			}

		case key == "policyType":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyType' column can only be used with $match")
			}
			switch policyType := value.(type) {
			case string:
				return "policy_type = ?", []any{policyType}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyType' column can only be used with string")
			}

		case key == "reconciliationStatus":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'reconciliationStatus' column can only be used with $match")
			}
			switch reconciliationStatus := value.(type) {
			case string:
				return "reconciliation_status = ?", []any{reconciliationStatus}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'reconciliationStatus' column can only be used with string")
			}

		default:
			return "", nil, errors.Wrapf(storageerrors.ErrInvalidQuery, "unknown key '%s' when building query", key)
		}
	}))
}

type ReconciliationFilters struct{}

type ListReconciliationsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[ReconciliationFilters]]

func NewListReconciliationsQuery(opts bunpaginate.PaginatedQueryOptions[ReconciliationFilters]) ListReconciliationsQuery {
	return ListReconciliationsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}
