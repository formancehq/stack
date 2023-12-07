package service

import (
	"context"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error
	CreatePolicy(ctx context.Context, policy *models.Policy) error
	DeletePolicy(ctx context.Context, id uuid.UUID) error
	GetPolicy(ctx context.Context, id uuid.UUID) (*models.Policy, error)
	ListPolicies(ctx context.Context, q storage.GetPoliciesQuery) (*api.Cursor[models.Policy], error)

	CreateReconciation(ctx context.Context, reco *models.Reconciliation) error
	GetReconciliation(ctx context.Context, id uuid.UUID) (*models.Reconciliation, error)
	ListReconciliations(ctx context.Context, q storage.GetReconciliationsQuery) (*api.Cursor[models.Reconciliation], error)
}

type Service struct {
	store  Store
	client *sdk.Formance
}

func NewService(store Store, client *sdk.Formance) *Service {
	return &Service{
		store:  store,
		client: client,
	}
}
