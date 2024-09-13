package storage

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type balance struct {
	bun.BaseModel `bun:"table:balances"`

	// Mandatory fields
	AccountID models.AccountID `bun:"account_id,pk,type:character varying,notnull"`
	CreatedAt time.Time        `bun:"created_at,pk,type:timestamp without time zone,notnull"`
	Asset     string           `bun:"asset,pk,type:text,notnull"`

	ConnectorID   models.ConnectorID `bun:"connector_id,type:character varying,notnull"`
	Balance       *big.Int           `bun:"balance,type:numeric,notnull"`
	LastUpdatedAt time.Time          `bun:"last_updated_at,type:timestamp without time zone,notnull"`
}

func (s *store) BalancesUpsert(ctx context.Context, balances []models.Balance) error {
	toInsert := fromBalancesModels(balances)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return e("failed to start transaction", err)
	}
	defer tx.Rollback()

	_, err = tx.NewInsert().
		Model((*models.Balance)(nil)).
		With("cte1", tx.NewValues(&toInsert)).
		Column(
			"account_id",
			"created_at",
			"asset",
			"connector_id",
			"balance",
			"last_updated_at",
		).
		TableExpr(`
		(SELECT *
			FROM cte1
			WHERE cte1.balance != COALESCE((SELECT balance FROM accounts.balances WHERE account_id = cte1.account_id AND last_updated_at < cte1.last_updated_at AND asset = cte1.asset ORDER BY last_updated_at DESC LIMIT 1), cte1.balance+1)
		) data`).
		On("CONFLICT (account_id, created_at, asset) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("failed to create balances", err)
	}

	// Always update the previous row in order to keep the balance history consistent.
	_, err = tx.NewUpdate().
		Model((*models.Balance)(nil)).
		With("cte1", s.db.NewValues(&toInsert)).
		TableExpr(`
					(SELECT (SELECT created_at FROM accounts.balances WHERE last_updated_at < cte1.last_updated_at AND account_id = cte1.account_id AND asset = cte1.asset ORDER BY last_updated_at DESC LIMIT 1), cte1.account_id, cte1.asset, cte1.last_updated_at FROM cte1) data
				`).
		Set("last_updated_at = data.last_updated_at").
		Where("balance.account_id = data.account_id AND balance.asset = data.asset AND balance.created_at = data.created_at").
		Exec(ctx)
	if err != nil {
		return e("failed to update balances", err)
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *store) BalancesDeleteForConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*balance)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)
	if err != nil {
		return e("delete balances", err)
	}

	return nil
}

type BalanceQuery struct {
	AccountID *models.AccountID
	Asset     string
	From      time.Time
	To        time.Time
}

func NewBalanceQuery() BalanceQuery {
	return BalanceQuery{}
}

func (b BalanceQuery) WithAccountID(accountID *models.AccountID) BalanceQuery {
	b.AccountID = accountID

	return b
}

func (b BalanceQuery) WithAsset(asset string) BalanceQuery {
	b.Asset = asset

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

func applyBalanceQuery(query *bun.SelectQuery, balanceQuery BalanceQuery) *bun.SelectQuery {
	if balanceQuery.AccountID != nil {
		query = query.Where("balance.account_id = ?", balanceQuery.AccountID)
	}

	if balanceQuery.Asset != "" {
		query = query.Where("balance.asset = ?", balanceQuery.Asset)
	}

	if !balanceQuery.From.IsZero() {
		query = query.Where("balance.last_updated_at >= ?", balanceQuery.From)
	}

	if !balanceQuery.To.IsZero() {
		query = query.Where("(balance.created_at <= ?)", balanceQuery.To)
	}

	return query
}

type ListBalancesQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[BalanceQuery]]

func NewListBalancesQuery(opts bunpaginate.PaginatedQueryOptions[BalanceQuery]) ListBalancesQuery {
	return ListBalancesQuery{
		Order:    bunpaginate.OrderAsc,
		PageSize: opts.PageSize,
		Options:  opts,
	}
}

func (s *store) BalancesList(ctx context.Context, q ListBalancesQuery) (*bunpaginate.Cursor[models.Balance], error) {
	cursor, err := paginateWithOffset[bunpaginate.PaginatedQueryOptions[BalanceQuery], balance](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[BalanceQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = applyBalanceQuery(query, q.Options.Options)

			query = query.Order("created_at DESC")

			return query
		},
	)
	if err != nil {
		return nil, e("failed to fetch balances", err)
	}

	balances := toBalancesModels(cursor.Data)

	return &bunpaginate.Cursor[models.Balance]{
		PageSize: cursor.PageSize,
		HasMore:  cursor.HasMore,
		Previous: cursor.Previous,
		Next:     cursor.Next,
		Data:     balances,
	}, nil
}

func (s *store) balancesListAssets(ctx context.Context, accountID models.AccountID) ([]string, error) {
	var assets []string

	err := s.db.NewSelect().
		ColumnExpr("DISTINCT asset").
		Model(&models.Balance{}).
		Where("account_id = ?", accountID).
		Scan(ctx, &assets)
	if err != nil {
		return nil, e("failed to list balance assets", err)
	}

	return assets, nil
}

func (s *store) balancesGetAtByAsset(ctx context.Context, accountID models.AccountID, asset string, at time.Time) (*models.Balance, error) {
	var balance balance

	err := s.db.NewSelect().
		Model(&balance).
		Where("account_id = ?", accountID).
		Where("asset = ?", asset).
		Where("created_at <= ?", at).
		Where("last_updated_at >= ?", at).
		Order("last_updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get balance", err)
	}

	return pointer.For(toBalanceModels(balance)), nil
}

func (s *store) BalancesGetAt(ctx context.Context, accountID models.AccountID, at time.Time) ([]*models.Balance, error) {
	assets, err := s.balancesListAssets(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to list balance assets: %w", err)
	}

	var balances []*models.Balance
	for _, currency := range assets {
		balance, err := s.balancesGetAtByAsset(ctx, accountID, currency, at)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			}
			return nil, fmt.Errorf("failed to get balance: %w", err)
		}

		balances = append(balances, balance)
	}

	return balances, nil
}

func fromBalancesModels(from []models.Balance) []balance {
	var to []balance
	for _, b := range from {
		to = append(to, fromBalanceModels(b))
	}
	return to
}

func fromBalanceModels(from models.Balance) balance {
	return balance{
		AccountID:     from.AccountID,
		CreatedAt:     from.CreatedAt,
		Asset:         from.Asset,
		ConnectorID:   from.AccountID.ConnectorID,
		Balance:       from.Balance,
		LastUpdatedAt: from.LastUpdatedAt,
	}
}

func toBalancesModels(from []balance) []models.Balance {
	var to []models.Balance
	for _, b := range from {
		to = append(to, toBalanceModels(b))
	}
	return to
}

func toBalanceModels(from balance) models.Balance {
	return models.Balance{
		AccountID:     from.AccountID,
		CreatedAt:     from.CreatedAt,
		Asset:         from.Asset,
		Balance:       from.Balance,
		LastUpdatedAt: from.LastUpdatedAt,
	}
}
