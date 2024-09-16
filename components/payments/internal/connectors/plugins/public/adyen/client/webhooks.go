package client

import (
	"context"
	"fmt"

	"github.com/adyen/adyen-go-api-library/v7/src/management"
	"github.com/adyen/adyen-go-api-library/v7/src/webhook"
	"github.com/formancehq/stack/libs/go-libs/pointer"
)

func (c *Client) SearchWebhook(ctx context.Context, connectorID string) error {
	pageSize := 50
	for page := 1; ; page++ {
		webhooks, raw, err := c.client.Management().WebhooksCompanyLevelApi.ListAllWebhooks(
			ctx,
			c.client.Management().WebhooksCompanyLevelApi.ListAllWebhooksInput(c.companyID).PageNumber(int32(page)).PageSize(int32(pageSize)),
		)
		if err != nil {
			return err
		}

		if raw.StatusCode >= 300 {
			return fmt.Errorf("failed to get webhooks: %d", raw.StatusCode)
		}

		if len(webhooks.Data) == 0 {
			break
		}

		for _, webhook := range webhooks.Data {
			if _, ok := c.webhooks[""]; ok {
				// already exists
				continue
			}

			if webhook.Description == nil {
				continue
			}

			if *webhook.Description != connectorID {
				continue
			}

			c.webhooks[merchantID] = webhook
		}

		if len(webhooks.Data) < pageSize {
			break
		}
	}

	return nil
}

func (c *Client) CreateWebhook(ctx context.Context, url string, connectorID string) error {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	if _, ok := c.webhooks[merchantID]; ok {
		return nil
	}

	if err := c.FillExistingWebhooks(ctx, connectorID, merchantID); err != nil {
		return err
	}

	if _, ok := c.webhooks[merchantID]; ok {
		return nil
	}

	req := management.CreateMerchantWebhookRequest{
		Active: true,
		AdditionalSettings: &management.AdditionalSettings{
			Properties: &map[string]bool{
				"hmacSignatureKey": true,
			},
		},
		CommunicationFormat: "json",
		Description:         pointer.For(connectorID),
		SslVersion:          pointer.For("TLSv1.3"),
		Type:                "standard",
		Url:                 url,
	}

	if c.webhookUsername != "" {
		req.Username = pointer.For(c.webhookUsername)
	}

	if c.webhookPassword != "" {
		req.Password = pointer.For(c.webhookPassword)
	}

	webhook, raw, err := c.client.Management().WebhooksMerchantLevelApi.SetUpWebhook(
		ctx,
		c.client.Management().WebhooksMerchantLevelApi.SetUpWebhookInput(merchantID).
			CreateMerchantWebhookRequest(req),
	)
	if err != nil {
		return err
	}

	if raw.StatusCode >= 300 {
		return fmt.Errorf("failed to create webhook: %d", raw.StatusCode)
	}

	c.webhooks[merchantID] = webhook

	return nil
}

func (c *Client) DeleteWebhook() error {
	return nil
}

func (c *Client) TranslateWebhook(req string) (*webhook.Webhook, error) {
	return webhook.HandleRequest(req)
}
