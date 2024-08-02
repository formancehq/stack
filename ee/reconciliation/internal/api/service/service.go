package service

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/reconciliation/internal/storage"
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
	client SDKFormance
}

func NewService(store Store, client SDKFormance) *Service {
	return &Service{
		store:  store,
		client: client,
	}
}

type SDKFormance interface {
	PaymentsgetServerInfo(ctx context.Context) (*operations.PaymentsgetServerInfoResponse, error)
	GetPoolBalances(ctx context.Context, req operations.GetPoolBalancesRequest) (*operations.GetPoolBalancesResponse, error)
	V2GetInfo(ctx context.Context) (*operations.V2GetInfoResponse, error)
	V2GetBalancesAggregated(ctx context.Context, req operations.V2GetBalancesAggregatedRequest) (*operations.V2GetBalancesAggregatedResponse, error)
}

type sdkFormanceClient struct {
	client *sdk.Formance
}

func NewSDKFormance(client *sdk.Formance) *sdkFormanceClient {
	return &sdkFormanceClient{
		client: client,
	}
}

func (s *sdkFormanceClient) PaymentsgetServerInfo(ctx context.Context) (*operations.PaymentsgetServerInfoResponse, error) {
	return s.client.Payments.V1.PaymentsgetServerInfo(ctx)
}

func (s *sdkFormanceClient) GetPoolBalances(ctx context.Context, req operations.GetPoolBalancesRequest) (*operations.GetPoolBalancesResponse, error) {
	return s.client.Payments.V1.GetPoolBalances(ctx, req)
}

func (s *sdkFormanceClient) V2GetInfo(ctx context.Context) (*operations.V2GetInfoResponse, error) {
	return s.client.Ledger.V2.GetInfo(ctx)
}

func (s *sdkFormanceClient) V2GetBalancesAggregated(ctx context.Context, req operations.V2GetBalancesAggregatedRequest) (*operations.V2GetBalancesAggregatedResponse, error) {
	return s.client.Ledger.V2.GetBalancesAggregated(ctx, req)
}

var _ SDKFormance = (*sdkFormanceClient)(nil)
