package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) ListPayments(ctx context.Context, pagination PaginatorQuery) ([]*models.Payment, PaginationDetails, error) {
	var payments []*models.Payment

	query := s.db.NewSelect().
		Model(&payments).
		Relation("Connector").
		Relation("Metadata").
		Relation("Adjustments")

	query = pagination.apply(query, "payment.created_at")

	err := query.Scan(ctx)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to list payments", err)
	}

	var (
		hasMore                       = len(payments) > pagination.pageSize
		hasPrevious                   bool
		firstReference, lastReference string
	)

	if hasMore {
		if pagination.cursor.Next || pagination.cursor.Reference == "" {
			payments = payments[:pagination.pageSize]
		} else {
			payments = payments[1:]
		}
	}

	sort.Slice(payments, func(i, j int) bool {
		return payments[i].CreatedAt.After(payments[j].CreatedAt)
	})

	if len(payments) > 0 {
		firstReference = payments[0].CreatedAt.Format(time.RFC3339Nano)
		lastReference = payments[len(payments)-1].CreatedAt.Format(time.RFC3339Nano)

		query = s.db.NewSelect().Model(&payments)

		hasPrevious, err = pagination.hasPrevious(ctx, query, "payment.created_at", firstReference)
		if err != nil {
			return nil, PaginationDetails{}, fmt.Errorf("failed to check if there is a previous page: %w", err)
		}
	}

	paginationDetails, err := pagination.paginationDetails(hasMore, hasPrevious, firstReference, lastReference)
	if err != nil {
		return nil, PaginationDetails{}, fmt.Errorf("failed to get pagination details: %w", err)
	}

	return payments, paginationDetails, nil
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

func (s *Storage) UpsertPayments(ctx context.Context, payments []*models.Payment) ([]*models.PaymentID, error) {
	if len(payments) == 0 {
		return nil, nil
	}

	var idsUpdated []string
	err := s.db.NewUpdate().
		With("_data",
			s.db.NewValues(&payments).
				Column(
					"id",
					"amount",
					"type",
					"scheme",
					"asset",
					"source_account_id",
					"destination_account_id",
					"status",
					"raw_data",
				),
		).
		Model((*models.Payment)(nil)).
		TableExpr("_data").
		Set("amount = _data.amount").
		Set("type = _data.type").
		Set("scheme = _data.scheme").
		Set("asset = _data.asset").
		Set("source_account_id = _data.source_account_id").
		Set("destination_account_id = _data.destination_account_id").
		Set("status = _data.status").
		Set("raw_data = _data.raw_data").
		Where(`(payment.id = _data.id) AND 
			(payment.amount != _data.amount OR payment.type != _data.type OR payment.scheme != _data.scheme OR
				payment.asset != _data.asset OR payment.source_account_id != _data.source_account_id OR 
				payment.destination_account_id != _data.destination_account_id OR payment.status != _data.status)`).
		Returning("payment.id").
		Scan(ctx, &idsUpdated)
	if err != nil {
		return nil, e("failed to update payments", err)
	}

	idsUpdatedMap := make(map[string]struct{})
	for _, id := range idsUpdated {
		idsUpdatedMap[id] = struct{}{}
	}

	paymentsToInsert := make([]*models.Payment, 0, len(payments))
	for _, payment := range payments {
		if _, ok := idsUpdatedMap[payment.ID.String()]; !ok {
			paymentsToInsert = append(paymentsToInsert, payment)
		}
	}

	var idsInserted []string
	if len(paymentsToInsert) > 0 {
		err = s.db.NewInsert().
			Model(&paymentsToInsert).
			On("CONFLICT (id) DO NOTHING").
			Returning("payment.id").
			Scan(ctx, &idsInserted)
		if err != nil {
			return nil, e("failed to create payments", err)
		}
	}

	res := make([]*models.PaymentID, 0, len(idsUpdated)+len(idsInserted))
	for _, id := range idsUpdated {
		res = append(res, models.MustPaymentIDFromString(id))
	}
	for _, id := range idsInserted {
		res = append(res, models.MustPaymentIDFromString(id))
	}

	var adjustments []*models.Adjustment
	var metadata []*models.Metadata

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
			return nil, e("failed to create adjustments", err)
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
			return nil, e("failed to create metadata", err)
		}
	}

	return res, nil
}
