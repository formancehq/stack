package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/query"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type payment struct {
	bun.BaseModel `bun:"table:payments"`

	// Mandatory fields
	ID            models.PaymentID     `bun:"id,pk,type:character varying,notnull"`
	ConnectorID   models.ConnectorID   `bun:"connector_id,type:character varying,notnull"`
	Reference     string               `bun:"reference,type:text,notnull"`
	CreatedAt     time.Time            `bun:"created_at,type:timestamp without time zone,notnull"`
	Type          models.PaymentType   `bun:"type,type:text,notnull"`
	InitialAmount *big.Int             `bun:"initial_amount,type:numeric,notnull"`
	Amount        *big.Int             `bun:"amount,type:numeric,notnull"`
	Asset         string               `bun:"asset,type:text,notnull"`
	Scheme        models.PaymentScheme `bun:"scheme,type:text,notnull"`

	// Scan only fields
	Status models.PaymentStatus `bun:"status,type:text,notnull,scanonly"`

	// Optional fields
	// c.f.: https://bun.uptrace.dev/guide/models.html#nulls
	SourceAccountID      *models.AccountID `bun:"source_account_id,type:character varying,nullzero"`
	DestinationAccountID *models.AccountID `bun:"destination_account_id,type:character varying,nullzero"`

	// Optional fields with default
	// c.f. https://bun.uptrace.dev/guide/models.html#default
	Metadata map[string]string `bun:"metadata,type:jsonb,nullzero,notnull,default:'{}'"`
}

type paymentAdjustment struct {
	bun.BaseModel `bun:"table:payment_adjustments"`

	// Mandatory fields
	ID        models.PaymentAdjustmentID `bun:"id,pk,type:character varying,notnull"`
	PaymentID models.PaymentID           `bun:"payment_id,type:character varying,notnull"`
	CreatedAt time.Time                  `bun:"created_at,type:timestamp without time zone,notnull"`
	Status    models.PaymentStatus       `bun:"status,type:text,notnull"`
	Raw       json.RawMessage            `bun:"raw,type:json,notnull"`

	// Optional fields
	// c.f.: https://bun.uptrace.dev/guide/models.html#nulls
	Amount *big.Int `bun:"amount,type:numeric,nullzero"`
	Asset  *string  `bun:"asset,type:text,nullzero"`

	// Optional fields with default
	// c.f. https://bun.uptrace.dev/guide/models.html#default
	Metadata map[string]string `bun:"metadata,type:jsonb,nullzero,notnull,default:'{}'"`
}

func (s *store) PaymentsUpsert(ctx context.Context, payments []models.Payment) error {
	paymentsToInsert := make([]payment, 0, len(payments))
	adjustmentsToInsert := make([]paymentAdjustment, 0)
	paymentsRefunded := make([]payment, 0)
	for _, p := range payments {
		paymentsToInsert = append(paymentsToInsert, fromPaymentModels(p))

		for _, a := range p.Adjustments {
			adjustmentsToInsert = append(adjustmentsToInsert, fromPaymentAdjustmentModels(a))
			switch a.Status {
			case models.PAYMENT_STATUS_REFUNDED:
				res := fromPaymentModels(p)
				res.Amount = a.Amount
				paymentsRefunded = append(paymentsRefunded, res)
			}
		}
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to create transaction")
	}
	defer tx.Rollback()

	if len(paymentsToInsert) > 0 {
		_, err = tx.NewInsert().
			Model(&paymentsToInsert).
			On("CONFLICT (id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return e("failed to insert payments", err)
		}
	}

	if len(paymentsRefunded) > 0 {
		_, err = tx.NewInsert().
			Model(&paymentsRefunded).
			On("CONFLICT (id) DO UPDATE").
			Set("amount = payment.amount - EXCLUDED.amount").
			Exec(ctx)
		if err != nil {
			return e("failed to update payment", err)
		}
	}

	if len(adjustmentsToInsert) > 0 {
		_, err = tx.NewInsert().
			Model(&adjustmentsToInsert).
			On("CONFLICT (id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return e("failed to insert adjustments", err)
		}
	}

	return e("failed to commit transactions", tx.Commit())
}

func (s *store) PaymentsUpdateMetadata(ctx context.Context, id models.PaymentID, metadata map[string]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return e("update payment metadata", err)
	}
	defer tx.Rollback()

	var payment payment
	err = tx.NewSelect().
		Model(&payment).
		Column("id", "metadata").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return e("update payment metadata", err)
	}

	if payment.Metadata == nil {
		payment.Metadata = make(map[string]string)
	}

	for k, v := range metadata {
		payment.Metadata[k] = v
	}

	_, err = tx.NewUpdate().
		Model(&payment).
		Column("metadata").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("update payment metadata", err)
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *store) PaymentsGet(ctx context.Context, id models.PaymentID) (*models.Payment, error) {
	var payment payment

	err := s.db.NewSelect().
		Model(&payment).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get payment", err)
	}

	var ajs []paymentAdjustment
	err = s.db.NewSelect().
		Model(&ajs).
		Where("payment_id = ?", id).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get payment adjustments", err)
	}

	adjustments := make([]models.PaymentAdjustment, 0, len(ajs))
	for _, a := range ajs {
		adjustments = append(adjustments, toPaymentAdjustmentModels(a))
	}

	status := models.PAYMENT_STATUS_PENDING
	if len(adjustments) > 0 {
		status = adjustments[len(adjustments)-1].Status
	}
	res := toPaymentModels(payment, status)
	res.Adjustments = adjustments
	return &res, nil
}

