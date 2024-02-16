package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

func (s *Storage) CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := tx.NewInsert().
		Column("id", "created_at", "scheduled_at", "description", "type", "destination_account_id", "provider", "connector_id", "initial_amount", "amount", "asset", "metadata").
		Model(transferInitiation)

	if transferInitiation.SourceAccountID != nil {
		query = query.Column("source_account_id")
	}

	_, err = query.Exec(ctx)
	if err != nil {
		return e("failed to create transfer initiation", err)
	}

	for _, adjustment := range transferInitiation.RelatedAdjustments {
		adj := adjustment
		if _, err := tx.NewInsert().Model(adj).Exec(ctx); err != nil {
			return e("failed to add adjustment", err)
		}
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *Storage) AddTransferInitiationAdjustment(ctx context.Context, adjustment *models.TransferInitiationAdjustment) error {
	if _, err := s.db.NewInsert().Model(adjustment).Exec(ctx); err != nil {
		return e("failed to add adjustment", err)
	}

	return nil
}

func (s *Storage) ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error) {
	var transferInitiation models.TransferInitiation

	query := s.db.NewSelect().
		Column("id", "created_at", "scheduled_at", "description", "type", "source_account_id", "destination_account_id", "provider", "connector_id", "amount", "asset", "metadata").
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

func (s *Storage) AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, createdAt time.Time, metadata map[string]string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if paymentID == nil {
		return errors.New("payment id is nil")
	}

	_, err = tx.NewInsert().
		Column("transfer_initiation_id", "payment_id", "created_at", "status").
		Model(&models.TransferInitiationPayment{
			TransferInitiationID: id,
			PaymentID:            *paymentID,
			CreatedAt:            createdAt,
			Status:               models.TransferInitiationStatusProcessing,
		}).
		Exec(ctx)
	if err != nil {
		return e("failed to add transfer initiation payment id", err)
	}

	if metadata != nil {
		_, err := tx.NewUpdate().
			Model((*models.TransferInitiation)(nil)).
			Set("metadata = ?", metadata).
			Where("id = ?", id).
			Exec(ctx)
		if err != nil {
			return e("failed to add metadata", err)
		}
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *Storage) UpdateTransferInitiationPaymentsStatus(
	ctx context.Context,
	id models.TransferInitiationID,
	paymentID *models.PaymentID,
	adjustment *models.TransferInitiationAdjustment,
) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if paymentID != nil {
		query := tx.NewUpdate().
			Model((*models.TransferInitiationPayment)(nil)).
			Set("status = ?", adjustment.Status)

		if adjustment.Error != "" {
			query = query.Set("error = ?", adjustment.Error)
		}

		_, err := query.
			Where("transfer_initiation_id = ?", id).
			Where("payment_id = ?", paymentID).
			Exec(ctx)
		if err != nil {
			return e("failed to update transfer initiation status", err)
		}
	}

	if _, err = tx.NewInsert().Model(adjustment).Exec(ctx); err != nil {
		return e("failed to add adjustment", err)
	}

	return e("failed to commit transaction", tx.Commit())
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
