package service

import (
	"context"

	models "github.com/formancehq/reconciliation/internal/models/v2"
	storage "github.com/formancehq/reconciliation/internal/storage/v2"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error

	CreateRule(ctx context.Context, rule *models.Rule) error
	DeleteRule(ctx context.Context, id string) error
	GetRule(ctx context.Context, id string) (*models.Rule, error)
	ListRules(ctx context.Context, q storage.ListRulesQuery) (*bunpaginate.Cursor[models.Rule], error)

	CreatePolicy(ctx context.Context, policy *models.Policy) error
	UpdatePolicyRules(ctx context.Context, id uuid.UUID, rules []string) error
	UpdatePolicyStatus(ctx context.Context, id uuid.UUID, enabled bool) error
	DeletePolicy(ctx context.Context, id uuid.UUID) error
	GetPolicy(ctx context.Context, id uuid.UUID) (*models.Policy, error)
	ListPolicies(ctx context.Context, q storage.ListPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}
