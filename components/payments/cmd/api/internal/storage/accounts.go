package storage

import (
	"context"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) ListAccounts(ctx context.Context, pagination PaginatorQuery) ([]*models.Account, PaginationDetails, error) {
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
