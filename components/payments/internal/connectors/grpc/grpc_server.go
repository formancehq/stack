package grpc

import (
	"context"

	"github.com/formancehq/payments/internal/connectors/grpc/proto/services"
)

var _ services.PluginServer = &GRPCServer{}

type GRPCServer struct {
	services.UnimplementedPluginServer
	// This is the real implementation
	Impl PSP
}

func (s *GRPCServer) Install(ctx context.Context, req *services.InstallRequest) (*services.InstallResponse, error) {
	return s.Impl.Install(ctx, req)
}

func (s *GRPCServer) Uninstall(ctx context.Context, req *services.UninstallRequest) (*services.UninstallResponse, error) {
	return s.Impl.Uninstall(ctx, req)
}

func (s *GRPCServer) FetchNextAccounts(ctx context.Context, req *services.FetchNextAccountsRequest) (*services.FetchNextAccountsResponse, error) {
	return s.Impl.FetchNextAccounts(ctx, req)
}

func (s *GRPCServer) FetchNextExternalAccounts(ctx context.Context, req *services.FetchNextExternalAccountsRequest) (*services.FetchNextExternalAccountsResponse, error) {
	return s.Impl.FetchNextExternalAccounts(ctx, req)
}

func (s *GRPCServer) FetchNextPayments(ctx context.Context, req *services.FetchNextPaymentsRequest) (*services.FetchNextPaymentsResponse, error) {
	return s.Impl.FetchNextPayments(ctx, req)
}

func (s *GRPCServer) FetchNextBalances(ctx context.Context, req *services.FetchNextBalancesRequest) (*services.FetchNextBalancesResponse, error) {
	return s.Impl.FetchNextBalances(ctx, req)
}

func (s *GRPCServer) FetchNextOthers(ctx context.Context, req *services.FetchNextOthersRequest) (*services.FetchNextOthersResponse, error) {
	return s.Impl.FetchNextOthers(ctx, req)
}

func (s *GRPCServer) CreateBankAccount(ctx context.Context, req *services.CreateBankAccountRequest) (*services.CreateBankAccountResponse, error) {
	return s.Impl.CreateBankAccount(ctx, req)
}

func (s *GRPCServer) CreateWebhooks(ctx context.Context, req *services.CreateWebhooksRequest) (*services.CreateWebhooksResponse, error) {
	return s.Impl.CreateWebhooks(ctx, req)
}

func (s *GRPCServer) TranslateWebhook(ctx context.Context, req *services.TranslateWebhookRequest) (*services.TranslateWebhookResponse, error) {
	return s.Impl.TranslateWebhook(ctx, req)
}
