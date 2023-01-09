package storage

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/app/models"
)

func (s *Storage) UpsertAccounts(ctx context.Context, provider models.ConnectorProvider, accounts []models.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	accountsMap := make(map[string]models.Account)
	for _, account := range accounts {
		accountsMap[account.Reference] = account
	}

	accounts = make([]models.Account, 0, len(accountsMap))
	for _, account := range accountsMap {
		accounts = append(accounts, account)
	}

	_, err := s.db.NewInsert().
		Model(&accounts).
		On("CONFLICT (reference) DO UPDATE").
		Set("provider = EXCLUDED.provider").
		Set("type = EXCLUDED.type").
		Exec(ctx)
	if err != nil {
		return e("failed to create accounts", err)
	}

	return nil
}

func (s *Storage) ListAccounts(ctx context.Context, sort Sorter, pagination Paginator) ([]*models.Account, error) {
	var accounts []*models.Account

	query := s.db.NewSelect().
		Model(&accounts)

	if sort != nil {
		query = sort.apply(query)
	}

	query = pagination.apply(query)

	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	return accounts, nil
}
