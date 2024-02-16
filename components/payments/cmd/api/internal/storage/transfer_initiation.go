package storage

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

func (s *Storage) GetTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error) {
	var transferInitiation models.TransferInitiation

	query := s.db.NewSelect().
		Column("id", "connector_id", "created_at", "scheduled_at", "description", "type", "source_account_id", "destination_account_id", "provider", "initial_amount", "amount", "asset", "metadata").
		Model(&transferInitiation).
		Relation("RelatedAdjustments").
		Where("id = ?", id)

	err := query.Scan(ctx)
	if err != nil {
		return nil, e("failed to get transfer initiation", err)
	}

	transferInitiation.SortRelatedAdjustments()

	transferInitiation.RelatedPayments, err = s.ReadTransferInitiationPayments(ctx, id)
	if err != nil {
		return nil, e("failed to get transfer initiation payments", err)
	}

	return &transferInitiation, nil
}

func (s *Storage) ReadTransferInitiationPayments(ctx context.Context, id models.TransferInitiationID) ([]*models.TransferInitiationPayment, error) {
	var payments []*models.TransferInitiationPayment

	query := s.db.NewSelect().
		Column("transfer_initiation_id", "payment_id", "created_at", "status", "error").
		Model(&payments).
		Where("transfer_initiation_id = ?", id).
		Order("created_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, e("failed to get transfer initiation payments", err)
	}

	return payments, nil
}

type TransferInitiationQuery struct{}

type ListTransferInitiationsQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[TransferInitiationQuery]]

func NewListTransferInitiationsQuery(opts PaginatedQueryOptions[TransferInitiationQuery]) ListTransferInitiationsQuery {
	return ListTransferInitiationsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *Storage) ListTransferInitiations(ctx context.Context, q ListTransferInitiationsQuery) (*api.Cursor[models.TransferInitiation], error) {
	return PaginateWithOffset[PaginatedQueryOptions[TransferInitiationQuery], models.TransferInitiation](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[TransferInitiationQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Column("id", "connector_id", "created_at", "scheduled_at", "description", "type", "source_account_id", "destination_account_id", "provider", "initial_amount", "amount", "asset", "metadata").
				Relation("RelatedAdjustments").
				Order("created_at DESC")

			if q.Options.QueryBuilder != nil {
				where, args, err := s.transferInitiationQueryContext(q.Options.QueryBuilder)
				if err != nil {
					// TODO: handle error
					panic(err)
				}
				query = query.Where(where, args...)
			}

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			}

			return query
		},
	)
}

func (s *Storage) transferInitiationQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "source_account_id", key == "destination_account_id":
			if operator != "$match" {
				return "", nil, fmt.Errorf("'%s' columns can only be used with $match", key)
			}

			switch accountID := value.(type) {
			case string:
				return fmt.Sprintf("%s = ?", key), []any{accountID}, nil
			default:
				return "", nil, fmt.Errorf("unexpected type %T for column '%s'", accountID, key)
			}
		default:
			return "", nil, fmt.Errorf("unknown key '%s' when building query", key)
		}
	}))
}
