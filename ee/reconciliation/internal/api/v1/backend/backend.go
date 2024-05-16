package backend

import (
	"context"

	"github.com/formancehq/reconciliation/internal/api/v1/service"
	"github.com/formancehq/reconciliation/internal/models"
	storage "github.com/formancehq/reconciliation/internal/storage/v1"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Service
type Service interface {
	Reconciliation(ctx context.Context, policyID string, req *service.ReconciliationRequest) (*models.Reconciliation, error)
	GetReconciliation(ctx context.Context, id string) (*models.Reconciliation, error)
	ListReconciliations(ctx context.Context, q storage.GetReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error)

	CreatePolicy(ctx context.Context, req *service.CreatePolicyRequest) (*models.Policy, error)
	DeletePolicy(ctx context.Context, id string) error
	GetPolicy(ctx context.Context, id string) (*models.Policy, error)
	ListPolicies(ctx context.Context, q storage.GetPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error)
}

type Backend interface {
	GetService() Service
}

type DefaultBackend struct {
	service Service
}

func (d DefaultBackend) GetService() Service {
	return d.service
}

func NewDefaultBackend(service Service) Backend {
	return &DefaultBackend{
		service: service,
	}
}
