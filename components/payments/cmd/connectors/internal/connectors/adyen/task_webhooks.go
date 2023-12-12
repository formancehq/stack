package adyen

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/adyen/adyen-go-api-library/v7/src/hmacvalidator"
	"github.com/adyen/adyen-go-api-library/v7/src/webhook"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/adyen/client"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
)

func taskHandleWebhook(client *client.Client, webhookID uuid.UUID) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		storageReader storage.Reader,
		ingester ingestion.Ingester,
	) error {
		logger.Info(taskNameMain)

		w, err := storageReader.GetWebhook(ctx, webhookID)
		if err != nil {
			return err
		}

		webhooks, err := client.CreateWebhookForRequest(string(w.RequestBody))
		if err != nil {
			return err
		}

		for _, item := range *webhooks.NotificationItems {
			logger.Infof("received webhook: %s -> %v", item.NotificationRequestItem.EventCode, item)
			logger.Infof("with signature: %s", (*item.NotificationRequestItem.AdditionalData)["hmacSignature"])
			logger.Infof("hmac key: %s", client.HMACKey)

			expectedSign, err := hmacvalidator.CalculateHmac(item.NotificationRequestItem, client.HMACKey)
			if err != nil {
				logger.Errorf("error calculating hmac: %v", err)
			}

			logger.Infof("expected signature: %s", expectedSign)

			if !hmacvalidator.ValidateHmac(item.NotificationRequestItem, client.HMACKey) {
				logger.Errorf("invalid hmac for webhook: %s", item.NotificationRequestItem.EventCode)
				continue
			}

			if err := handleNotificationRequestItem(
				ctx,
				connectorID,
				storageReader,
				ingester,
				item.NotificationRequestItem,
			); err != nil {
				return err
			}
		}

		return nil
	}
}

func handleNotificationRequestItem(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	switch item.EventCode {
	case webhook.EventCodeAuthorisation:
		return handleAuthorisation(ctx, connectorID, ingester, item)
	case webhook.EventCodeAuthorisationAdjustment:
		return handleAuthorisationAdjustment(ctx, connectorID, storageReader, ingester, item)
	case webhook.EventCodeCancellation:
		return handleCancellation(ctx, connectorID, storageReader, ingester, item)
	case webhook.EventCodeCapture:
		return handleCapture(ctx, connectorID, storageReader, ingester, item)
	case webhook.EventCodeCaptureFailed:
		return handleCaptureFailed(ctx, connectorID, storageReader, ingester, item)
	case webhook.EventCodeRefund:
		return handleRefund(ctx, connectorID, storageReader, ingester, item)
	case webhook.EventCodeRefundFailed:
		return handleRefundFailed()
	case webhook.EventCodeRefundedReversed:
		return handleRefundedReversed(ctx, connectorID, storageReader, ingester, item)
	case webhook.EventCodeRefundWithData:
		return handleRefundWithData(ctx, connectorID, storageReader, ingester, item)
	case webhook.EventCodePayoutThirdparty:
		return handlePayoutThirdparty(ctx, connectorID, ingester, item)
	case webhook.EventCodePayoutDecline:
		return handlePayoutDecline(ctx, connectorID, ingester, item)
	case webhook.EventCodePayoutExpire:
		return handlePayoutExpire(ctx, connectorID, ingester, item)
	}

	return nil
}

func handleAuthorisation(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	raw, err := json.Marshal(item)
	if err != nil {
		return err
	}

	status := models.PaymentStatusPending
	if item.Success == "false" {
		status = models.PaymentStatusFailed
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.PspReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		CreatedAt:     *item.EventDate,
		Reference:     item.PspReference,
		Amount:        big.NewInt(item.Amount.Value),
		InitialAmount: big.NewInt(item.Amount.Value),
		Type:          models.PaymentTypePayIn,
		Status:        status,
		Scheme:        parseScheme(item.PaymentMethod),
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, item.Amount.Currency),
		RawData:       raw,
		DestinationAccountID: &models.AccountID{
			Reference:   item.MerchantAccountCode,
			ConnectorID: connectorID,
		},
	}

	if err := ingester.IngestPayments(
		ctx,
		connectorID,
		ingestion.PaymentBatch{{Payment: payment}},
		struct{}{},
	); err != nil {
		return err
	}

	return nil
}

