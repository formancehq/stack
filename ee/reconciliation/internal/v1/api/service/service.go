package service

import (
	"context"

	"github.com/formancehq/reconciliation/internal/client"
	"github.com/formancehq/reconciliation/internal/v1/models"
	"github.com/formancehq/reconciliation/internal/v1/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error
	CreatePolicy(ctx context.Context, policy *models.Policy) error
	DeletePolicy(ctx context.Context, id uuid.UUID) error
	GetPolicy(ctx context.Context, id uuid.UUID) (*models.Policy, error)
	ListPolicies(ctx context.Context, q storage.GetPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error)

	CreateReconciation(ctx context.Context, reco *models.Reconciliation) error
	GetReconciliation(ctx context.Context, id uuid.UUID) (*models.Reconciliation, error)
	ListReconciliations(ctx context.Context, q storage.GetReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error)
}

type Service struct {
	store  Store
	client client.SDKFormance
}

func NewService(store Store, client client.SDKFormance) *Service {
	return &Service{
		store:  store,
		client: client,
	}
}
