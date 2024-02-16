package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error {
	if len(balances) == 0 {
		return nil
	}

	query := s.db.NewInsert().
		Model((*models.Balance)(nil)).
		With("cte1", s.db.NewValues(&balances)).
		Column(
			"created_at",
			"account_id",
			"balance",
			"currency",
			"last_updated_at",
		)
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

	_, err := query.Exec(ctx)
	if err != nil {
		return e("failed to create balances", err)
	}

	// Always update the previous row in order to keep the balance history consistent.
	_, err = s.db.NewUpdate().
		Model((*models.Balance)(nil)).
		With("cte1", s.db.NewValues(&balances)).
		TableExpr(`
					(SELECT (SELECT created_at FROM accounts.balances WHERE last_updated_at < cte1.last_updated_at AND account_id = cte1.account_id AND currency = cte1.currency ORDER BY last_updated_at DESC LIMIT 1), cte1.account_id, cte1.currency, cte1.last_updated_at FROM cte1) data
				`).
		Set("last_updated_at = data.last_updated_at").
		Where("balance.account_id = data.account_id AND balance.currency = data.currency AND balance.created_at = data.created_at").
		Exec(ctx)
	if err != nil {
		return e("failed to update balances", err)
	}

	return nil
}

func (s *Storage) GetBalancesForAccountID(ctx context.Context, accountID models.AccountID) ([]*models.Balance, error) {
	var balances []*models.Balance

	err := s.db.NewSelect().
		Model(&balances).
		Where("account_id = ?", accountID).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get balances", err)
	}

	return balances, nil
}
