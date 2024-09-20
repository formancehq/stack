package service

import (
	"context"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/google/uuid"
)

type mockSDKFormanceClient struct {
	ledgerVersion    string
	ledgerBalances   map[string]*big.Int
	paymentsVersion  string
	paymentsBalances map[string]*big.Int
}

func newMockSDKFormanceClient(
	ledgerVersion string,
	ledgerBalances map[string]*big.Int,
	paymentsVersion string,
	paymentsBalances map[string]*big.Int,
) *mockSDKFormanceClient {
	return &mockSDKFormanceClient{
		ledgerVersion:    ledgerVersion,
		ledgerBalances:   ledgerBalances,
		paymentsVersion:  paymentsVersion,
		paymentsBalances: paymentsBalances,
	}
}

func (s *mockSDKFormanceClient) PaymentsgetServerInfo(ctx context.Context) (*operations.PaymentsgetServerInfoResponse, error) {
	return &operations.PaymentsgetServerInfoResponse{
		ServerInfo: &shared.ServerInfo{
			Version: s.paymentsVersion,
		},
		StatusCode: http.StatusOK,
	}, nil
}

func (s *mockSDKFormanceClient) GetPoolBalances(ctx context.Context, req operations.GetPoolBalancesRequest) (*operations.GetPoolBalancesResponse, error) {
	poolBalances := make([]shared.PoolBalance, 0, len(s.paymentsBalances))
	for assetCode, balance := range s.paymentsBalances {
		poolBalances = append(poolBalances, shared.PoolBalance{
			Amount: balance,
			Asset:  assetCode,
		})
	}

	return &operations.GetPoolBalancesResponse{
		PoolBalancesResponse: &shared.PoolBalancesResponse{
			Data: shared.PoolBalances{
				Balances: poolBalances,
			},
		},
		StatusCode: http.StatusOK,
	}, nil
}

func (s *mockSDKFormanceClient) V2GetInfo(ctx context.Context) (*operations.V2GetInfoResponse, error) {
	return &operations.V2GetInfoResponse{
		StatusCode: http.StatusOK,
		V2ConfigInfoResponse: &shared.V2ConfigInfoResponse{
			Version: s.ledgerVersion,
		},
	}, nil
}

func (s *mockSDKFormanceClient) V2GetBalancesAggregated(ctx context.Context, req operations.V2GetBalancesAggregatedRequest) (*operations.V2GetBalancesAggregatedResponse, error) {
	balances := make(map[string]*big.Int)
	for assetCode, balance := range s.ledgerBalances {
		balances[assetCode] = balance
	}

	return &operations.V2GetBalancesAggregatedResponse{
		StatusCode: http.StatusOK,
		V2AggregateBalancesResponse: &shared.V2AggregateBalancesResponse{
			Data: balances,
		},
	}, nil
}

type mockStore struct {
}

func newMockStore() *mockStore {
	return &mockStore{}
}

func (s *mockStore) Ping() error {
	return nil
}

func (s *mockStore) CreatePolicy(ctx context.Context, policy *models.Policy) error {
	return nil
}

func (s *mockStore) DeletePolicy(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *mockStore) GetPolicy(ctx context.Context, id uuid.UUID) (*models.Policy, error) {
	return &models.Policy{
		ID:             id,
		CreatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Name:           "test",
		LedgerName:     "default",
		LedgerQuery:    map[string]interface{}{},
		PaymentsPoolID: uuid.New(),
	}, nil
}

func (s *mockStore) ListPolicies(ctx context.Context, q storage.GetPoliciesQuery) (*bunpaginate.Cursor[models.Policy], error) {
	return nil, nil
}

func (s *mockStore) CreateReconciation(ctx context.Context, reco *models.Reconciliation) error {
	return nil
}

func (s *mockStore) GetReconciliation(ctx context.Context, id uuid.UUID) (*models.Reconciliation, error) {
	return nil, nil
}

func (s *mockStore) ListReconciliations(ctx context.Context, q storage.GetReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error) {
	return nil, nil
}

var _ Store = (*mockStore)(nil)
