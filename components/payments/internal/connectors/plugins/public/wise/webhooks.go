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
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/wise/client"
	"github.com/formancehq/payments/internal/models"
)

type webhookConfig struct {
	triggerOn string
	urlPath   string
	fn        func(context.Context, models.TranslateWebhookRequest) (models.WebhookResponse, error)
	version   string
}

var webhookConfigs map[string]webhookConfig

func (p Plugin) createWebhooks(ctx context.Context, req models.CreateWebhooksRequest) (models.CreateWebhooksResponse, error) {
	var from client.Profile
	if req.FromPayload == nil {
		return models.CreateWebhooksResponse{}, errors.New("missing from payload when creating webhooks")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.CreateWebhooksResponse{}, err
	}

	stackPublicURL := os.Getenv("STACK_PUBLIC_URL")
	if stackPublicURL == "" {
		err := errors.New("STACK_PUBLIC_URL is not set")
		return models.CreateWebhooksResponse{}, err
	}

	webhookURL := fmt.Sprintf("%s/api/payments/v3/connectors/webhooks/%s", stackPublicURL, req.ConnectorID)
	others := make([]models.PSPOther, 0, len(webhookConfigs))
	for name, config := range webhookConfigs {
		url := fmt.Sprintf("%s%s", webhookURL, config.urlPath)
		resp, err := p.client.CreateWebhook(ctx, from.ID, name, config.triggerOn, url, config.version)
		if err != nil {
			return models.CreateWebhooksResponse{}, err
		}

		raw, err := json.Marshal(resp)
		if err != nil {
			return models.CreateWebhooksResponse{}, err
		}

		others = append(others, models.PSPOther{
			ID:    resp.ID,
			Other: raw,
		})
	}

	return models.CreateWebhooksResponse{
		Others: others,
	}, nil
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

func (p Plugin) translateBalanceUpdateWebhook(ctx context.Context, req models.TranslateWebhookRequest) (models.WebhookResponse, error) {
	update, err := p.client.TranslateBalanceUpdateWebhook(ctx, req.Webhook.Body)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	raw, err := json.Marshal(update)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	occuredAt, err := time.Parse(time.RFC3339, update.Data.OccurredAt)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to parse created time: %w", err)
	}

	var paymentType models.PaymentType
	if update.Data.TransactionType == "credit" {
		paymentType = models.PAYMENT_TYPE_PAYIN
	} else {
		paymentType = models.PAYMENT_TYPE_PAYOUT
	}

	precision, ok := supportedCurrenciesWithDecimal[update.Data.Currency]
	if !ok {
		return models.WebhookResponse{}, nil
	}

	amount, err := currency.GetAmountWithPrecisionFromString(update.Data.Amount.String(), precision)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	payment := models.PSPPayment{
		Reference: update.Data.TransferReference,
		CreatedAt: occuredAt,
		Type:      paymentType,
		Amount:    amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, update.Data.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    models.PAYMENT_STATUS_SUCCEEDED,
		Raw:       raw,
	}

	switch paymentType {
	case models.PAYMENT_TYPE_PAYIN:
		payment.SourceAccountReference = pointer.For(fmt.Sprintf("%d", update.Data.BalanceID))
	case models.PAYMENT_TYPE_PAYOUT:
		payment.DestinationAccountReference = pointer.For(fmt.Sprintf("%d", update.Data.BalanceID))
	}

	return models.WebhookResponse{
		Payment: &payment,
	}, nil
}

func (p Plugin) verifySignature(body []byte, signature string) error {
	msgHash := sha256.New()
	_, err := msgHash.Write(body)
	if err != nil {
		return err
	}
	msgHashSum := msgHash.Sum(nil)

	data, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(p.config.webhookPublicKey, crypto.SHA256, msgHashSum, data)
	if err != nil {
		return err
	}

	return nil
}
