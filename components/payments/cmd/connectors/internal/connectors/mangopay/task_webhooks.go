package mangopay

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func handleWebhooks(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connectorContext := task.ConnectorContextFromContext(r.Context())
		webhookID := connectors.WebhookIDFromContext(r.Context())
		span := trace.SpanFromContext(r.Context())

		// Mangopay does not send us the event inside the body, but using
		// URL query.
		eventType := r.URL.Query().Get("EventType")
		resourceID := r.URL.Query().Get("RessourceId")

		hook := client.Webhook{
			ResourceID: resourceID,
			EventType:  client.EventType(eventType),
		}

		body, err := json.Marshal(hook)
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}

		if err := store.UpdateWebhookRequestBody(r.Context(), webhookID, body); err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}

		detachedCtx, _ := contextutil.DetachedWithTimeout(r.Context(), 30*time.Second)
		taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:      "handle webhook",
			Key:       taskNameHandleWebhook,
			WebhookID: webhookID,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}

		err = connectorContext.Scheduler().Schedule(detachedCtx, taskDescriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func taskCreateWebhooks(c *client.Client) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		storageReader storage.Reader,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"mangopay.taskCreateWebhooks",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		stackPublicURL := os.Getenv("STACK_PUBLIC_URL")
		if stackPublicURL == "" {
			err := errors.New("STACK_PUBLIC_URL is not set")
			otel.RecordError(span, err)
			return err
		}

		webhookURL := fmt.Sprintf("%s/api/payments/connectors/webhooks/mangopay/%s/", stackPublicURL, &connectorID)
		logger.Infof("creating webhook for mangopay with url %s", webhookURL)

		alreadyExistingHooks, err := c.ListAllHooks(ctx)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		activeHooks := make(map[client.EventType]*client.Hook)
		for _, hook := range alreadyExistingHooks {
			// Mangopay allows only one active hook per event type.
			if hook.Validity == "VALID" {
				activeHooks[hook.EventType] = hook
			}
		}

		for _, eventType := range client.AllEventTypes {
			if v, ok := activeHooks[eventType]; ok {
				// Already created, continue

				if v.URL != webhookURL {
					// If the URL is different, update it
					err := c.UpdateHook(ctx, v.ID, webhookURL)
					if err != nil {
						otel.RecordError(span, err)
						return err
					}

					logger.Infof("updated webhook for mangopay with event type %s", eventType)
				}

				continue
			}

			// Otherwise, create it
			err := c.CreateHook(ctx, eventType, webhookURL)
			if err != nil {
				otel.RecordError(span, err)
				return err
			}

			logger.Infof("created webhook for mangopay with event type %s", eventType)
		}

		return nil
	}
}

func taskHandleWebhooks(c *client.Client, webhookID uuid.UUID) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		storageReader storage.Reader,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"mangopay.taskHandleWebhooks",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("webhookID", webhookID.String()),
		)
		defer span.End()

		w, err := storageReader.GetWebhook(ctx, webhookID)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		webhook, err := c.UnmarshalWebhooks((string(w.RequestBody)))
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		switch webhook.EventType {
		case client.EventTypeTransferNormalCreated,
			client.EventTypeTransferNormalFailed,
			client.EventTypeTransferNormalSucceeded:
			logger.WithField("webhook", webhook).Info("handling transfer webhook")
			return handleTransfer(
				ctx,
				c,
				connectorID,
				webhook,
				ingester,
			)

		case client.EventTypePayoutNormalCreated,
			client.EventTypePayoutNormalFailed,
			client.EventTypePayoutNormalSucceeded,
			client.EventTypePayoutInstantFailed,
			client.EventTypePayoutInstantSucceeded:
			logger.WithField("webhook", webhook).Info("handling payout webhook")
			return handlePayout(
				ctx,
				c,
				connectorID,
				webhook,
				ingester,
			)

		case client.EventTypePayinNormalCreated,
			client.EventTypePayinNormalSucceeded,
			client.EventTypePayinNormalFailed:
			logger.WithField("webhook", webhook).Info("handling payin webhook")
			return handlePayIn(
				ctx,
				c,
				connectorID,
				webhook,
				ingester,
			)

		case client.EventTypeTransferRefundFailed,
			client.EventTypeTransferRefundSucceeded,
			client.EventTypePayOutRefundFailed,
			client.EventTypePayOutRefundSucceeded,
			client.EventTypePayinRefundFailed,
			client.EventTypePayinRefundSucceeded:
			logger.WithField("webhook", webhook).Info("handling refunds webhook")
			return handleRefunds(
				ctx,
				c,
				connectorID,
				webhook,
				ingester,
				storageReader,
			)

		case client.EventTypeTransferRefundCreated,
			client.EventTypePayOutRefundCreated,
			client.EventTypePayinRefundCreated:
			// NOTE: we don't handle these events, as we are only interested in
			// the refund successed or failures.

		default:
			// ignore unknown events
			logger.Errorf("unknown event type: %s", webhook.EventType)
			return nil
		}

		return nil
	}
}

