package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type connector struct {
	bun.BaseModel `bun:"table:connectors"`

	// Mandatory fields
	ID        models.ConnectorID `bun:"id,pk,type:character varying,notnull"`
	Name      string             `bun:"name,type:text,notnull"`
	CreatedAt time.Time          `bun:"created_at,type:timestamp without time zone,notnull"`
	Provider  string             `bun:"provider,type:text,notnull"`

	// EncryptedConfig is a PGP-encrypted JSON string.
	EncryptedConfig string `bun:"config,type:bytea,notnull"`

	// Config is a decrypted config. It is not stored in the database.
	DecryptedConfig json.RawMessage `bun:"decrypted_config,scanonly"`
}

func (s *store) ConnectorsInstall(ctx context.Context, c models.Connector) error {
	toInsert := connector{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		Provider:  c.Provider,
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}
	defer tx.Rollback()

	_, err = tx.NewInsert().
		Model(&toInsert).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("failed to insert connector", err)
	}

	_, err = tx.NewUpdate().
		Model((*connector)(nil)).
		Set("config = pgp_sym_encrypt(?::TEXT, ?, ?)", c.Config, s.configEncryptionKey, encryptionOptions).
		Where("id = ?", toInsert.ID).
		Exec(ctx)
	if err != nil {
		return e("failed to encrypt config", err)
	}

	return e("failed to commit transaction", tx.Commit())
}

// TODO(polo): find a better way to delete all data
func (s *store) ConnectorsUninstall(ctx context.Context, id models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*connector)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return e("failed to delete connector", err)
}

func (s *store) ConnectorsGet(ctx context.Context, id models.ConnectorID) (*models.Connector, error) {
	var connector connector

	err := s.db.NewSelect().
		Model(&connector).
		ColumnExpr("*, pgp_sym_decrypt(config, ?, ?) AS decrypted_config", s.configEncryptionKey, encryptionOptions).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to fetch connector", err)
	}

	return &models.Connector{
		ID:        connector.ID,
		Name:      connector.Name,
		CreatedAt: connector.CreatedAt,
		Provider:  connector.Provider,
		Config:    connector.DecryptedConfig,
	}, nil
}

type ConnectorQuery struct{}

type ListConnectorsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[ConnectorQuery]]

func NewListConnectorsQuery(opts bunpaginate.PaginatedQueryOptions[ConnectorQuery]) ListConnectorsQuery {
	return ListConnectorsQuery{
		Order:    bunpaginate.OrderAsc,
		PageSize: opts.PageSize,
		Options:  opts,
	}
}

func (s *store) connectorsQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "name",
			key == "provider":
			return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil
		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))
}

func (s *store) ConnectorsList(ctx context.Context, q ListConnectorsQuery) (*bunpaginate.Cursor[models.Connector], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.connectorsQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	cursor, err := paginateWithOffset[bunpaginate.PaginatedQueryOptions[ConnectorQuery], connector](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[ConnectorQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			if where != "" {
				query = query.Where(where, args...)
			}

			query = query.ColumnExpr("*, pgp_sym_decrypt(config, ?, ?) AS decrypted_config", s.configEncryptionKey, encryptionOptions)

			// TODO(polo): sorter ?
			query = query.Order("created_at DESC")

			return query
		},
	)
	if err != nil {
		return nil, e("failed to fetch connectors", err)
	}

	connectors := make([]models.Connector, 0, len(cursor.Data))
	for _, c := range cursor.Data {
		connectors = append(connectors, models.Connector{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
			Provider:  c.Provider,
			Config:    c.DecryptedConfig,
		})
	}

	return &bunpaginate.Cursor[models.Connector]{
		PageSize: cursor.PageSize,
		HasMore:  cursor.HasMore,
		Previous: cursor.Previous,
		Next:     cursor.Next,
		Data:     connectors,
	}, nil
}
