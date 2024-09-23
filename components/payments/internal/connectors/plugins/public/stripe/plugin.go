package stripe

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/connectors/plugins/public/stripe/client"
	"github.com/formancehq/payments/internal/models"
	stripesdk "github.com/stripe/stripe-go/v79"
)

const PageLimit = int64(20)

type Plugin struct {
	StripeAPIBackend stripesdk.Backend // override in tests to mock
	client           client.Client
}

func (p *Plugin) SetClient(client client.Client) error {
	if p.client != nil {
		return fmt.Errorf("client is not intended to be overwritten after install")
	}
	p.client = client
	return nil
}

func (p *Plugin) Install(_ context.Context, req models.InstallRequest) (models.InstallResponse, error) {
	config, err := unmarshalAndValidateConfig(req.Config)
	if err != nil {
		return models.InstallResponse{}, err
	}

	p.client = client.New(p.StripeAPIBackend, config.APIKey)
	return models.InstallResponse{
		Capabilities: capabilities,
		Workflow:     workflow(),
	}, nil
}

func (p *Plugin) Uninstall(ctx context.Context, req models.UninstallRequest) (models.UninstallResponse, error) {
	return models.UninstallResponse{}, nil
}

func (p *Plugin) FetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	if p.client == nil {
		return models.FetchNextExternalAccountsResponse{}, plugins.ErrNotYetInstalled
	}
	return p.fetchNextExternalAccounts(ctx, req)
}
