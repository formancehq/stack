package backend

import (
	"context"

	"github.com/formancehq/reconciliation/internal/v2/api/service"
	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Service
type Service interface {
	CreateRule(ctx context.Context, req *service.CreateRuleRequest) (*models.Rule, error)
	DeleteRule(ctx context.Context, id string) error
	GetRule(ctx context.Context, id string) (*models.Rule, error)
	ListRules(ctx context.Context, q storage.ListRulesQuery) (*bunpaginate.Cursor[models.Rule], error)

	CreatePolicy(ctx context.Context, req *service.CreatePolicyRequest) (*models.Policy, error)
	UpdatePolicyRules(ctx context.Context, id string, req *service.UpdatePolicyRulesRequest) error
	EnablePolicy(ctx context.Context, id string) error
	DisablePolicy(ctx context.Context, id string) error
	DeletePolicy(ctx context.Context, id string) error
	GetPolicy(ctx context.Context, id string) (*models.Policy, error)
	ListPolicies(ctx context.Context, q storage.ListPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error)
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
