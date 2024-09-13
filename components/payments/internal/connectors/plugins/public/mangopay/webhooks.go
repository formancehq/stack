package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/mangopay/client"
	"github.com/formancehq/payments/internal/models"
)

type webhookTranslateRequest struct {
	req     models.TranslateWebhookRequest
	webhook *client.Webhook
}

type webhookConfig struct {
	urlPath string
	fn      func(context.Context, webhookTranslateRequest) (models.WebhookResponse, error)
}

var webhookConfigs map[client.EventType]webhookConfig

func (p Plugin) initWebhookConfig() {
	webhookConfigs = map[client.EventType]webhookConfig{
		client.EventTypeTransferNormalCreated: {
			urlPath: "/transfer/created",
			fn:      p.translateTransfer,
		},
		client.EventTypeTransferNormalFailed: {
			urlPath: "/transfer/failed",
			fn:      p.translateTransfer,
		},
		client.EventTypeTransferNormalSucceeded: {
			urlPath: "/transfer/succeeded",
			fn:      p.translateTransfer,
		},

		client.EventTypePayoutNormalCreated: {
			urlPath: "/payout/normal/created",
			fn:      p.translatePayout,
		},
		client.EventTypePayoutNormalFailed: {
			urlPath: "/payout/normal/failed",
			fn:      p.translatePayout,
		},
		client.EventTypePayoutNormalSucceeded: {
			urlPath: "/payout/normal/succeeded",
			fn:      p.translatePayout,
		},
		client.EventTypePayoutInstantFailed: {
			urlPath: "/payout/instant/failed",
			fn:      p.translatePayout,
		},
		client.EventTypePayoutInstantSucceeded: {
			urlPath: "/payout/instant/succeeded",
			fn:      p.translatePayout,
		},

		client.EventTypePayinNormalCreated: {
			urlPath: "/payin/normal/created",
			fn:      p.translatePayin,
		},
		client.EventTypePayinNormalSucceeded: {
			urlPath: "/payin/normal/succeeded",
			fn:      p.translatePayin,
		},
		client.EventTypePayinNormalFailed: {
			urlPath: "/payin/normal/failed",
			fn:      p.translatePayin,
		},

		client.EventTypeTransferRefundFailed: {
			urlPath: "/refund/transfer/failed",
			fn:      p.translateRefund,
		},
		client.EventTypeTransferRefundSucceeded: {
			urlPath: "/refund/transfer/succeeded",
			fn:      p.translateRefund,
		},
		client.EventTypePayOutRefundFailed: {
			urlPath: "/refund/payout/failed",
			fn:      p.translateRefund,
		},
		client.EventTypePayOutRefundSucceeded: {
			urlPath: "/refund/payout/succeeded",
			fn:      p.translateRefund,
		},
		client.EventTypePayinRefundFailed: {
			urlPath: "/refund/payin/failed",
			fn:      p.translateRefund,
		},
		client.EventTypePayinRefundSucceeded: {
			urlPath: "/refund/payin/succeeded",
			fn:      p.translateRefund,
		},
	}
}

