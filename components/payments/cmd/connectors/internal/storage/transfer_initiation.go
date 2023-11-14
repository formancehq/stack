package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
)

func (s *Storage) CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error {
	_, err := s.db.NewInsert().
		Column("id", "created_at", "scheduled_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "connector_id", "amount", "asset", "status", "error").
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
		Column("id", "created_at", "scheduled_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "connector_id", "amount", "asset", "status", "error").
		Model(&transferInitiation).
		Where("id = ?", id)

	err := query.Scan(ctx)
	if err != nil {
		return nil, e("failed to get transfer initiation", err)
	}

	transferInitiation.RelatedPayments, err = s.ReadTransferInitiationPayments(ctx, id)
	if err != nil {
		return nil, e("failed to get transfer initiation payments", err)
	}

	return &transferInitiation, nil
}

func (s *Storage) ReadTransferInitiationPayments(ctx context.Context, id models.TransferInitiationID) ([]*models.TransferInitiationPayments, error) {
	var payments []*models.TransferInitiationPayments

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

func (s *Storage) ListTransferInitiations(ctx context.Context, pagination PaginatorQuery) ([]*models.TransferInitiation, PaginationDetails, error) {
	var tfs []*models.TransferInitiation

	query := s.db.NewSelect().
		Column("id", "created_at", "scheduled_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "connector_id", "amount", "asset", "status", "error").
		Model(&tfs)

	if pagination.queryBuilder != nil {
		where, args, err := s.transferInitiationQueryContext(pagination.queryBuilder)
		if err != nil {
			// TODO: handle error
			panic(err)
		}
		query = query.Where(where, args...)
	}

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
			Column("id", "created_at", "updated_at", "description", "type", "source_account_id", "destination_account_id", "provider", "connector_id", "amount", "asset", "status", "error").
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

func (s *Storage) AddTransferInitiationPaymentID(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, createdAt time.Time) error {
	if paymentID == nil {
		return errors.New("payment id is nil")
	}

	_, err := s.db.NewInsert().
		Column("transfer_initiation_id", "payment_id", "created_at", "status").
		Model(&models.TransferInitiationPayments{
			TransferInitiationID: id,
			PaymentID:            *paymentID,
			CreatedAt:            createdAt,
			Status:               models.TransferInitiationStatusProcessing,
		}).
		Exec(ctx)
	if err != nil {
		return e("failed to add transfer initiation payment id", err)
	}

	return nil
}

func (s *Storage) updateTransferInitiationStatus(
	ctx context.Context,
	id models.TransferInitiationID,
	status models.TransferInitiationStatus,
	errorMessage string,
	attempts int,
	updatedAt time.Time,
) error {
	query := s.db.NewUpdate().
		Model((*models.TransferInitiation)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", updatedAt).
		Set("error = ?", errorMessage)

	if attempts > 0 {
		query = query.Set("attempts = ?", attempts)
	}

	_, err := query.Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("failed to update transfer initiation status", err)
	}

	return nil
}

func (s *Storage) UpdateTransferInitiationPaymentsStatus(
	ctx context.Context,
	id models.TransferInitiationID,
	paymentID *models.PaymentID,
	status models.TransferInitiationStatus,
	errorMessage string,
	attempts int,
	updatedAt time.Time,
) error {
	if paymentID != nil {
		query := s.db.NewUpdate().
			Model((*models.TransferInitiationPayments)(nil)).
			Set("status = ?", status)

		if errorMessage != "" {
			query = query.Set("error = ?", errorMessage)
		}

		_, err := query.
			Where("transfer_initiation_id = ?", id).
			Where("payment_id = ?", paymentID).
			Exec(ctx)
		if err != nil {
			return e("failed to update transfer initiation status", err)
		}
	}

	return s.updateTransferInitiationStatus(ctx, id, status, errorMessage, attempts, updatedAt)
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

func (s *Storage) transferInitiationQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "source_account_id", key == "destination_account_id":
			if operator != "$match" {
				return "", nil, fmt.Errorf("'%s' columns can only be used with $match", key)
			}

			switch accountID := value.(type) {
			case string:
				return fmt.Sprintf("%s = ?", key), []any{accountID}, nil
			default:
				return "", nil, fmt.Errorf("unexpected type %T for column '%s'", accountID, key)
			}
		default:
			return "", nil, fmt.Errorf("unknown key '%s' when building query", key)
		}
	}))
}
