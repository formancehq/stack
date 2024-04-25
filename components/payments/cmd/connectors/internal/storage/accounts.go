package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) UpsertAccounts(ctx context.Context, accounts []*models.Account) ([]models.AccountID, error) {
	if len(accounts) == 0 {
		return nil, nil
	}

	var idsUpdated []string
	err := s.db.NewUpdate().
		With("_data",
			s.db.NewValues(&accounts).
				Column(
					"id",
					"default_currency",
					"account_name",
					"metadata",
				),
		).
		Model((*models.Account)(nil)).
		TableExpr("_data").
		Set("default_currency = _data.default_currency").
		Set("account_name = _data.account_name").
		Set("metadata = _data.metadata").
		Where(`(account.id = _data.id) AND
			(account.default_currency != _data.default_currency OR account.account_name != _data.account_name OR (account.metadata != _data.metadata))`).
		Returning("account.id").
		Scan(ctx, &idsUpdated)
	if err != nil {
		return nil, e("failed to update accounts", err)
	}

	idsUpdatedMap := make(map[string]struct{})
	for _, id := range idsUpdated {
		idsUpdatedMap[id] = struct{}{}
	}

	accountsToInsert := make([]*models.Account, 0, len(accounts))
	for _, account := range accounts {
		if _, ok := idsUpdatedMap[account.ID.String()]; !ok {
			accountsToInsert = append(accountsToInsert, account)
		}
	}

	var idsInserted []string
	if len(accountsToInsert) > 0 {
		err = s.db.NewInsert().
			Model(&accountsToInsert).
			On("CONFLICT (id) DO NOTHING").
			Returning("account.id").
			Scan(ctx, &idsInserted)
		if err != nil {
			return nil, e("failed to create accounts", err)
		}
	}

	res := make([]models.AccountID, 0, len(idsUpdated)+len(idsInserted))
	for _, id := range idsUpdated {
		res = append(res, models.MustAccountIDFromString(id))
	}
	for _, id := range idsInserted {
		res = append(res, models.MustAccountIDFromString(id))
	}

	return res, nil
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