func handleTransfer(
	ctx context.Context,
	c *client.Client,
	connectorID models.ConnectorID,
	webhook *client.Webhook,
	ingester ingestion.Ingester,
) error {
	transfer, err := c.GetWalletTransfer(ctx, webhook.ResourceID)
	if err != nil {
		return err
	}

	if err := fetchWallet(ctx, c, connectorID, transfer.CreditedWalletID, ingester); err != nil {
		return err
	}

	if err := fetchWallet(ctx, c, connectorID, transfer.DebitedWalletID, ingester); err != nil {
		return err
	}

	var amount big.Int
	_, ok := amount.SetString(transfer.DebitedFunds.Amount.String(), 10)
	if !ok {
		return fmt.Errorf("failed to parse amount %s", transfer.DebitedFunds.Amount.String())
	}
	paymentStatus := matchPaymentStatus(transfer.Status)
	raw, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer: %w", err)
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: transfer.ID,
				Type:      models.PaymentTypeTransfer,
			},
			ConnectorID: connectorID,
		},
		CreatedAt:     time.Unix(transfer.CreationDate, 0),
		Reference:     transfer.ID,
		Amount:        &amount,
		InitialAmount: &amount,
		ConnectorID:   connectorID,
		Type:          models.PaymentTypeTransfer,
		Status:        paymentStatus,
		Scheme:        models.PaymentSchemeOther,
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, transfer.DebitedFunds.Currency),
		RawData:       raw,
	}

	if transfer.DebitedWalletID != "" {
		payment.SourceAccountID = &models.AccountID{
			Reference:   transfer.DebitedWalletID,
			ConnectorID: connectorID,
		}
	}

	if transfer.CreditedWalletID != "" {
		payment.DestinationAccountID = &models.AccountID{
			Reference:   transfer.CreditedWalletID,
			ConnectorID: connectorID,
		}
	}

	err = ingester.IngestPayments(ctx, []ingestion.PaymentBatchElement{{
		Payment: payment,
	}})
	if err != nil {
		return err
	}

	return nil
}

func handlePayIn(
	ctx context.Context,
	c *client.Client,
	connectorID models.ConnectorID,
	webhook *client.Webhook,
	ingester ingestion.Ingester,
) error {
	payin, err := c.GetPayin(ctx, webhook.ResourceID)
	if err != nil {
		return err
	}

	// In case of a payin, there is no debited wallet id, so we can only
	// fetch the credited wallet.
	if err := fetchWallet(ctx, c, connectorID, payin.CreditedWalletID, ingester); err != nil {
		return err
	}

	var amount big.Int
	_, ok := amount.SetString(payin.DebitedFunds.Amount.String(), 10)
	if !ok {
		return fmt.Errorf("failed to parse amount %s", payin.DebitedFunds.Amount.String())
	}
	paymentStatus := matchPaymentStatus(payin.Status)
	raw, err := json.Marshal(payin)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer: %w", err)
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: payin.ID,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		},
		CreatedAt:     time.Unix(payin.CreationDate, 0),
		Reference:     payin.ID,
		Amount:        &amount,
		InitialAmount: &amount,
		ConnectorID:   connectorID,
		Type:          models.PaymentTypePayIn,
		Status:        paymentStatus,
		Scheme:        models.PaymentSchemeOther,
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, payin.DebitedFunds.Currency),
		RawData:       raw,
	}

	if payin.CreditedWalletID != "" {
		payment.DestinationAccountID = &models.AccountID{
			Reference:   payin.CreditedWalletID,
			ConnectorID: connectorID,
		}
	}

	err = ingester.IngestPayments(ctx, []ingestion.PaymentBatchElement{{
		Payment: payment,
	}})
	if err != nil {
		return err
	}

	return nil
}

