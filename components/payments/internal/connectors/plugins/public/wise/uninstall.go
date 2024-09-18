package wise

import (
	"context"
	"strings"

	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) uninstall(ctx context.Context, req models.UninstallRequest) (models.UninstallResponse, error) {
	profiles, err := p.client.GetProfiles(ctx)
	if err != nil {
		return models.UninstallResponse{}, err
	}

	for _, profile := range profiles {
		webhooks, err := p.client.ListWebhooksSubscription(ctx, profile.ID)
		if err != nil {
			return models.UninstallResponse{}, err
		}

		for _, webhook := range webhooks {
			if !strings.Contains(webhook.Delivery.URL, req.ConnectorID) {
				continue
			}

			if err := p.client.DeleteWebhooks(ctx, profile.ID, webhook.ID); err != nil {
				return models.UninstallResponse{}, err
			}
		}
	}

	return models.UninstallResponse{}, nil
}
