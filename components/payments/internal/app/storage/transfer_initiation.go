package storage

import (
	"context"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/app/models"
)

func (s *Storage) CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error {
	_, err := s.db.NewInsert().
		Column("id", "payment_id", "created_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "amount", "asset", "status", "error").
		Model(transferInitiation).
		Exec(ctx)
	if err != nil {
		return e("failed to create transfer initiation", err)
	}

	return nil
}

func (s *Storage) ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error) {
	var transferInitiation models.TransferInitiation

	query := s.db.NewSelect().
		Column("id", "payment_id", "created_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "amount", "asset", "status", "error").
		Model(&transferInitiation).
		Where("id = ?", id)

	err := query.Scan(ctx)
	if err != nil {
		return nil, e("failed to get transfer initiation", err)
	}

	return &transferInitiation, nil
}

func (s *Storage) ListTransferInitiations(ctx context.Context, pagination Paginator) ([]*models.TransferInitiation, PaginationDetails, error) {
	var tfs []*models.TransferInitiation

	query := s.db.NewSelect().
		Column("id", "payment_id", "created_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "amount", "asset", "status", "error").
		Model(&tfs)

	query = pagination.apply(query, "transfer_initiation.created_at")

	err := query.Scan(ctx)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to list payments", err)
	}

	var (
		hasMore                       = len(tfs) > pagination.pageSize
		hasPrevious                   bool
		firstReference, lastReference string
	)

	if hasMore {
		if pagination.cursor.Next || pagination.cursor.Reference == "" {
			tfs = tfs[:pagination.pageSize]
		} else {
			tfs = tfs[1:]
		}
	}

	sort.Slice(tfs, func(i, j int) bool {
		return tfs[i].CreatedAt.After(tfs[j].CreatedAt)
	})

	if len(tfs) > 0 {
		firstReference = tfs[0].CreatedAt.Format(time.RFC3339Nano)
		lastReference = tfs[len(tfs)-1].CreatedAt.Format(time.RFC3339Nano)

		query = s.db.NewSelect().
			Column("id", "payment_id", "created_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "amount", "asset", "status", "error").
			Model(&tfs)

		hasPrevious, err = pagination.hasPrevious(ctx, query, "transfer_initiation.created_at", firstReference)
		if err != nil {
			return nil, PaginationDetails{}, e("failed to check if there is a previous page", err)
		}
	}

	paginationDetails, err := pagination.paginationDetails(hasMore, hasPrevious, firstReference, lastReference)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to get pagination details", err)
	}

	return tfs, paginationDetails, nil
}

func (s *Storage) UpdateTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID models.PaymentID, updatedAt time.Time) error {
	_, err := s.db.NewUpdate().
		Model((*models.TransferInitiation)(nil)).
		Set("payment_id = ?", paymentID).
		Set("updated_at = ?", updatedAt).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("failed to update transfer initiation payment id", err)
	}

	return nil
}

func (s *Storage) UpdateTransferInitiationError(ctx context.Context, id models.TransferInitiationID, error string, updatedAt time.Time) error {
	_, err := s.db.NewUpdate().
		Model((*models.TransferInitiation)(nil)).
		Set("status = ?", models.TransferInitiationStatusFailed).
		Set("error = ?", error).
		Set("updated_at = ?", updatedAt).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("failed to update transfer initiation error", err)
	}

	return nil
}

func (s *Storage) UpdateTransferInitiationStatus(ctx context.Context, id models.TransferInitiationID, status models.TransferInitiationStatus, errorMessage string, updatedAt time.Time) error {
	query := s.db.NewUpdate().
		Model((*models.TransferInitiation)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", updatedAt)

	if errorMessage != "" {
		query = query.Set("error = ?", errorMessage)
	}

	_, err := query.Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("failed to update transfer initiation status", err)
	}

	return nil
}

func (s *Storage) DeleteTransferInitiation(ctx context.Context, id models.TransferInitiationID) error {
	_, err := s.db.NewDelete().
		Model((*models.TransferInitiation)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("failed to delete transfer initiation", err)
	}

	return nil
}
