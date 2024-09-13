package grpc

import (
	"context"
	"os"

	"github.com/formancehq/payments/internal/connectors/grpc/proto/services"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type PSP interface {
	Install(ctx context.Context, in *services.InstallRequest) (*services.InstallResponse, error)
	Uninstall(ctx context.Context, in *services.UninstallRequest) (*services.UninstallResponse, error)

	FetchNextOthers(ctx context.Context, in *services.FetchNextOthersRequest) (*services.FetchNextOthersResponse, error)
	FetchNextPayments(ctx context.Context, in *services.FetchNextPaymentsRequest) (*services.FetchNextPaymentsResponse, error)
	FetchNextAccounts(ctx context.Context, in *services.FetchNextAccountsRequest) (*services.FetchNextAccountsResponse, error)
	FetchNextBalances(ctx context.Context, in *services.FetchNextBalancesRequest) (*services.FetchNextBalancesResponse, error)
	FetchNextExternalAccounts(ctx context.Context, in *services.FetchNextExternalAccountsRequest) (*services.FetchNextExternalAccountsResponse, error)

	CreateBankAccount(ctx context.Context, in *services.CreateBankAccountRequest) (*services.CreateBankAccountResponse, error)

	CreateWebhooks(ctx context.Context, in *services.CreateWebhooksRequest) (*services.CreateWebhooksResponse, error)
	TranslateWebhook(ctx context.Context, in *services.TranslateWebhookRequest) (*services.TranslateWebhookResponse, error)
}

type PSPGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl PSP
}

func (p *PSPGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	services.RegisterPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *PSPGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: services.NewPluginClient(c)}, nil
}

var PluginMap = map[string]plugin.Plugin{
	"psp": &PSPGRPCPlugin{},
}

var _ plugin.GRPCPlugin = &PSPGRPCPlugin{}

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "PLUGIN_KEY",
	MagicCookieValue: os.Getenv("PLUGIN_MAGIC_COOKIE"),
}
