package storage

import (
	"context"
	"database/sql"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Storage) CreateTransferReversal(ctx context.Context, transferReversal *models.TransferReversal) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.NewInsert().Model(transferReversal).Exec(ctx)
	if err != nil {
		return e("failed to create transfer reversal", err)
	}

	adjustment := &models.TransferInitiationAdjustment{
		ID:                   uuid.New(),
		TransferInitiationID: transferReversal.TransferInitiationID,
		CreatedAt:            time.Now().UTC(),
		Status:               models.TransferInitiationStatusAskReversed,
		Error:                "",
		Metadata:             transferReversal.Metadata,
	}

	if _, err = tx.NewInsert().Model(adjustment).Exec(ctx); err != nil {
		return e("failed to create adjustment", err)
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *Storage) UpdateTransferReversalStatus(
	ctx context.Context,
	transfer *models.TransferInitiation,
	transferReversal *models.TransferReversal,
	adjustment *models.TransferInitiationAdjustment,
) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC()

	_, err = tx.NewUpdate().
		Model(transferReversal).
		Set("status = ?", transferReversal.Status).
		Set("error = ?", transferReversal.Error).
		Set("updated_at = ?", now).
		Where("id = ?", transferReversal.ID).
		Exec(ctx)
	if err != nil {
		return e("failed to update transfer reversal status", err)
	}

	if transferReversal.Status == models.TransferReversalStatusProcessed {
		var amount *big.Int
		err = tx.NewUpdate().
			Model((*models.TransferInitiation)(nil)).
			Set("amount = amount - ?", transferReversal.Amount).
			Where("id = ?", transferReversal.TransferInitiationID).
			Returning("amount").
			Scan(ctx, &amount)
		if err != nil {
			return e("failed to update transfer initiation amount", err)
		}

		switch amount.Cmp(big.NewInt(0)) {
		case 0:
			// amount == 0, so we can mark the transfer as reversed
			adjustment.Status = models.TransferInitiationStatusReversed
		case 1:
			// amount > 0, so we can mark the transfer as partially reversed
			adjustment.Status = models.TransferInitiationStatusPartiallyReversed
		case -1:
			// Should not happened since we have checks in postgres
			return errors.New("transfer reversal amount is greater than transfer initiation amount")
		}

		transfer.Amount = amount
	}

	if _, err := tx.NewInsert().Model(adjustment).Exec(ctx); err != nil {
		return e("failed to add adjustment", err)
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *Storage) GetTransferReversal(ctx context.Context, id models.TransferReversalID) (*models.TransferReversal, error) {
	var ret models.TransferReversal
	err := s.db.NewSelect().
		Model(&ret).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get transfer reversal", err)
	}

	return &ret, nil
}
