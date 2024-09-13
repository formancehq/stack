package webhooks

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/formancehq/payments/internal/connectors/engine/workflow"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type Webhooks interface {
	RegisterWebhooks(connectorID models.ConnectorID, webhooks []models.WebhookConfig)
	UnregisterWebhooks(connectorID models.ConnectorID)
	HandleWebhook(ctx context.Context, urlPath string, webhook models.Webhook) error
}

type webhooks struct {
	temporalClient client.Client

	registeredWebhooksConfigs map[string][]models.WebhookConfig
	rwMutex                   sync.RWMutex
}

func New() *webhooks {
	return &webhooks{
		registeredWebhooksConfigs: make(map[string][]models.WebhookConfig),
	}
}

func (w *webhooks) RegisterWebhooks(connectorID models.ConnectorID, webhooks []models.WebhookConfig) {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	w.registeredWebhooksConfigs[connectorID.String()] = webhooks
}

func (w *webhooks) UnregisterWebhooks(connectorID models.ConnectorID) {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	delete(w.registeredWebhooksConfigs, connectorID.String())
}

func (w *webhooks) HandleWebhook(
	ctx context.Context,
	urlPath string,
	webhook models.Webhook,
) error {
	w.rwMutex.RLock()
	defer w.rwMutex.RUnlock()

	webhooksConfigs, ok := w.registeredWebhooksConfigs[webhook.ConnectorID.String()]
	if !ok {
		return errors.New("connector not found")
	}

	for _, config := range webhooksConfigs {
		if !strings.Contains(urlPath, config.URLPath) {
			continue
		}

		ctx, _ := contextutil.Detached(ctx)
		if _, err := w.temporalClient.ExecuteWorkflow(
			ctx,
			client.StartWorkflowOptions{
				ID:                                       fmt.Sprintf("webhook-%s-%s", webhook.ConnectorID.String(), webhook.ID),
				TaskQueue:                                webhook.ConnectorID.Reference,
				WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
				WorkflowExecutionErrorWhenAlreadyStarted: false,
			},
			workflow.RunHandleWebhooks,
			workflow.HandleWebhooks{
				ConnectorID:   webhook.ConnectorID,
				WebhookConfig: config,
				Webhook:       webhook,
			},
		); err != nil {
			return err
		}

		// Nothing to do after that
		return nil
	}

	return nil
}

var _ Webhooks = &webhooks{}
