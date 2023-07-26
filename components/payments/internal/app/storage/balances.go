package storage

import (
	"context"
	"time"

	"github.com/formancehq/payments/internal/app/models"
)

func (s *Storage) InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error {
	if len(balances) == 0 {
		return nil
	}

	type insertedBalance struct {
		CreatedAt time.Time `bun:"created_at"`
		AccountID string    `bun:"account_id"`
		Currency  string    `bun:"currency"`
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
		for _, balance := range balances {
			found := false
			for _, insertedBalance := range insertedBalances {
				if balance.AccountID.String() == insertedBalance.AccountID &&
					balance.Currency == insertedBalance.Currency &&
					balance.CreatedAt == insertedBalance.CreatedAt {
					found = true
					break
				}
			}
			if !found {
				balancesToUpdate = append(balancesToUpdate, balance)
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

func (s *Storage) ListBalances(ctx context.Context, query BalanceQuery) ([]*models.Balance, PaginationDetails, error) {
	// TODO(polo)
	return nil, PaginationDetails{}, nil
}

type BalanceQuery struct {
	AccountID  *models.AccountID
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
