package integration

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/app/models"
)

type InMemoryConnectorStore struct {
	installed map[models.ConnectorProvider]bool
	disabled  map[models.ConnectorProvider]bool
	configs   map[models.ConnectorProvider]json.RawMessage
}

func (i *InMemoryConnectorStore) Uninstall(_ context.Context, name models.ConnectorProvider) error {
	delete(i.installed, name)
	delete(i.configs, name)
	delete(i.disabled, name)

	return nil
}

func (i *InMemoryConnectorStore) ListConnectors(_ context.Context) ([]*models.Connector, error) {
	return []*models.Connector{}, nil
}

func (i *InMemoryConnectorStore) IsInstalled(_ context.Context, name models.ConnectorProvider) (bool, error) {
	return i.installed[name], nil
}

func (i *InMemoryConnectorStore) Install(_ context.Context, name models.ConnectorProvider, config json.RawMessage) error {
	i.installed[name] = true
	i.configs[name] = config
	i.disabled[name] = false

	return nil
}

func (i *InMemoryConnectorStore) UpdateConfig(_ context.Context, name models.ConnectorProvider, config json.RawMessage) error {
	i.configs[name] = config

	return nil
}

func (i *InMemoryConnectorStore) Enable(_ context.Context, name models.ConnectorProvider) error {
	i.disabled[name] = false

	return nil
}

func (i *InMemoryConnectorStore) Disable(_ context.Context, name models.ConnectorProvider) error {
	i.disabled[name] = true

	return nil
}

func (i *InMemoryConnectorStore) IsEnabled(_ context.Context, name models.ConnectorProvider) (bool, error) {
	disabled, ok := i.disabled[name]
	if !ok {
		return false, nil
	}

	return !disabled, nil
}

func (i *InMemoryConnectorStore) GetConnector(_ context.Context, name models.ConnectorProvider) (*models.Connector, error) {
	cfg, ok := i.configs[name]
	if !ok {
		return nil, ErrNotFound
	}

	return &models.Connector{
		Config: cfg,
	}, nil
}

func (i *InMemoryConnectorStore) ReadConfig(ctx context.Context, name models.ConnectorProvider, to interface{}) error {
	connector, err := i.GetConnector(ctx, name)
	if err != nil {
		return err
	}

	return connector.ParseConfig(to)
}

func (i *InMemoryConnectorStore) CreateNewTransfer(_ context.Context, _ models.ConnectorProvider, _ string, _ string, _ string, _ int64) (models.Transfer, error) {
	return models.Transfer{}, nil
}

func (i *InMemoryConnectorStore) ListTransfers(_ context.Context, _ models.ConnectorProvider) ([]models.Transfer, error) {
	return []models.Transfer{}, nil
}

var _ Repository = &InMemoryConnectorStore{}

func NewInMemoryStore() *InMemoryConnectorStore {
	return &InMemoryConnectorStore{
		installed: make(map[models.ConnectorProvider]bool),
		disabled:  make(map[models.ConnectorProvider]bool),
		configs:   make(map[models.ConnectorProvider]json.RawMessage),
	}
}
