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

func (s *Storage) CreatePolicy(ctx context.Context, policy *models.Policy) error {
	_, err := s.db.NewInsert().
		Model(policy).
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to create policy", err)
	}
	return nil
}

func (s *Storage) UpdatePolicyRules(ctx context.Context, id uuid.UUID, rules []string) error {
	_, err := s.db.NewUpdate().
		Model((*models.Policy)(nil)).
		Set("rules = ?", rules).
		Where("id = ?", id).Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to update policy rules", err)
	}
	return nil
}

func (s *Storage) UpdatePolicyStatus(ctx context.Context, id uuid.UUID, enabled bool) error {
	_, err := s.db.NewUpdate().
		Model((*models.Policy)(nil)).
		Set("enabled = ?", enabled).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to change enabled status for policy", err)
	}
	return nil
}

func (s *Storage) DeletePolicy(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model((*models.Policy)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to delete policy", err)
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
		return nil, storageerrors.E("failed to get policy", err)
	}
	return &policy, nil
}

func (s *Storage) buildPolicyListQuery(selectQuery *bun.SelectQuery, q ListPoliciesQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = selectQuery.
		Order("created_at DESC")

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Storage) ListPolicies(ctx context.Context, q ListPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error) {
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

	return paginateWithOffset[bunpaginate.PaginatedQueryOptions[PoliciesFilters], models.Policy](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[PoliciesFilters]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildPolicyListQuery(query, q, where, args)
		},
	)
}

func (s *Storage) policyQueryContext(qb query.Builder, q ListPoliciesQuery) (string, []any, error) {
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
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'ledgernameQuery' column can only be used with string")
			}
		case key == "enabled":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'enabled' column can only be used with $match")
			}
			switch name := value.(type) {
			case bool:
				return "enabled = ?", []any{name}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'enabled' column can only be used with string")
			}
		default:
			return "", nil, errors.Wrapf(storageerrors.ErrInvalidQuery, "unknown key '%s' when building query", key)
		}
	}))
}

type PoliciesFilters struct{}

type ListPoliciesQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[PoliciesFilters]]

func NewListPoliciesQuery(opts bunpaginate.PaginatedQueryOptions[PoliciesFilters]) ListPoliciesQuery {
	return ListPoliciesQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}
