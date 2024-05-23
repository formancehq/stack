package storage

import (
	"context"

	storageerrors "github.com/formancehq/reconciliation/internal/utils/storage/errors"
	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (s *Storage) CreateRule(ctx context.Context, rule *models.Rule) error {
	_, err := s.db.NewInsert().
		Model(rule).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to create rule", err)
	}
	return nil
}

func (s *Storage) DeleteRule(ctx context.Context, id uint32) error {
	_, err := s.db.NewDelete().
		Model((*models.Rule)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to delete rule", err)
	}
	return nil
}

func (s *Storage) GetRule(ctx context.Context, id uint32) (*models.Rule, error) {
	var rule models.Rule

	err := s.db.NewSelect().
		Model(&rule).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.E("failed to get rule", err)
	}

	return &rule, nil
}

func (s *Storage) buildRuleListQuery(selectQuery *bun.SelectQuery, q ListRulesQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = selectQuery.
		Order("created_at DESC")

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Storage) ListRules(ctx context.Context, q ListRulesQuery) (*bunpaginate.Cursor[models.Rule], error) {
	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = s.ruleQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffset[bunpaginate.PaginatedQueryOptions[RulesFilters], models.Rule](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[RulesFilters]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildRuleListQuery(query, q, where, args)
		},
	)
}

func (s *Storage) ruleQueryContext(qb query.Builder, q ListRulesQuery) (string, []any, error) {
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
		case key == "ruleType":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'ruleType' column can only be used with $match")
			}
			switch ruleType := value.(type) {
			case string:
				return "rule_type = ?", []any{ruleType}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'ruleType' column can only be used with string")
			}
		default:
			return "", nil, errors.Wrapf(storageerrors.ErrInvalidQuery, "unknown key '%s' when building query", key)
		}
	}))
}

type RulesFilters struct{}

type ListRulesQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[RulesFilters]]

func NewListRulesQuery(opts bunpaginate.PaginatedQueryOptions[RulesFilters]) ListRulesQuery {
	return ListRulesQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}
