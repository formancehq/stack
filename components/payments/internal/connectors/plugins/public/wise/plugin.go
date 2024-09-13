package wise

import (
	"context"
	"errors"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/connectors/plugins/public/wise/client"
	"github.com/formancehq/payments/internal/models"
)

type Plugin struct {
	config Config
	client *client.Client
}

func (p *Plugin) Install(ctx context.Context, req models.InstallRequest) (models.InstallResponse, error) {
	config, err := unmarshalAndValidateConfig(req.Config)
	if err != nil {
		return models.InstallResponse{}, err
	}

	client, err := client.New(config.APIKey)
	if err != nil {
		return models.InstallResponse{}, fmt.Errorf("failed to install wise plugin %w", err)
	}
	p.client = client
	p.config = config

	webhookConfigs = map[string]webhookConfig{
		"transfer_state_changed": {
			triggerOn: "transfers#state-change",
			urlPath:   "/transferstatechanged",
			fn:        p.translateTransferStateChangedWebhook,
			version:   "2.0.0",
		},
		"balance_update": {
			triggerOn: "balances#update",
			urlPath:   "/balanceupdate",
			fn:        p.translateBalanceUpdateWebhook,
			version:   "2.2.0",
		},
	}

	configs := make([]models.PSPWebhookConfig, 0, len(webhookConfigs))
	for name, config := range webhookConfigs {
		configs = append(configs, models.PSPWebhookConfig{
			Name:    name,
			URLPath: config.urlPath,
		})
	}

	return models.InstallResponse{
		Capabilities:    capabilities,
		Workflow:        workflow(),
		WebhooksConfigs: configs,
	}, nil
}

func (p Plugin) Uninstall(ctx context.Context, req models.UninstallRequest) (models.UninstallResponse, error) {
	return p.uninstall(ctx, req)
}

func (p Plugin) FetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	if p.client == nil {
		return models.FetchNextAccountsResponse{}, plugins.ErrNotYetInstalled
	}
	return p.fetchNextAccounts(ctx, req)
}

func (p Plugin) FetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	if p.client == nil {
		return models.FetchNextBalancesResponse{}, plugins.ErrNotYetInstalled
	}
	return p.fetchNextBalances(ctx, req)
}

func (p Plugin) FetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	if p.client == nil {
		return models.FetchNextExternalAccountsResponse{}, plugins.ErrNotYetInstalled
	}
	return p.fetchExternalAccounts(ctx, req)
}

func (p Plugin) FetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	if p.client == nil {
		return models.FetchNextPaymentsResponse{}, plugins.ErrNotYetInstalled
	}
	return p.fetchNextPayments(ctx, req)
}

func (p Plugin) FetchNextOthers(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	if p.client == nil {
		return models.FetchNextOthersResponse{}, plugins.ErrNotYetInstalled
	}

	switch req.Name {
	case fetchProfileName:
		return p.fetchNextProfiles(ctx, req)
	default:
		return models.FetchNextOthersResponse{}, plugins.ErrNotImplemented
	}
}

func (p Plugin) CreateBankAccount(ctx context.Context, req models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	return models.CreateBankAccountResponse{}, plugins.ErrNotImplemented
}

func (p Plugin) CreateWebhooks(ctx context.Context, req models.CreateWebhooksRequest) (models.CreateWebhooksResponse, error) {
	if p.client == nil {
		return models.CreateWebhooksResponse{}, plugins.ErrNotYetInstalled
	}
	return p.createWebhooks(ctx, req)
}

func (p Plugin) TranslateWebhook(ctx context.Context, req models.TranslateWebhookRequest) (models.TranslateWebhookResponse, error) {
	if p.client == nil {
		return models.TranslateWebhookResponse{}, plugins.ErrNotYetInstalled
	}

	testNotif, ok := req.Webhook.Headers["X-Test-Notification"]
	if ok && len(testNotif) > 0 {
		if testNotif[0] == "true" {
			return models.TranslateWebhookResponse{}, nil
		}
	}

	v, ok := req.Webhook.Headers["X-Delivery-Id"]
	if !ok || len(v) == 0 {
		return models.TranslateWebhookResponse{}, errors.New("missing X-Delivery-Id header")
	}

	signatures, ok := req.Webhook.Headers["X-Signature-Sha256"]
	if !ok || len(signatures) == 0 {
		return models.TranslateWebhookResponse{}, errors.New("missing X-Signature-Sha256 header")
	}

	err := p.verifySignature(req.Webhook.Body, signatures[0])
	if err != nil {
		return models.TranslateWebhookResponse{}, err
	}

	config, ok := webhookConfigs[req.Name]
	if !ok {
		return models.TranslateWebhookResponse{}, errors.New("unknown webhook name")
	}

	res, err := config.fn(ctx, req)
	if err != nil {
		return models.TranslateWebhookResponse{}, err
	}

	res.IdempotencyKey = v[0]

	return models.TranslateWebhookResponse{
		Responses: []models.WebhookResponse{res},
	}, nil
}

var _ models.Plugin = &Plugin{}
