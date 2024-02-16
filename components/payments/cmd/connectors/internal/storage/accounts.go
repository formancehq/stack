package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

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

func (s *Storage) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account

	err := s.db.NewSelect().
		Model(&account).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get account", err)
	}

	return &account, nil
}
