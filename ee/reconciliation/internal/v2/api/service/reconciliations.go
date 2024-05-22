package service

import (
	"context"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
)

func (s *Service) GetReconciliation(ctx context.Context, id string) (*models.Reconciliation, error) {
	reconciliationID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	reconciliation, err := s.store.GetReconciliation(ctx, reconciliationID)
	if err != nil {
		return nil, newStorageError(err, "getting reconciliation")
	}

	return reconciliation, nil
}

func (s *Service) ListReconciliations(ctx context.Context, q storage.ListReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error) {
	reconciliations, err := s.store.ListReconciliations(ctx, q)
	return reconciliations, newStorageError(err, "listing reconciliations")
}
