package worker

import (
	"context"

	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error

	CreateTransactionBasedReconciliation(ctx context.Context, transactionBasedReconciliation *models.ReconciliationTransactionBased) (bool, error)
	GetRule(ctx context.Context, id uuid.UUID) (*models.Rule, error)
	ListPolicies(ctx context.Context, q storage.ListPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error)
	ListTransactionBasedReconciliations(ctx context.Context, q storage.ListTransactionBasedReconciliationsQuery) (*bunpaginate.Cursor[models.ReconciliationTransactionBased], error)
}
