package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/app/models"
)

func (s *Storage) UpsertAccounts(ctx context.Context, provider models.ConnectorProvider, accounts []*models.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	accountsMap := make(map[string]*models.Account)
	for _, account := range accounts {
		accountsMap[account.Reference] = account
	}

	connector, err := s.GetConnector(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to get connector: %w", err)
	}

	accounts = make([]*models.Account, 0, len(accountsMap))
	for _, account := range accountsMap {
		account.ConnectorID = connector.ID
		accounts = append(accounts, account)
	}

	_, err = s.db.NewInsert().
		Model(&accounts).
		On("CONFLICT (id) DO UPDATE").
		Set("connector_id = EXCLUDED.connector_id").
		Set("provider = EXCLUDED.provider").
		Set("raw_data = EXCLUDED.raw_data").
		Set("default_currency = EXCLUDED.default_currency").
		Set("account_name = EXCLUDED.account_name").
		Exec(ctx)
	if err != nil {
		return e("failed to create accounts", err)
	}

	return nil
}

func (s *Storage) ListAccounts(ctx context.Context, pagination Paginator) ([]*models.Account, PaginationDetails, error) {
	var accounts []*models.Account

	query := s.db.NewSelect().
		Model(&accounts)

	query = pagination.apply(query, "account.created_at")

	err := query.Scan(ctx)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to list payments", err)
	}

	var (
		hasMore                       = len(accounts) > pagination.pageSize
		hasPrevious                   bool
		firstReference, lastReference string
	)

	if hasMore {
		if pagination.cursor.Next || pagination.cursor.Reference == "" {
			accounts = accounts[:pagination.pageSize]
		} else {
			accounts = accounts[1:]
		}
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].CreatedAt.After(accounts[j].CreatedAt)
	})

	if len(accounts) > 0 {
		firstReference = accounts[0].CreatedAt.Format(time.RFC3339Nano)
		lastReference = accounts[len(accounts)-1].CreatedAt.Format(time.RFC3339Nano)

		query = s.db.NewSelect().Model(&accounts)

		hasPrevious, err = pagination.hasPrevious(ctx, query, "account.created_at", firstReference)
		if err != nil {
			return nil, PaginationDetails{}, e("failed to check if there is a previous page", err)
		}
	}

	paginationDetails, err := pagination.paginationDetails(hasMore, hasPrevious, firstReference, lastReference)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to get pagination details", err)
	}

	return accounts, paginationDetails, nil
}
