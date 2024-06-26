package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
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

func (s *Storage) paymentsQueryContext(qb query.Builder) (map[string]string, string, []any, error) {
	metadata := make(map[string]string)

	where, args, err := qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "reference":
			return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil

		case key == "type",
			key == "status",
			key == "asset":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'type' column can only be used with $match")
			}
			return fmt.Sprintf("%s = ?", key), []any{value}, nil

		case key == "connectorID":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'type' column can only be used with $match")
			}
			return "connector_id = ?", []any{value}, nil

		case key == "amount":
			return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil

		case key == "initialAmount":
			return fmt.Sprintf("initial_amount %s ?", query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil

		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'metadata' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			valueString, ok := value.(string)
			if !ok {
				return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("metadata value must be a string, got %T", value))
			}

			metadata[match[0][1]] = valueString

			// Do nothing here, as we don't want to add this to the query
			return "", nil, nil

		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))

	return metadata, where, args, err
}

func (s *Storage) ListPayments(ctx context.Context, q ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error) {
	var (
		metadata map[string]string
		where    string
		args     []any
		err      error
	)
	if q.Options.QueryBuilder != nil {
		metadata, where, args, err = s.paymentsQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	return PaginateWithOffset[PaginatedQueryOptions[PaymentQuery], models.Payment](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PaymentQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Relation("Metadata").
				Relation("Connector").
				Relation("Adjustments")

			if where != "" {
				query = query.Where(where, args...)
			}

			if len(metadata) > 0 {
				metadataQuery := s.db.NewSelect().Model((*models.PaymentMetadata)(nil))
				for key, value := range metadata {
					metadataQuery = metadataQuery.Where("payment_metadata.key = ? AND payment_metadata.value = ?", key, value)
				}
				query = query.With("_metadata", metadataQuery)
				query = query.Where("payment.id IN (SELECT payment_id FROM _metadata)")
			}

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			} else {
				query = query.Order("created_at DESC")
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
			Where("metadata.value != EXCLUDED.value").
			Exec(ctx)
		if err != nil {
			return e("failed to create metadata", err)
		}
	}

	return nil
}
