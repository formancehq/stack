package webhooks

import (
	"errors"
	"sync"

	"github.com/formancehq/payments/internal/models"
)

type Webhooks interface {
	RegisterWebhooks(connectorID models.ConnectorID, webhooks []models.WebhookConfig)
	UnregisterWebhooks(connectorID models.ConnectorID)
	GetConfigs(connectorID models.ConnectorID, urlPath string) ([]models.WebhookConfig, error)
}

type webhooks struct {
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

func (w *webhooks) GetConfigs(connectorID models.ConnectorID, urlPath string) ([]models.WebhookConfig, error) {
	w.rwMutex.RLock()
	defer w.rwMutex.RUnlock()

	webhooksConfigs, ok := w.registeredWebhooksConfigs[connectorID.String()]
	if !ok {
		return nil, errors.New("connector not found")
	}

	return webhooksConfigs, nil
}

var _ Webhooks = &webhooks{}