func (p Plugin) createWebhooks(ctx context.Context, req models.CreateWebhooksRequest) error {
	stackPublicURL := os.Getenv("STACK_PUBLIC_URL")
	if stackPublicURL == "" {
		err := errors.New("STACK_PUBLIC_URL is not set")
		return err
	}

	activeHooks, err := p.getActiveHooks(ctx)
	if err != nil {
		return err
	}

	webhookURL := fmt.Sprintf("%s/api/payments/v3/connectors/webhooks/%s", stackPublicURL, req.ConnectorID)
	for eventType, config := range webhookConfigs {
		url, err := url.JoinPath(webhookURL, config.urlPath)
		if err != nil {
			return err
		}

		if v, ok := activeHooks[eventType]; ok {
			// Already created, continue

			if v.URL != webhookURL {
				// If the URL is different, update it
				err := p.client.UpdateHook(ctx, v.ID, url)
				if err != nil {
					return err
				}
			}

			continue
		}

		// Otherwise, create it
		err = p.client.CreateHook(ctx, eventType, url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Plugin) getActiveHooks(ctx context.Context) (map[client.EventType]*client.Hook, error) {
	alreadyExistingHooks, err := p.client.ListAllHooks(ctx)
	if err != nil {
		return nil, err
	}

	activeHooks := make(map[client.EventType]*client.Hook)
	for _, hook := range alreadyExistingHooks {
		// Mangopay allows only one active hook per event type.
		if hook.Validity == "VALID" {
			activeHooks[hook.EventType] = hook
		}
	}

	return activeHooks, nil
}

func (p Plugin) translateTransfer(ctx context.Context, req webhookTranslateRequest) (models.WebhookResponse, error) {
	transfer, err := p.client.GetWalletTransfer(ctx, req.webhook.ResourceID)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	raw, err := json.Marshal(transfer)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal transfer: %w", err)
	}

	paymentStatus := matchPaymentStatus(transfer.Status)

	var amount big.Int
	_, ok := amount.SetString(transfer.DebitedFunds.Amount.String(), 10)
	if !ok {
		return models.WebhookResponse{}, fmt.Errorf("failed to parse amount %s", transfer.DebitedFunds.Amount.String())
	}

	payment := models.PSPPayment{
		Reference: transfer.ID,
		CreatedAt: time.Unix(transfer.CreationDate, 0),
		Type:      models.PAYMENT_TYPE_TRANSFER,
		Amount:    &amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, transfer.DebitedFunds.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    paymentStatus,
		Raw:       raw,
	}

	if transfer.DebitedWalletID != "" {
		payment.SourceAccountReference = &transfer.DebitedWalletID
	}

	if transfer.CreditedWalletID != "" {
		payment.DestinationAccountReference = &transfer.CreditedWalletID
	}

	return models.WebhookResponse{
		IdempotencyKey: fmt.Sprintf("%s-%s-%d", req.webhook.ResourceID, string(req.webhook.EventType), req.webhook.Date),
		Payment:        &payment,
	}, nil
}

func (p Plugin) translatePayout(ctx context.Context, req webhookTranslateRequest) (models.WebhookResponse, error) {
	payout, err := p.client.GetPayout(ctx, req.webhook.ResourceID)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	raw, err := json.Marshal(payout)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal transfer: %w", err)
	}

	paymentStatus := matchPaymentStatus(payout.Status)

	var amount big.Int
	_, ok := amount.SetString(payout.DebitedFunds.Amount.String(), 10)
	if !ok {
		return models.WebhookResponse{}, fmt.Errorf("failed to parse amount %s", payout.DebitedFunds.Amount.String())
	}

	payment := models.PSPPayment{
		Reference: payout.ID,
		CreatedAt: time.Unix(payout.CreationDate, 0),
		Type:      models.PAYMENT_TYPE_PAYOUT,
		Amount:    &amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, payout.DebitedFunds.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    paymentStatus,
		Raw:       raw,
	}

	if payout.DebitedWalletID != "" {
		payment.DestinationAccountReference = &payout.DebitedWalletID
	}

	return models.WebhookResponse{
		IdempotencyKey: fmt.Sprintf("%s-%s-%d", req.webhook.ResourceID, string(req.webhook.EventType), req.webhook.Date),
		Payment:        &payment,
	}, nil
}

func (p Plugin) translatePayin(ctx context.Context, req webhookTranslateRequest) (models.WebhookResponse, error) {
	payin, err := p.client.GetPayin(ctx, req.webhook.ResourceID)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	raw, err := json.Marshal(payin)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal transfer: %w", err)
	}

	paymentStatus := matchPaymentStatus(payin.Status)

	var amount big.Int
	_, ok := amount.SetString(payin.DebitedFunds.Amount.String(), 10)
	if !ok {
		return models.WebhookResponse{}, fmt.Errorf("failed to parse amount %s", payin.DebitedFunds.Amount.String())
	}

	payment := models.PSPPayment{
		Reference: payin.ID,
		CreatedAt: time.Unix(payin.CreationDate, 0),
		Type:      models.PAYMENT_TYPE_PAYIN,
		Amount:    &amount,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, payin.DebitedFunds.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    paymentStatus,
		Raw:       raw,
	}

	if payin.CreditedWalletID != "" {
		payment.DestinationAccountReference = &payin.CreditedWalletID
	}

	return models.WebhookResponse{
		IdempotencyKey: fmt.Sprintf("%s-%s-%d", req.webhook.ResourceID, string(req.webhook.EventType), req.webhook.Date),
		Payment:        &payment,
	}, nil
}

func (p Plugin) translateRefund(ctx context.Context, req webhookTranslateRequest) (models.WebhookResponse, error) {
	refund, err := p.client.GetRefund(ctx, req.webhook.ResourceID)
	if err != nil {
		return models.WebhookResponse{}, err
	}

	raw, err := json.Marshal(refund)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal transfer: %w", err)
	}

	paymentType := matchPaymentType(refund.InitialTransactionType)

	var amountRefunded big.Int
	_, ok := amountRefunded.SetString(refund.DebitedFunds.Amount.String(), 10)
	if !ok {
		return models.WebhookResponse{}, fmt.Errorf("failed to parse amount %s", refund.DebitedFunds.Amount.String())
	}

	payment := models.PSPPayment{
		Reference: refund.InitialTransactionID,
		CreatedAt: time.Unix(refund.CreationDate, 0),
		Type:      paymentType,
		Amount:    &amountRefunded,
		Asset:     currency.FormatAsset(supportedCurrenciesWithDecimal, refund.DebitedFunds.Currency),
		Scheme:    models.PAYMENT_SCHEME_OTHER,
		Status:    models.PAYMENT_STATUS_REFUNDED,
		Raw:       raw,
	}

	return models.WebhookResponse{
		IdempotencyKey: fmt.Sprintf("%s-%s-%d", req.webhook.ResourceID, string(req.webhook.EventType), req.webhook.Date),
		Payment:        &payment,
	}, nil
}
