package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) ListPayments(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Payment, storage.PaginationDetails, error) {
	payments, paginationDetails, err := s.store.ListPayments(ctx, pagination)
	return payments, paginationDetails, newStorageError(err, "listing payments")
}

func (s *Service) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	_, err := models.PaymentIDFromString(id)
	if err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	payment, err := s.store.GetPayment(ctx, id)
	return payment, newStorageError(err, "getting payment")
}

type UpdateMetadataRequest map[string]string

func (s *Service) UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error {
	err := s.store.UpdatePaymentMetadata(ctx, paymentID, metadata)
	return newStorageError(err, "updating payment metadata")
}
