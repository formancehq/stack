package service

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ListTransferInitiations(ctx context.Context, q storage.ListTransferInitiationsQuery) (*bunpaginate.Cursor[models.TransferInitiation], error) {
	cursor, err := s.store.ListTransferInitiations(ctx, q)
	return cursor, newStorageError(err, "listing transfer initiations")
}

func (s *Service) ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error) {
	transferInitiation, err := s.store.GetTransferInitiation(ctx, id)
	return transferInitiation, newStorageError(err, "reading transfer initiation")
}
