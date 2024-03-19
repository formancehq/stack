package storage

import (
	"context"

	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (s *Storage) CreatePolicy(ctx context.Context, policy *models.Policy) error {
	_, err := s.db.NewInsert().
		Model(policy).
		Exec(ctx)
	if err != nil {
		return e("failed to create policy", err)
	}

	return nil
}

func (s *Storage) DeletePolicy(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&models.Policy{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("failed to delete policy", err)
	}

	return nil
}

func (s *Storage) GetPolicy(ctx context.Context, id uuid.UUID) (*models.Policy, error) {
	var policy models.Policy
	err := s.db.NewSelect().
		Model(&policy).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get policy", err)
	}

	return &policy, nil
}

func (s *Storage) buildPolicyListQuery(selectQuery *bun.SelectQuery, q GetPoliciesQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = selectQuery.
		Order("created_at DESC")

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

// TODO(polo): add pagination from go libs
func (s *Storage) ListPolicies(ctx context.Context, q GetPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error) {
	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = s.policyQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffset[PaginatedQueryOptions[PoliciesFilters], models.Policy](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PoliciesFilters]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildPolicyListQuery(query, q, where, args)
		},
	)
}

func (s *Storage) policyQueryContext(qb query.Builder, q GetPoliciesQuery) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "ledgerQuery":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrInvalidQuery, "'ledgerQuery' column can only be used with $match")
			}
			switch ledgerQuery := value.(type) {
			case string:
				return "ledger_query = ?", []any{ledgerQuery}, nil
			default:
				return "", nil, errors.Wrap(ErrInvalidQuery, "'ledgerQuery' column can only be used with string")
			}
		case key == "ledgerName":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrInvalidQuery, "'ledgerName' column can only be used with $match")
			}
			switch name := value.(type) {
			case string:
				return "ledger_name = ?", []any{name}, nil
			default:
				return "", nil, errors.Wrap(ErrInvalidQuery, "'ledgerName' column can only be used with string")
			}
		case key == "paymentsPoolID":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrInvalidQuery, "'paymentsPoolID' column can only be used with $match")
			}
			switch pID := value.(type) {
			case string:
				return "payments_pool_id = ?", []any{pID}, nil
			default:
				return "", nil, errors.Wrap(ErrInvalidQuery, "'paymentsPoolID' column can only be used with string")
			}
		default:
			return "", nil, errors.Wrapf(ErrInvalidQuery, "unknown key '%s' when building query", key)
		}
	}))
}

type PoliciesFilters struct{}

type GetPoliciesQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PoliciesFilters]]

func NewGetPoliciesQuery(opts PaginatedQueryOptions[PoliciesFilters]) GetPoliciesQuery {
	return GetPoliciesQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}
