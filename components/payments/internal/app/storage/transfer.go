package storage

import (
	"context"

	"github.com/pkg/errors"

	"github.com/uptrace/bun"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"

	"github.com/formancehq/payments/internal/app/models"
)

func (s *Storage) CreateNewTransfer(ctx context.Context, name models.ConnectorProvider,
	source, destination, currency string, amount int64,
) (models.Transfer, error) {
	connector, err := s.GetConnector(ctx, name)
	if err != nil {
		return models.Transfer{}, err
	}

	newTransfer := models.Transfer{
		ConnectorID: connector.ID,
		Amount:      amount,
		Status:      models.TransferStatusPending,
		Currency:    currency,
		Source:      source,
		Destination: destination,
	}

	_, err = s.db.NewInsert().Model(&newTransfer).Exec(ctx)
	if err != nil {
		return models.Transfer{}, e("failed to create new transfer", err)
	}

	return newTransfer, nil
}

func (s *Storage) ListTransfers(ctx context.Context, name models.ConnectorProvider) ([]models.Transfer, error) {
	connector, err := s.GetConnector(ctx, name)
	if err != nil {
		return nil, err
	}

	var transfers []models.Transfer
	err = s.db.NewSelect().Model(&transfers).Where("connector_id = ?", connector.ID).Scan(ctx)
	if err != nil {
		return nil, e("failed to list transfers", err)
	}

	return transfers, nil
}

func (s *Storage) GetTransfer(ctx context.Context, transferID uuid.UUID) (models.Transfer, error) {
	var transfer models.Transfer

	err := s.db.NewSelect().Model(&transfer).Where("id = ?", transferID).Scan(ctx)
	if err != nil {
		return models.Transfer{}, e("failed to get transfer", err)
	}

	return transfer, nil
}

func (s *Storage) UpdateTransferStatus(ctx context.Context, transferID uuid.UUID,
	status models.TransferStatus, reference, transferErr string,
) error {
	_, err := s.db.NewUpdate().Model(&models.Transfer{}).
		Set("status = ?", status).
		Set("reference = ?", reference).
		Set("error = ?", transferErr).
		Where("id = ?", transferID).
		Exec(ctx)
	if err != nil {
		return e("failed to update transfer status", err)
	}

	return nil
}

func (s *Storage) UpdateTransfersFromPayments(ctx context.Context, payments []*models.Payment) error {
	var transfers []models.Transfer

	if len(payments) == 0 {
		return nil
	}

	paymentReferences := make([]string, len(payments))
	for paymentIdx := range payments {
		paymentReferences[paymentIdx] = payments[paymentIdx].Reference
	}

	err := s.db.NewSelect().Model(&transfers).
		Where("reference IN (?)", bun.In(paymentReferences)).
		Where("payment_id IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return e("failed to get transfer", err)
	}

	if len(transfers) == 0 {
		return nil
	}

	for transferIdx := range transfers {
		if transfers[transferIdx].Reference == nil {
			continue
		}

		for paymentIdx := range payments {
			if payments[paymentIdx].Reference == *transfers[transferIdx].Reference {
				transfers[transferIdx].PaymentID = &payments[paymentIdx].ID
			}
		}
	}

	_, err = s.db.NewUpdate().Model(&transfers).Exec(ctx)
	if err != nil {
		return e("failed to update transfers", err)
	}

	return nil
}