func handleAuthorisationAdjustment(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success == "true" {
		payment, err := storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.OriginalReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Amount = big.NewInt(item.Amount.Value)
		payment.InitialAmount = big.NewInt(item.Amount.Value)

		if err := ingester.IngestPayments(
			ctx,
			connectorID,
			ingestion.PaymentBatch{{Payment: payment}},
			struct{}{},
		); err != nil {
			return err
		}
	}

	return nil
}

func handleCancellation(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success == "true" {
		payment, err := storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.OriginalReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Status = models.PaymentStatusCancelled

		if err := ingester.IngestPayments(
			ctx,
			connectorID,
			ingestion.PaymentBatch{{Payment: payment}},
			struct{}{},
		); err != nil {
			return err
		}
	}

	return nil
}

func handleCapture(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success == "true" {
		payment, err := storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.OriginalReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Status = models.PaymentStatusSucceeded
		payment.Amount = big.NewInt(item.Amount.Value)

		if err := ingester.IngestPayments(
			ctx,
			connectorID,
			ingestion.PaymentBatch{{Payment: payment}},
			struct{}{},
		); err != nil {
			return err
		}
	}

	return nil
}

func handleCaptureFailed(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success == "true" {
		payment, err := storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.OriginalReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Status = models.PaymentStatusFailed

		if err := ingester.IngestPayments(
			ctx,
			connectorID,
			ingestion.PaymentBatch{{Payment: payment}},
			struct{}{},
		); err != nil {
			return err
		}
	}

	return nil
}

func handleRefund(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success == "true" {
		payment, err := storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.OriginalReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Amount = payment.Amount.Sub(payment.Amount, big.NewInt(item.Amount.Value))
		if payment.Amount.Cmp(big.NewInt(0)) == 0 {
			payment.Status = models.PaymentStatusRefunded
		}

		if err := ingester.IngestPayments(
			ctx,
			connectorID,
			ingestion.PaymentBatch{{Payment: payment}},
			struct{}{},
		); err != nil {
			return err
		}
	}

	return nil
}

func handleRefundFailed() error {
	// Nothing to do for now (while waiting to enhance the payment adjustment model)
	return nil
}

func handleRefundedReversed(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success == "true" {
		payment, err := storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.PspReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Amount = payment.Amount.Add(payment.Amount, big.NewInt(item.Amount.Value))

		if err := ingester.IngestPayments(
			ctx,
			connectorID,
			ingestion.PaymentBatch{{Payment: payment}},
			struct{}{},
		); err != nil {
			return err
		}
	}

	return nil
}

func handleRefundWithData(
	ctx context.Context,
	connectorID models.ConnectorID,
	storageReader storage.Reader,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success == "true" {
		payment, err := storageReader.GetPayment(ctx, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.OriginalReference,
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		}.String())
		if err != nil {
			return err
		}

		payment.Amount = payment.Amount.Sub(payment.Amount, big.NewInt(item.Amount.Value))
		if payment.Amount.Cmp(big.NewInt(0)) == 0 {
			payment.Status = models.PaymentStatusRefunded
		}

		if err := ingester.IngestPayments(
			ctx,
			connectorID,
			ingestion.PaymentBatch{{Payment: payment}},
			struct{}{},
		); err != nil {
			return err
		}
	}

	return nil
}

func handlePayoutThirdparty(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	raw, err := json.Marshal(item)
	if err != nil {
		return err
	}

	status := models.PaymentStatusSucceeded
	if item.Success == "false" {
		status = models.PaymentStatusFailed
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.PspReference,
				Type:      models.PaymentTypePayOut,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		CreatedAt:     *item.EventDate,
		Reference:     item.PspReference,
		Amount:        big.NewInt(item.Amount.Value),
		InitialAmount: big.NewInt(item.Amount.Value),
		Type:          models.PaymentTypePayOut,
		Status:        status,
		Scheme:        parseScheme(item.PaymentMethod),
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, item.Amount.Currency),
		RawData:       raw,
		SourceAccountID: &models.AccountID{
			Reference:   item.MerchantAccountCode,
			ConnectorID: connectorID,
		},
	}

	if err := ingester.IngestPayments(
		ctx,
		connectorID,
		ingestion.PaymentBatch{{Payment: payment}},
		struct{}{},
	); err != nil {
		return err
	}

	return nil
}

