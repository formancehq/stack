package service

import (
	"context"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) GetTransactionBasedReconciliationDetails(ctx context.Context, q storage.ListTransactionBasedReconciliationsQuery) (*bunpaginate.Cursor[models.ReconciliationTransactionBased], error) {
	reconciliations, err := s.store.ListTransactionBasedReconciliations(ctx, q)
	return reconciliations, newStorageError(err, "failed to list transaction based reconciliations")
}

func (s *Service) handleTransactionBasedReconciliation(ctx context.Context, reconciliation *models.Reconciliation) error {
	// TODO(polo): add history fetching
	return nil
}
