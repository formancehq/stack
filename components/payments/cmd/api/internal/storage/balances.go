package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type BalanceQuery struct {
	AccountID *models.AccountID
	Currency  string
	From      time.Time
	To        time.Time
}

func NewBalanceQuery() BalanceQuery {
	return BalanceQuery{}
}

func (b BalanceQuery) WithAccountID(accountID *models.AccountID) BalanceQuery {
	b.AccountID = accountID

	return b
}

func (b BalanceQuery) WithCurrency(currency string) BalanceQuery {
	b.Currency = currency

	return b
}

func (b BalanceQuery) WithFrom(from time.Time) BalanceQuery {
	b.From = from

	return b
}

func (b BalanceQuery) WithTo(to time.Time) BalanceQuery {
	b.To = to

	return b
}

func applyBalanceQuery(query *bun.SelectQuery, balanceQuery BalanceQuery) *bun.SelectQuery {
	if balanceQuery.AccountID != nil {
		query = query.Where("balance.account_id = ?", balanceQuery.AccountID)
	}

	if balanceQuery.Currency != "" {
		query = query.Where("balance.currency = ?", balanceQuery.Currency)
	}

	if !balanceQuery.From.IsZero() {
		query = query.Where("balance.last_updated_at >= ?", balanceQuery.From)
	}

	if !balanceQuery.To.IsZero() {
		query = query.Where("(balance.created_at <= ?)", balanceQuery.To)
	}

	return query
}

type ListBalancesQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[BalanceQuery]]

func NewListBalancesQuery(opts PaginatedQueryOptions[BalanceQuery]) ListBalancesQuery {
	return ListBalancesQuery{
		Order:    bunpaginate.OrderAsc,
		PageSize: opts.PageSize,
		Options:  opts,
	}
}

func (s *Storage) ListBalances(ctx context.Context, q ListBalancesQuery) (*api.Cursor[models.Balance], error) {
	return PaginateWithOffset[PaginatedQueryOptions[BalanceQuery], models.Balance](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[BalanceQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Order("created_at DESC")

			query = applyBalanceQuery(query, q.Options.Options)

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			}

			return query
		},
	)
}

func (s *Storage) ListBalanceCurrencies(ctx context.Context, accountID models.AccountID) ([]string, error) {
	var currencies []string

	err := s.db.NewSelect().
		ColumnExpr("DISTINCT currency").
		Model(&models.Balance{}).
		Where("account_id = ?", accountID).
		Scan(ctx, &currencies)
	if err != nil {
		return nil, e("failed to list balance currencies", err)
	}

	return currencies, nil
}

func (s *Storage) GetBalanceAtByCurrency(ctx context.Context, accountID models.AccountID, currency string, at time.Time) (*models.Balance, error) {
	var balance models.Balance

	err := s.db.NewSelect().
		Model(&balance).
		Where("account_id = ?", accountID).
		Where("currency = ?", currency).
		Where("created_at <= ?", at).
		Where("last_updated_at >= ?", at).
		Order("last_updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get balance", err)
	}

	return &balance, nil
}

func (s *Storage) GetBalancesAt(ctx context.Context, accountID models.AccountID, at time.Time) ([]*models.Balance, error) {
	currencies, err := s.ListBalanceCurrencies(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to list balance currencies: %w", err)
	}

	var balances []*models.Balance
	for _, currency := range currencies {
		balance, err := s.GetBalanceAtByCurrency(ctx, accountID, currency, at)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			}
			return nil, fmt.Errorf("failed to get balance: %w", err)
		}

		balances = append(balances, balance)
	}

	return balances, nil
}
