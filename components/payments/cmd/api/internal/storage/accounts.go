package storage

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type AccountQuery struct{}

type ListAccountsQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[AccountQuery]]

func NewListAccountsQuery(opts PaginatedQueryOptions[AccountQuery]) ListAccountsQuery {
	return ListAccountsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *Storage) accountQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "reference":
			return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil
		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'metadata' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			key := "metadata"
			return key + " @> ?", []any{map[string]any{
				match[0][1]: value,
			}}, nil
		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))
}

func (s *Storage) ListAccounts(ctx context.Context, q ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.accountQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	return PaginateWithOffset[PaginatedQueryOptions[AccountQuery], models.Account](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[AccountQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.Relation("PoolAccounts")

			if where != "" {
				query = query.Where(where, args...)
			}

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			} else {
				query = query.Order("created_at DESC")
			}

			return query
		},
	)
}

func (s *Storage) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account

	err := s.db.NewSelect().
		Model(&account).
		Relation("PoolAccounts").
		Where("account.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get account", err)
	}

	return &account, nil
}

func (s *Storage) UpsertAccounts(ctx context.Context, accounts []*models.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	_, err := s.db.NewInsert().
		Model(&accounts).
		On("CONFLICT (id) DO UPDATE").
		Set("connector_id = EXCLUDED.connector_id").
		Set("raw_data = EXCLUDED.raw_data").
		Set("default_currency = EXCLUDED.default_currency").
		Set("account_name = EXCLUDED.account_name").
		Set("metadata = EXCLUDED.metadata").
		Exec(ctx)
	if err != nil {
		return e("failed to create accounts", err)
	}

	return nil
}
