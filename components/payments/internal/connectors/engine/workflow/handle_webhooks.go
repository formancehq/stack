package workflow

import (
	"fmt"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type HandleWebhooks struct {
	ConnectorID   models.ConnectorID
	WebhookConfig models.WebhookConfig
	Webhook       models.Webhook
}

func (w Workflow) runHandleWebhooks(
	ctx workflow.Context,
	handleWebhooks HandleWebhooks,
) error {
	err := activities.StorageWebhooksStore(infiniteRetryContext(ctx), handleWebhooks.Webhook)
	if err != nil {
		return errors.Wrap(err, "failed to store webhook")
	}

	resp, err := activities.PluginTranslateWebhook(
		infiniteRetryContext(ctx),
		handleWebhooks.ConnectorID,
		models.TranslateWebhookRequest{
			Name: handleWebhooks.WebhookConfig.Name,
			Webhook: models.PSPWebhook{
				QueryValues: handleWebhooks.Webhook.QueryValues,
				Headers:     handleWebhooks.Webhook.Headers,
				Body:        handleWebhooks.Webhook.Body,
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to translate webhook")
	}

	for _, response := range resp.Responses {
		if err := workflow.ExecuteChildWorkflow(
			workflow.WithChildOptions(
				ctx,
				workflow.ChildWorkflowOptions{
					WorkflowID:            fmt.Sprintf("store-webhook-%s-%s", handleWebhooks.ConnectorID.String(), response.IdempotencyKey),
					TaskQueue:             handleWebhooks.ConnectorID.String(),
					ParentClosePolicy:     enums.PARENT_CLOSE_POLICY_ABANDON,
					WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
					SearchAttributes: map[string]interface{}{
						SearchAttributeStack: w.stack,
					},
				},
			),
			RunStoreWebhookTranslation,
			StoreWebhookTranslation{
				ConnectorID:     handleWebhooks.ConnectorID,
				Account:         response.Account,
				ExternalAccount: response.ExternalAccount,
				Payment:         response.Payment,
			},
		).Get(ctx, nil); err != nil {
			applicationError := &temporal.ApplicationError{}
			if errors.As(err, &applicationError) {
				if applicationError.Type() != "ChildWorkflowExecutionAlreadyStartedError" {
					return err
				}
			} else {
				return errors.Wrap(err, "running store workflow")
			}
		}
	}

	return nil
}

var RunHandleWebhooks any

type StoreWebhookTranslation struct {
	ConnectorID     models.ConnectorID
	Account         *models.PSPAccount
	ExternalAccount *models.PSPAccount
	Payment         *models.PSPPayment
}

func (w Workflow) runStoreWebhookTranslation(
	ctx workflow.Context,
	storeWebhookTranslation StoreWebhookTranslation,
) error {
	if storeWebhookTranslation.Account != nil {
		err := activities.StorageAccountsStore(
			infiniteRetryContext(ctx),
			models.FromPSPAccounts(
				[]models.PSPAccount{*storeWebhookTranslation.Account},
				models.ACCOUNT_TYPE_INTERNAL,
				storeWebhookTranslation.ConnectorID,
			),
		)
		if err != nil {
			return errors.Wrap(err, "storing next accounts")
		}
	}

	if storeWebhookTranslation.ExternalAccount != nil {
		err := activities.StorageAccountsStore(
			infiniteRetryContext(ctx),
			models.FromPSPAccounts(
				[]models.PSPAccount{*storeWebhookTranslation.ExternalAccount},
				models.ACCOUNT_TYPE_EXTERNAL,
				storeWebhookTranslation.ConnectorID,
			),
		)
		if err != nil {
			return errors.Wrap(err, "storing next accounts")
		}
	}

	if storeWebhookTranslation.Payment != nil {
		err := activities.StoragePaymentsStore(
			infiniteRetryContext(ctx),
			models.FromPSPPayments(
				[]models.PSPPayment{*storeWebhookTranslation.Payment},
				storeWebhookTranslation.ConnectorID,
			),
		)
		if err != nil {
			return errors.Wrap(err, "storing next accounts")
		}
	}

	return nil
}

var RunStoreWebhookTranslation any

func init() {
	RunHandleWebhooks = Workflow{}.runHandleWebhooks
	RunStoreWebhookTranslation = Workflow{}.runStoreWebhookTranslation
}