func (s *store) PaymentsDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*payment)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)

	return e("failed to delete payments", err)
}

type PaymentQuery struct{}

type ListPaymentsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[PaymentQuery]]

func NewListPaymentsQuery(opts bunpaginate.PaginatedQueryOptions[PaymentQuery]) ListPaymentsQuery {
	return ListPaymentsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *store) paymentsQueryContext(qb query.Builder) (string, []any, error) {
	where, args, err := qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "reference",
			key == "connector_id",
			key == "type",
			key == "asset",
			key == "scheme",
			key == "status",
			key == "source_account_id",
			key == "destination_account_id":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'type' column can only be used with $match")
			}
			return fmt.Sprintf("%s = ?", key), []any{value}, nil

		case key == "initial_amount",
			key == "amount":
			return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil
		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'metadata' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			key := "metadata"
			return key + " @> ?", []any{map[string]any{
				match[0][1]: value,
			}}, nil
		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))

	return where, args, err
}

func (s *store) PaymentsList(ctx context.Context, q ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.paymentsQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	// TODO(polo): should fetch the adjustments and get the last status and amount?
	cursor, err := paginateWithOffset[bunpaginate.PaginatedQueryOptions[PaymentQuery], payment](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[PaymentQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			if where != "" {
				query = query.Where(where, args...)
			}

			query.Column("payment.*", "apd.status").
				Join(`join lateral (
				select status
				from payment_adjustments apd
				where payment_id = payment.id
				order by created_at desc
				limit 1
			) apd on true`)

			// TODO(polo): sorter ?
			query = query.Order("created_at DESC")

			return query
		},
	)
	if err != nil {
		return nil, e("failed to fetch payments", err)
	}

	payments := make([]models.Payment, 0, len(cursor.Data))
	for _, p := range cursor.Data {
		payments = append(payments, toPaymentModels(p, p.Status))
	}

	return &bunpaginate.Cursor[models.Payment]{
		PageSize: cursor.PageSize,
		HasMore:  cursor.HasMore,
		Previous: cursor.Previous,
		Next:     cursor.Next,
		Data:     payments,
	}, nil
}

func fromPaymentModels(from models.Payment) payment {
	return payment{
		ID:                   from.ID,
		ConnectorID:          from.ConnectorID,
		Reference:            from.Reference,
		CreatedAt:            from.CreatedAt.UTC(),
		Type:                 from.Type,
		InitialAmount:        from.InitialAmount,
		Amount:               from.Amount,
		Asset:                from.Asset,
		Scheme:               from.Scheme,
		SourceAccountID:      from.SourceAccountID,
		DestinationAccountID: from.DestinationAccountID,
		Metadata:             from.Metadata,
	}
}

func toPaymentModels(payment payment, status models.PaymentStatus) models.Payment {
	return models.Payment{
		ID:                   payment.ID,
		ConnectorID:          payment.ConnectorID,
		InitialAmount:        payment.InitialAmount,
		Reference:            payment.Reference,
		CreatedAt:            payment.CreatedAt.UTC(),
		Type:                 payment.Type,
		Amount:               payment.Amount,
		Asset:                payment.Asset,
		Scheme:               payment.Scheme,
		Status:               status,
		SourceAccountID:      payment.SourceAccountID,
		DestinationAccountID: payment.DestinationAccountID,
		Metadata:             payment.Metadata,
	}
}

func fromPaymentAdjustmentModels(from models.PaymentAdjustment) paymentAdjustment {
	return paymentAdjustment{
		ID:        from.ID,
		PaymentID: from.PaymentID,
		CreatedAt: from.CreatedAt.UTC(),
		Status:    from.Status,
		Amount:    from.Amount,
		Asset:     from.Asset,
		Metadata:  from.Metadata,
		Raw:       from.Raw,
	}
}

func toPaymentAdjustmentModels(from paymentAdjustment) models.PaymentAdjustment {
	return models.PaymentAdjustment{
		ID:        from.ID,
		PaymentID: from.PaymentID,
		CreatedAt: from.CreatedAt.UTC(),
		Status:    from.Status,
		Amount:    from.Amount,
		Asset:     from.Asset,
		Metadata:  from.Metadata,
		Raw:       from.Raw,
	}
}
