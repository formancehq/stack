package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/uptrace/bun"
)

func (s *Storage) InsertBalances(ctx context.Context, balances []*models.Balance, checkIfAccountExists bool) error {
	if len(balances) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			logging.FromContext(ctx).Error("failed to rollback transaction", err)
		}
	}()

	for _, balance := range balances {
		if err := s.insertBalances(ctx, tx, balance, checkIfAccountExists); err != nil {
			return err
		}
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *Storage) insertBalances(ctx context.Context, tx bun.Tx, balance *models.Balance, checkIfAccountExists bool) error {
	var account models.Account
	err := tx.NewSelect().
		Model(&account).
		Where("id = ?", balance.AccountID).
		Scan(ctx)
	if err != nil {
		pErr := e("failed to get account", err)
		if !errors.Is(pErr, ErrNotFound) {
			return pErr
		}

		// error not found here
		if !checkIfAccountExists {
			// return an error here to keep the same behavior as before when
			// checkIfAccountExists is false
			return pErr
		} else {
			// if checkIfAccountExists is true, we should ignore the balance
			// if the account does not exist
			return nil
		}
	}

	var lastBalance models.Balance
	found := true
	err = tx.NewSelect().
		Model(&lastBalance).
		Where("account_id = ? AND currency = ?", balance.AccountID, balance.Asset).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		pErr := e("failed to get account", err)
		if !errors.Is(pErr, ErrNotFound) {
			return pErr
		}
		found = false
	}

	if found && lastBalance.CreatedAt.After(balance.CreatedAt) {
		// Do not insert balance if the last balance is newer
		return nil
	}

	switch {
	case found && lastBalance.Balance.Cmp(balance.Balance) == 0:
		// same balance, no need to have a new entry, just update the last one
		_, err = tx.NewUpdate().
			Model((*models.Balance)(nil)).
			Set("last_updated_at = ?", balance.LastUpdatedAt).
			Where("account_id = ? AND created_at = ? AND currency = ?", lastBalance.AccountID, lastBalance.CreatedAt, lastBalance.Asset).
			Exec(ctx)
		if err != nil {
			return e("failed to update balance", err)
		}

	case found && lastBalance.Balance.Cmp(balance.Balance) != 0:
		// different balance, insert a new entry
		_, err = tx.NewInsert().
			Model(balance).
			Exec(ctx)
		if err != nil {
			return e("failed to insert balance", err)
		}

		// and update last row last updated at to this created at
		_, err = tx.NewUpdate().
			Model(&lastBalance).
			Set("last_updated_at = ?", balance.CreatedAt).
			Where("account_id = ? AND created_at = ? AND currency = ?", lastBalance.AccountID, lastBalance.CreatedAt, lastBalance.Asset).
			Exec(ctx)
		if err != nil {
			return e("failed to update balance", err)
		}

	case !found:
		// no balance found, insert a new entry
		_, err = tx.NewInsert().
			Model(balance).
			Exec(ctx)
		if err != nil {
			return e("failed to insert balance", err)
		}
	}

	return nil
}

func (s *Storage) GetBalancesForAccountID(ctx context.Context, accountID models.AccountID) ([]*models.Balance, error) {
	var balances []*models.Balance

	err := s.db.NewSelect().
		Model(&balances).
		Where("account_id = ?", accountID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get balances", err)
	}

	return balances, nil
}
