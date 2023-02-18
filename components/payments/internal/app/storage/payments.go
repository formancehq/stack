package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/uptrace/bun"

	"github.com/formancehq/payments/internal/app/models"
)

func (s *Storage) ListPayments(ctx context.Context, pagination Paginator) ([]*models.Payment, PaginationDetails, error) {
	var payments []*models.Payment

	query := s.db.NewSelect().
		Model(&payments).
		Relation("Account").
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
		return nil, fmt.Errorf("failed to get payment %s: %w", id, err)
	}

	return &payment, nil
}

func (s *Storage) UpsertPayments(ctx context.Context, provider models.ConnectorProvider, payments []*models.Payment) error {
	if len(payments) == 0 {
		return nil
	}

	connector, err := s.GetConnector(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to get connector: %w", err)
	}

	var accountReferences []string

	for i := range payments {
		payments[i].ConnectorID = connector.ID

		if payments[i].Account != nil && payments[i].Account.Reference != "" {
			accountReferences = append(accountReferences, payments[i].Account.Reference)
		}
	}

	if len(accountReferences) > 0 {
		var accounts []models.Account

		err = s.db.NewSelect().Model(&accounts).
			Where("reference IN (?)", bun.In(accountReferences)).
			Scan(ctx)
		if err != nil {
			return e("failed to get accounts", err)
		}

		for i := range payments {
			if payments[i].Account != nil && payments[i].Account.Reference != "" {
				for j := range accounts {
					if accounts[j].Reference == payments[i].Account.Reference {
						payments[i].AccountID = accounts[j].ID
					}
				}
			}
		}
	}

	_, err = s.db.NewInsert().
		Model(&payments).
		On("CONFLICT (reference) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("failed to create payments", err)
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

	// TODO: Not working with dummy pay connector
	// Need to check with a real connector
	//err = s.UpdateTransfersFromPayments(ctx, payments)
	//if err != nil {
	//	return e("failed to update transfers", err)
	//}

	return nil
}