func handlePayoutDecline(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success != "true" {
		return nil
	}

	raw, err := json.Marshal(item)
	if err != nil {
		return err
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.PspReference,
				Type:      models.PaymentTypePayOut,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		CreatedAt:     *item.EventDate,
		Reference:     item.PspReference,
		Amount:        big.NewInt(item.Amount.Value),
		InitialAmount: big.NewInt(item.Amount.Value),
		Type:          models.PaymentTypePayOut,
		Status:        models.PaymentStatusCancelled,
		Scheme:        parseScheme(item.PaymentMethod),
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, item.Amount.Currency),
		RawData:       raw,
		SourceAccountID: &models.AccountID{
			Reference:   item.MerchantAccountCode,
			ConnectorID: connectorID,
		},
	}

	if err := ingester.IngestPayments(
		ctx,
		connectorID,
		ingestion.PaymentBatch{{Payment: payment}},
		struct{}{},
	); err != nil {
		return err
	}

	return nil
}

func handlePayoutExpire(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	item webhook.NotificationRequestItem,
) error {
	if item.Success != "true" {
		return nil
	}

	raw, err := json.Marshal(item)
	if err != nil {
		return err
	}

	payment := &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: item.PspReference,
				Type:      models.PaymentTypePayOut,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		CreatedAt:     *item.EventDate,
		Reference:     item.PspReference,
		Amount:        big.NewInt(item.Amount.Value),
		InitialAmount: big.NewInt(item.Amount.Value),
		Type:          models.PaymentTypePayOut,
		Status:        models.PaymentStatusExpired,
		Scheme:        parseScheme(item.PaymentMethod),
		Asset:         currency.FormatAsset(supportedCurrenciesWithDecimal, item.Amount.Currency),
		RawData:       raw,
		SourceAccountID: &models.AccountID{
			Reference:   item.MerchantAccountCode,
			ConnectorID: connectorID,
		},
	}

	if err := ingester.IngestPayments(
		ctx,
		connectorID,
		ingestion.PaymentBatch{{Payment: payment}},
		struct{}{},
	); err != nil {
		return err
	}

	return nil
}

func parseScheme(scheme string) models.PaymentScheme {
	switch {
	case strings.HasPrefix(scheme, "visa"):
		return models.PaymentSchemeCardVisa
	case strings.HasPrefix(scheme, "electron"):
		return models.PaymentSchemeCardVisa
	case strings.HasPrefix(scheme, "amex"):
		return models.PaymentSchemeCardAmex
	case strings.HasPrefix(scheme, "alipay"):
		return models.PaymentSchemeCardAlipay
	case strings.HasPrefix(scheme, "cup"):
		return models.PaymentSchemeCardCUP
	case strings.HasPrefix(scheme, "discover"):
		return models.PaymentSchemeCardDiscover
	case strings.HasPrefix(scheme, "doku"):
		return models.PaymentSchemeDOKU
	case strings.HasPrefix(scheme, "dragonpay"):
		return models.PaymentSchemeDragonPay
	case strings.HasPrefix(scheme, "jcb"):
		return models.PaymentSchemeCardJCB
	case strings.HasPrefix(scheme, "maestro"):
		return models.PaymentSchemeMaestro
	case strings.HasPrefix(scheme, "mc"):
		return models.PaymentSchemeCardMasterCard
	case strings.HasPrefix(scheme, "molpay"):
		return models.PaymentSchemeMolPay
	case strings.HasPrefix(scheme, "diners"):
		return models.PaymentSchemeCardDiners
	default:
		return models.PaymentSchemeUnknown
	}
}