func handlePayout(
	ctx context.Context,
	c *client.Client,
	connectorID models.ConnectorID,
	webhook *client.Webhook,
	ingester ingestion.Ingester,
) error {
	payout, err := c.GetPayout(ctx, webhook.ResourceID)
	if err != nil {
		return err
	}

	// In case of a payout, there is no credited wallet id, so we can only
	// fetch the debited wallet.
	if err := fetchWallet(ctx, c, connectorID, payout.DebitedWalletID, ingester); err != nil {
		return err
	}

	var amount big.Int
	_, ok := amount.SetString(payout.DebitedFunds.Amount.String(), 10)
	if !ok {
		return fmt.Errorf("failed to parse amount %s", payout.DebitedFunds.Amount.String())
	}
	paymentStatus := matchPaymentStatus(payout.Status)
	raw, err := json.Marshal(payout)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer: %w", err)
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: payout.ID,
				Type:      models.PaymentTypePayOut,
			},
			ConnectorID: connectorID,
		},
		CreatedAt:     time.Unix(payout.CreationDate, 0),
		Reference:     payout.ID,
		Amount:        &amount,
		InitialAmount: &amount,
		ConnectorID:   connectorID,
		Type:          models.PaymentTypePayOut,
		Status:        paymentStatus,
		Scheme:        models.PaymentSchemeOther,
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, payout.DebitedFunds.Currency),
		RawData:       raw,
	}

	if payout.DebitedWalletID != "" {
		payment.DestinationAccountID = &models.AccountID{
			Reference:   payout.DebitedWalletID,
			ConnectorID: connectorID,
		}
	}

	err = ingester.IngestPayments(ctx, []ingestion.PaymentBatchElement{{
		Payment: payment,
	}})
	if err != nil {
		return err
	}

	return nil
}

func handleRefunds(
	ctx context.Context,
	c *client.Client,
	connectorID models.ConnectorID,
	webhook *client.Webhook,
	ingester ingestion.Ingester,
	storageReader storage.Reader,
) error {
	refund, err := c.GetRefund(ctx, webhook.ResourceID)
	if err != nil {
		return err
	}

	var amountRefunded big.Int
	_, ok := amountRefunded.SetString(refund.DebitedFunds.Amount.String(), 10)
	if !ok {
		return fmt.Errorf("failed to parse amount %s", refund.DebitedFunds.Amount.String())
	}
	paymentType := matchPaymentType(refund.InitialTransactionType)

	var payment *models.Payment
	if webhook.EventType == client.EventTypePayOutRefundSucceeded {
		payment, err = storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: refund.InitialTransactionID,
				Type:      paymentType,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Amount = payment.Amount.Sub(payment.Amount, &amountRefunded)
	}

	paymentStatus := models.PaymentStatusRefundedFailure
	if webhook.EventType == client.EventTypePayOutRefundSucceeded {
		paymentStatus = models.PaymentStatusRefunded
	}

	raw, err := json.Marshal(refund)
	if err != nil {
		return fmt.Errorf("failed to marshal refund: %w", err)
	}

	adjustment := &models.PaymentAdjustment{
		PaymentID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: refund.InitialTransactionID,
				Type:      paymentType,
			},
			ConnectorID: connectorID,
		},
		CreatedAt: time.Unix(refund.CreationDate, 0),
		Reference: refund.ID,
		Amount:    &amountRefunded,
		Status:    paymentStatus,
		RawData:   raw,
	}

	if err := ingester.IngestPayments(ctx, []ingestion.PaymentBatchElement{{
		Payment:    payment,
		Adjustment: adjustment,
	}}); err != nil {
		return err
	}

	return nil
}

func fetchWallet(
	ctx context.Context,
	c *client.Client,
	connectorID models.ConnectorID,
	walletID string,
	ingester ingestion.Ingester,
) error {
	if walletID == "" {
		return nil
	}

	wallet, err := c.GetWallet(ctx, walletID)
	if err != nil {
		return err
	}

	raw, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	userID := ""
	if len(wallet.Owners) > 0 {
		userID = wallet.Owners[0]
	}

	account := &models.Account{
		ID: models.AccountID{
			Reference:   wallet.ID,
			ConnectorID: connectorID,
		},
		CreatedAt:    time.Unix(wallet.CreationDate, 0),
		Reference:    wallet.ID,
		ConnectorID:  connectorID,
		DefaultAsset: currency.FormatAsset(supportedCurrenciesWithDecimal, wallet.Currency),
		AccountName:  wallet.Description,
		// Wallets are internal accounts on our side, since we
		// can have their balances.
		Type: models.AccountTypeInternal,
		Metadata: map[string]string{
			"user_id": userID,
		},
		RawData: raw,
	}

	var amount big.Int
	_, ok := amount.SetString(wallet.Balance.Amount.String(), 10)
	if !ok {
		return fmt.Errorf("failed to parse amount: %s", wallet.Balance.Amount.String())
	}

	now := time.Now()
	balance := &models.Balance{
		AccountID: models.AccountID{
			Reference:   wallet.ID,
			ConnectorID: connectorID,
		},
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, wallet.Balance.Currency),
		Balance:       &amount,
		CreatedAt:     now,
		LastUpdatedAt: now,
		ConnectorID:   connectorID,
	}

	if err := ingester.IngestAccounts(ctx, []*models.Account{account}); err != nil {
		return err
	}

	if err := ingester.IngestBalances(ctx, []*models.Balance{balance}, false); err != nil {
		return err
	}

	return nil
}
