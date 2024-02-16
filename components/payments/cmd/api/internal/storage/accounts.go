package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
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

func (s *Storage) ListAccounts(ctx context.Context, q ListAccountsQuery) (*api.Cursor[models.Account], error) {
	return PaginateWithOffset[PaginatedQueryOptions[AccountQuery], models.Account](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[AccountQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Relation("PoolAccounts").
				Order("created_at DESC")

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			}

			return query
		},
	)
}

func (s *Storage) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account

	err := s.db.NewSelect().
		Model(&account).
		Relation("Connector").
		Relation("PoolAccounts").
		Where("account.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get account", err)
	}

	return &account, nil
}
