package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/uptrace/bun"
)

func (s *Storage) InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error {
	if len(balances) == 0 {
		return nil
	}

	type insertedBalance struct {
		CreatedAt time.Time    `bun:"created_at"`
		AccountID string       `bun:"account_id"`
		Asset     models.Asset `bun:"currency"`
	}
	insertedBalances := make([]insertedBalance, 0)

	query := s.db.NewInsert().
		Model((*models.Balance)(nil)).
		With("cte1", s.db.NewValues(&balances)).
		Column(
			"created_at",
			"account_id",
			"balance",
			"currency",
			"last_updated_at",
		).
		Returning("created_at, account_id, currency")
	if checkIfAccountExists {
		query = query.TableExpr(`
		(SELECT *
			FROM cte1
			WHERE EXISTS (SELECT 1 FROM accounts.account WHERE id = cte1.account_id)
				AND cte1.balance != COALESCE((SELECT balance FROM accounts.balances WHERE account_id = cte1.account_id AND last_updated_at < cte1.last_updated_at AND currency = cte1.currency ORDER BY last_updated_at DESC LIMIT 1), cte1.balance+1)
		) data`)
	} else {
		query = query.TableExpr(`
		(SELECT *
			FROM cte1
			WHERE cte1.balance != COALESCE((SELECT balance FROM accounts.balances WHERE account_id = cte1.account_id AND last_updated_at < cte1.last_updated_at AND currency = cte1.currency ORDER BY last_updated_at DESC LIMIT 1), cte1.balance+1)
		) data`)
	}

	err := query.Scan(ctx, &insertedBalances)
	if err != nil {
		return e("failed to create balances", err)
	}

	if len(insertedBalances) != len(balances) {
		balancesToUpdate := make([]*models.Balance, 0)
		for i := range balances {
			found := false
			for _, insertedBalance := range insertedBalances {
				if balances[i].AccountID.String() == insertedBalance.AccountID &&
					balances[i].Asset == insertedBalance.Asset {
					found = true
					break
				}
			}
			if !found {
				balancesToUpdate = append(balancesToUpdate, balances[i])
			}
		}

		if len(balancesToUpdate) > 0 {
			_, err := s.db.NewUpdate().
				Model((*models.Balance)(nil)).
				With("cte1", s.db.NewValues(&balancesToUpdate)).
				TableExpr(`
					(SELECT (SELECT created_at FROM accounts.balances WHERE last_updated_at < cte1.last_updated_at AND account_id = cte1.account_id AND currency = cte1.currency ORDER BY last_updated_at DESC LIMIT 1), cte1.account_id, cte1.currency, cte1.last_updated_at FROM cte1) data
				`).
				Set("last_updated_at = data.last_updated_at").
				Where("balance.account_id = data.account_id AND balance.currency = data.currency AND balance.created_at = data.created_at").
				Exec(ctx)
			if err != nil {
				return e("failed to update balances", err)
			}
		}
	}

	return nil
}

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
	Pagination Paginator
}

func NewBalanceQuery(pagination Paginator) BalanceQuery {
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
