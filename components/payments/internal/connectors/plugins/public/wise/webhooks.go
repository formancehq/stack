package wise

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/formancehq/payments/internal/connectors/plugins/public/wise/client"
	"github.com/formancehq/payments/internal/models"
)

type webhookConfig struct {
	triggerOn string
	urlPath   string
	fn        func(context.Context, models.TranslateWebhookRequest) (models.WebhookResponse, error)
}

var webhookConfigs map[string]webhookConfig

func (p Plugin) createWebhooks(ctx context.Context, req models.CreateWebhooksRequest) error {
	var from client.Profile
	if req.FromPayload == nil {
		return errors.New("missing from payload when creating webhooks")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return err
	}

	stackPublicURL := os.Getenv("STACK_PUBLIC_URL")
	if stackPublicURL == "" {
		err := errors.New("STACK_PUBLIC_URL is not set")
		return err
	}

	webhookURL := fmt.Sprintf("%s/api/payments/v3/connectors/webhooks/%s", stackPublicURL, req.ConnectorID)

	for name, config := range webhookConfigs {
		url := fmt.Sprintf("%s%s", webhookURL, config.urlPath)
		if err := p.client.CreateWebhook(ctx, from.ID, name, config.triggerOn, url); err != nil {
			return err
		}
	}

	return nil
}

func (p Plugin) translateTransferStateChangedWebhook(ctx context.Context, req models.TranslateWebhookRequest) (models.WebhookResponse, error) {
	transfer, err := p.client.TranslateTransferStateChangedWebhook(ctx, req.Webhook.Body)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	payment, err := fromTransferToPayment(transfer)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	return models.WebhookResponse{
		Payment: payment,
	}, nil
}

func (p Plugin) verifySignature(payload []byte, signature string) error {
	msgHash := sha256.New()
	_, err := msgHash.Write(payload)
	if err != nil {
		return err
	}
	msgHashSum := msgHash.Sum(nil)

	data, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	err = rsa.VerifyPSS(p.config.webhookPublicKey, crypto.SHA256, msgHashSum, data, nil)
	if err != nil {
		return err
	}

	return nil
}
