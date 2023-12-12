package client

import "github.com/adyen/adyen-go-api-library/v7/src/webhook"

func (c *Client) CreateWebhookForRequest(req string) (*webhook.Webhook, error) {
	return webhook.HandleRequest(req)
}
