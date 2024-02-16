package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/uptrace/bun"
)

type PaymentQuery struct{}

type ListPaymentsQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PaymentQuery]]

func NewListPaymentsQuery(opts PaginatedQueryOptions[PaymentQuery]) ListPaymentsQuery {
	return ListPaymentsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *Storage) ListPayments(ctx context.Context, q ListPaymentsQuery) (*api.Cursor[models.Payment], error) {
	return PaginateWithOffset[PaginatedQueryOptions[PaymentQuery], models.Payment](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PaymentQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Relation("Connector").
				Relation("Metadata").
				Relation("Adjustments").
				Order("created_at DESC")

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			}

			return query
		},
	)
}

func (s *Storage) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	var payment models.Payment

	err := s.db.NewSelect().
		Model(&payment).
		Relation("Connector").
		Relation("Metadata").
		Relation("Adjustments").
		Where("payment.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e(fmt.Sprintf("failed to get payment %s", id), err)
	}

	return &payment, nil
}

func (s *Storage) UpsertPayments(ctx context.Context, payments []*models.Payment) error {
	if len(payments) == 0 {
		return nil
	}

	_, err := s.db.NewInsert().
		Model(&payments).
		On("CONFLICT (reference) DO UPDATE").
		Set("amount = EXCLUDED.amount").
		Set("type = EXCLUDED.type").
		Set("status = EXCLUDED.status").
		Set("raw_data = EXCLUDED.raw_data").
		Set("scheme = EXCLUDED.scheme").
		Set("asset = EXCLUDED.asset").
		Set("source_account_id = EXCLUDED.source_account_id").
		Set("destination_account_id = EXCLUDED.destination_account_id").
		Exec(ctx)
	if err != nil {
		return e("failed to create payments", err)
	}

	var adjustments []*models.PaymentAdjustment
	var metadata []*models.PaymentMetadata

	for i := range payments {
		for _, adjustment := range payments[i].Adjustments {
			if adjustment.Reference == "" {
				continue
			}

			adjustment.PaymentID = payments[i].ID

			adjustments = append(adjustments, adjustment)
		}

		for _, data := range payments[i].Metadata {
			data.PaymentID = payments[i].ID
			data.Changelog = append(data.Changelog,
				models.MetadataChangelog{
					CreatedAt: time.Now(),
					Value:     data.Value,
				})

			metadata = append(metadata, data)
		}
	}

	if len(adjustments) > 0 {
		_, err = s.db.NewInsert().
			Model(&adjustments).
			On("CONFLICT (reference) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return e("failed to create adjustments", err)
		}
	}

	if len(metadata) > 0 {
		_, err = s.db.NewInsert().
			Model(&metadata).
			On("CONFLICT (payment_id, key) DO UPDATE").
			Set("value = EXCLUDED.value").
			Set("changelog = metadata.changelog || EXCLUDED.changelog").
			Exec(ctx)
		if err != nil {
			return e("failed to create metadata", err)
		}
	}

	return nil
}
