package service

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ListTransferInitiations(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.TransferInitiation, storage.PaginationDetails, error) {
	transferInitiations, paginationDetails, err := s.store.ListTransferInitiations(ctx, pagination)
	return transferInitiations, paginationDetails, newStorageError(err, "listing transfer initiations")
}

func (s *Service) ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error) {
	transferInitiation, err := s.store.GetTransferInitiation(ctx, id)
	return transferInitiation, newStorageError(err, "reading transfer initiation")
}
