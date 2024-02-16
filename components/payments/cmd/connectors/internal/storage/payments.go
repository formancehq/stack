package storage

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/models"
)

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
					"created_at",
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
		Set("created_at = _data.created_at").
		Where(`(payment.id = _data.id) AND 
			(payment.created_at != _data.created_at OR payment.amount != _data.amount OR payment.type != _data.type OR
				payment.scheme != _data.scheme OR payment.asset != _data.asset OR payment.source_account_id != _data.source_account_id OR 
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

	return res, nil
}

func (s *Storage) UpsertPaymentsAdjustments(ctx context.Context, adjustments []*models.PaymentAdjustment) error {
	if len(adjustments) == 0 {
		return nil
	}

	_, err := s.db.NewInsert().
		Model(&adjustments).
		On("CONFLICT (reference) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("failed to create adjustments", err)
	}

	return nil
}

func (s *Storage) UpsertPaymentsMetadata(ctx context.Context, metadata []*models.PaymentMetadata) error {
	if len(metadata) == 0 {
		return nil
	}

	_, err := s.db.NewInsert().
		Model(&metadata).
		On("CONFLICT (payment_id, key) DO UPDATE").
		Set("value = EXCLUDED.value").
		Set("changelog = metadata.changelog || EXCLUDED.changelog").
		Exec(ctx)
	if err != nil {
		return e("failed to create metadata", err)
	}

	return nil
}
