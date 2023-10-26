package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/uptrace/bun"
)

func (s *Storage) ListBalances(ctx context.Context, balanceQuery BalanceQuery) ([]*models.Balance, PaginationDetails, error) {
	var balances []*models.Balance

	query := s.db.NewSelect().
		Model(&balances)

	query = balanceQuery.Pagination.apply(query, "balance.created_at")

	query = applyBalanceQuery(query, balanceQuery)

	err := query.Scan(ctx)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to list balances", err)
	}

	var (
		hasMore                       = len(balances) > balanceQuery.Pagination.pageSize
		hasPrevious                   bool
		firstReference, lastReference string
	)

	if hasMore {
		if balanceQuery.Pagination.cursor.Next || balanceQuery.Pagination.cursor.Reference == "" {
			balances = balances[:balanceQuery.Pagination.pageSize]
		} else {
			balances = balances[1:]
		}
	}

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].CreatedAt.After(balances[j].CreatedAt)
	})

	if len(balances) > 0 {
		firstReference = balances[0].CreatedAt.Format(time.RFC3339Nano)
		lastReference = balances[len(balances)-1].CreatedAt.Format(time.RFC3339Nano)

		query = s.db.NewSelect().Model(&balances)
		query = applyBalanceQuery(query, balanceQuery)

		hasPrevious, err = balanceQuery.Pagination.hasPrevious(ctx, query, "created_at", firstReference)
		if err != nil {
			return nil, PaginationDetails{}, fmt.Errorf("failed to check if there is a previous page: %w", err)
		}
	}

	paginationDetails, err := balanceQuery.Pagination.paginationDetails(hasMore, hasPrevious, firstReference, lastReference)
	if err != nil {
		return nil, PaginationDetails{}, fmt.Errorf("failed to get pagination details: %w", err)
	}

	return balances, paginationDetails, nil
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
		query = query.Where("balance.last_updated_at < ?", balanceQuery.To)
	}

	if balanceQuery.Limit > 0 {
		query = query.Limit(balanceQuery.Limit)
	}

	return query
}

type BalanceQuery struct {
	AccountID  *models.AccountID
	Currency   string
	From       time.Time
	To         time.Time
	Limit      int
	Pagination PaginatorQuery
}

func NewBalanceQuery(pagination PaginatorQuery) BalanceQuery {
	return BalanceQuery{
		Pagination: pagination,
	}
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

func (b BalanceQuery) WithLimit(limit int) BalanceQuery {
	b.Limit = limit

	return b
}
