package client

import (
	"context"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
)

type sdkFormanceClient struct {
	client *sdk.Formance
}

func NewSDKFormance(client *sdk.Formance) *sdkFormanceClient {
	return &sdkFormanceClient{
		client: client,
	}
}

type SDKSearchFormance interface {
	RawSearch(context.Context, map[string]interface{}) (*operations.SearchResponse, error)
}

func (s *sdkFormanceClient) RawSearch(ctx context.Context, rawQuery map[string]interface{}) (*operations.SearchResponse, error) {
	return s.client.Search.Search(ctx, shared.Query{Raw: rawQuery})
}

type SDKFormance interface {
	PaymentsgetServerInfo(ctx context.Context) (*operations.PaymentsgetServerInfoResponse, error)
	GetPoolBalances(ctx context.Context, req operations.GetPoolBalancesRequest) (*operations.GetPoolBalancesResponse, error)
	V2GetInfo(ctx context.Context) (*operations.V2GetInfoResponse, error)
	V2GetBalancesAggregated(ctx context.Context, req operations.V2GetBalancesAggregatedRequest) (*operations.V2GetBalancesAggregatedResponse, error)
}

func (s *sdkFormanceClient) PaymentsgetServerInfo(ctx context.Context) (*operations.PaymentsgetServerInfoResponse, error) {
	return s.client.Payments.PaymentsgetServerInfo(ctx)
}

func (s *sdkFormanceClient) GetPoolBalances(ctx context.Context, req operations.GetPoolBalancesRequest) (*operations.GetPoolBalancesResponse, error) {
	return s.client.Payments.GetPoolBalances(ctx, req)
}

func (s *sdkFormanceClient) V2GetInfo(ctx context.Context) (*operations.V2GetInfoResponse, error) {
	return s.client.Ledger.V2GetInfo(ctx)
}

func (s *sdkFormanceClient) V2GetBalancesAggregated(ctx context.Context, req operations.V2GetBalancesAggregatedRequest) (*operations.V2GetBalancesAggregatedResponse, error) {
	return s.client.Ledger.V2GetBalancesAggregated(ctx, req)
}

var _ SDKFormance = (*sdkFormanceClient)(nil)
