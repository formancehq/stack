package mangopay

import (
	"context"
	"errors"
	"strconv"

	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/connectors/plugins/public/mangopay/client"
	"github.com/formancehq/payments/internal/models"
)

type Plugin struct {
	client *client.Client
}

func (p *Plugin) Install(_ context.Context, req models.InstallRequest) (models.InstallResponse, error) {
	config, err := unmarshalAndValidateConfig(req.Config)
	if err != nil {
		return models.InstallResponse{}, err
	}

	client, err := client.New(config.ClientID, config.APIKey, config.Endpoint)
	if err != nil {
		return models.InstallResponse{}, err
	}
	p.client = client
	p.initWebhookConfig()

	configs := make([]models.PSPWebhookConfig, 0, len(webhookConfigs))
	for name, config := range webhookConfigs {
		configs = append(configs, models.PSPWebhookConfig{
			Name:    string(name),
			URLPath: config.urlPath,
		})
	}

	return models.InstallResponse{
		Capabilities:    capabilities,
		WebhooksConfigs: configs,
		Workflow:        workflow(),
	}, nil
}

func (p Plugin) Uninstall(ctx context.Context, req models.UninstallRequest) (models.UninstallResponse, error) {
	return models.UninstallResponse{}, nil
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
	return p.fetchNextExternalAccounts(ctx, req)
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
	case fetchUsersName:
		return p.fetchNextUsers(ctx, req)
	default:
		return models.FetchNextOthersResponse{}, plugins.ErrNotImplemented
	}
}

func (p Plugin) CreateBankAccount(ctx context.Context, req models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	if p.client == nil {
		return models.CreateBankAccountResponse{}, plugins.ErrNotYetInstalled
	}
	return p.createBankAccount(ctx, req)
}

func (p Plugin) CreateWebhooks(ctx context.Context, req models.CreateWebhooksRequest) (models.CreateWebhooksResponse, error) {
	if p.client == nil {
		return models.CreateWebhooksResponse{}, plugins.ErrNotYetInstalled
	}
	err := p.createWebhooks(ctx, req)
	return models.CreateWebhooksResponse{}, err
}

func (p Plugin) TranslateWebhook(ctx context.Context, req models.TranslateWebhookRequest) (models.TranslateWebhookResponse, error) {
	// Mangopay does not send us the event inside the body, but using
	// URL query.
	eventType, ok := req.Webhook.QueryValues["EventType"]
	if !ok || len(eventType) == 0 {
		return models.TranslateWebhookResponse{}, errors.New("missing EventType query parameter")
	}
	resourceID, ok := req.Webhook.QueryValues["RessourceId"]
	if !ok || len(resourceID) == 0 {
		return models.TranslateWebhookResponse{}, errors.New("missing RessourceId query parameter")
	}
	v, ok := req.Webhook.QueryValues["Date"]
	if !ok || len(v) == 0 {
		return models.TranslateWebhookResponse{}, errors.New("missing Date query parameter")
	}
	date, err := strconv.ParseInt(v[0], 10, 64)
	if err != nil {
		return models.TranslateWebhookResponse{}, errors.New("invalid Date query parameter")
	}

	webhook := client.Webhook{
		ResourceID: resourceID[0],
		Date:       date,
		EventType:  client.EventType(eventType[0]),
	}

	config, ok := webhookConfigs[webhook.EventType]
	if !ok {
		return models.TranslateWebhookResponse{}, errors.New("unsupported webhook event type")
	}

	webhookResponse, err := config.fn(ctx, webhookTranslateRequest{
		req:     req,
		webhook: &webhook,
	})
	if err != nil {
		return models.TranslateWebhookResponse{}, err
	}

	return models.TranslateWebhookResponse{
		Responses: []models.WebhookResponse{webhookResponse},
	}, nil
}

var _ models.Plugin = &Plugin{}
