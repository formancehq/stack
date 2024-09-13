package plugins

import (
	"context"
	"errors"

	"github.com/formancehq/payments/internal/connectors/grpc"
	"github.com/formancehq/payments/internal/connectors/grpc/proto/services"
	"github.com/formancehq/payments/internal/models"
)

type impl struct {
	pluginClient grpc.PSP
}

func (i *impl) Install(ctx context.Context, req models.InstallRequest) (models.InstallResponse, error) {
	resp, err := i.pluginClient.Install(ctx, &services.InstallRequest{
		Config: req.Config,
	})
	if err != nil {
		return models.InstallResponse{}, err
	}

	capabilities := make([]models.Capability, 0, len(resp.Capabilities))
	for _, capability := range resp.Capabilities {
		capabilities = append(capabilities, models.Capability(capability))
	}

	webhooksConfigs := make([]models.PSPWebhookConfig, 0, len(resp.WebhooksConfigs))
	for _, webhook := range resp.WebhooksConfigs {
		webhooksConfigs = append(webhooksConfigs, models.PSPWebhookConfig{
			Name:    webhook.Name,
			URLPath: webhook.UrlPath,
		})
	}

	return models.InstallResponse{
		Capabilities:    capabilities,
		Workflow:        grpc.TranslateProtoWorkflow(resp.Workflow),
		WebhooksConfigs: webhooksConfigs,
	}, nil
}

func (i *impl) Uninstall(ctx context.Context, req models.UninstallRequest) (models.UninstallResponse, error) {
	_, err := i.pluginClient.Uninstall(ctx, &services.UninstallRequest{
		ConnectorId: req.ConnectorID,
	})
	if err != nil {
		return models.UninstallResponse{}, err
	}

	return models.UninstallResponse{}, nil
}

func (i *impl) FetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	resp, err := i.pluginClient.FetchNextAccounts(ctx, &services.FetchNextAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextAccountsResponse{}, err
	}

	accounts := make([]models.PSPAccount, 0, len(resp.Accounts))
	for _, account := range resp.Accounts {
		accounts = append(accounts, grpc.TranslateProtoAccount(account))
	}

	return models.FetchNextAccountsResponse{
		Accounts: accounts,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	resp, err := i.pluginClient.FetchNextBalances(ctx, &services.FetchNextBalancesRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	balances := make([]models.PSPBalance, 0, len(resp.Balances))
	for _, balance := range resp.Balances {
		b, err := grpc.TranslateProtoBalance(balance)
		if err != nil {
			return models.FetchNextBalancesResponse{}, err
		}
		balances = append(balances, b)
	}

	return models.FetchNextBalancesResponse{
		Balances: balances,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	resp, err := i.pluginClient.FetchNextExternalAccounts(ctx, &services.FetchNextExternalAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}

	externalAccounts := make([]models.PSPAccount, 0, len(resp.Accounts))
	for _, account := range resp.Accounts {
		externalAccounts = append(externalAccounts, grpc.TranslateProtoAccount(account))
	}

	return models.FetchNextExternalAccountsResponse{
		ExternalAccounts: externalAccounts,
		NewState:         resp.NewState,
		HasMore:          resp.HasMore,
	}, nil
}

func (i *impl) FetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	resp, err := i.pluginClient.FetchNextPayments(ctx, &services.FetchNextPaymentsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}

	payments := make([]models.PSPPayment, 0, len(resp.Payments))
	for _, payment := range resp.Payments {
		p, err := grpc.TranslateProtoPayment(payment)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}
		payments = append(payments, p)
	}

	return models.FetchNextPaymentsResponse{
		Payments: payments,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextOthers(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	resp, err := i.pluginClient.FetchNextOthers(ctx, &services.FetchNextOthersRequest{
		Name:        req.Name,
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	others := make([]models.PSPOther, 0, len(resp.Others))
	for _, other := range resp.Others {
		others = append(others, models.PSPOther{
			ID:    other.Id,
			Other: other.Other,
		})
	}

	return models.FetchNextOthersResponse{
		Others:   others,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) CreateBankAccount(ctx context.Context, req models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	resp, err := i.pluginClient.CreateBankAccount(ctx, &services.CreateBankAccountRequest{
		BankAccount: grpc.TranslateBankAccount(req.BankAccount),
	})
	if err != nil {
		return models.CreateBankAccountResponse{}, err
	}

	return models.CreateBankAccountResponse{
		RelatedAccount: grpc.TranslateProtoAccount(resp.RelatedAccount),
	}, nil
}

func (i *impl) CreateWebhooks(ctx context.Context, req models.CreateWebhooksRequest) (models.CreateWebhooksResponse, error) {
	resp, err := i.pluginClient.CreateWebhooks(ctx, &services.CreateWebhooksRequest{
		ConnectorId: req.ConnectorID,
		FromPayload: req.FromPayload,
	})
	if err != nil {
		return models.CreateWebhooksResponse{}, err
	}

	others := make([]models.PSPOther, 0, len(resp.Others))
	for _, other := range resp.Others {
		others = append(others, models.PSPOther{
			ID:    other.Id,
			Other: other.Other,
		})
	}

	return models.CreateWebhooksResponse{
		Others: others,
	}, nil
}

func (i *impl) TranslateWebhook(ctx context.Context, req models.TranslateWebhookRequest) (models.TranslateWebhookResponse, error) {
	resp, err := i.pluginClient.TranslateWebhook(ctx, &services.TranslateWebhookRequest{
		Name:    req.Name,
		Webhook: grpc.TranslateWebhook(req.Webhook),
	})
	if err != nil {
		return models.TranslateWebhookResponse{}, err
	}

	responses := make([]models.WebhookResponse, 0, len(resp.Responses))
	for _, response := range resp.Responses {
		r := models.WebhookResponse{
			IdempotencyKey: response.IdempotencyKey,
		}

		switch v := response.Translated.(type) {
		case *services.TranslateWebhookResponse_Response_Payment:
			p, err := grpc.TranslateProtoPayment(v.Payment)
			if err != nil {
				return models.TranslateWebhookResponse{}, err
			}
			r.Payment = &p
		case *services.TranslateWebhookResponse_Response_Account:
			a := grpc.TranslateProtoAccount(v.Account)
			r.Account = &a
		case *services.TranslateWebhookResponse_Response_ExternalAccount:
			a := grpc.TranslateProtoAccount(v.ExternalAccount)
			r.ExternalAccount = &a
		default:
			return models.TranslateWebhookResponse{}, errors.New("unknown translated webhook type")
		}

		responses = append(responses, r)
	}

	return models.TranslateWebhookResponse{
		Responses: responses,
	}, nil
}

var _ models.Plugin = &impl{}
