package integration

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/formancehq/payments/internal/app/models"
)

type InMemoryConnectorStore struct {
	installed map[models.ConnectorProvider]bool
	disabled  map[models.ConnectorProvider]bool
	configs   map[models.ConnectorProvider]json.RawMessage
}

func (i *InMemoryConnectorStore) Uninstall(ctx context.Context, name models.ConnectorProvider) error {
	delete(i.installed, name)
	delete(i.configs, name)
	delete(i.disabled, name)

	return nil
}

func (i *InMemoryConnectorStore) ListConnectors(_ context.Context) ([]*models.Connector, error) {
	return []*models.Connector{}, nil
}

func (i *InMemoryConnectorStore) IsInstalled(ctx context.Context, name models.ConnectorProvider) (bool, error) {
	return i.installed[name], nil
}

func (i *InMemoryConnectorStore) Install(ctx context.Context, name models.ConnectorProvider, config json.RawMessage) error {
	i.installed[name] = true
	i.configs[name] = config
	i.disabled[name] = false

	return nil
}

func (i *InMemoryConnectorStore) UpdateConfig(ctx context.Context, name models.ConnectorProvider, config json.RawMessage) error {
	i.configs[name] = config

	return nil
}

func (i *InMemoryConnectorStore) Enable(ctx context.Context, name models.ConnectorProvider) error {
	i.disabled[name] = false

	return nil
}

func (i *InMemoryConnectorStore) Disable(ctx context.Context, name models.ConnectorProvider) error {
	i.disabled[name] = true

	return nil
}

func (i *InMemoryConnectorStore) IsEnabled(ctx context.Context, name models.ConnectorProvider) (bool, error) {
	disabled, ok := i.disabled[name]
	if !ok {
		return false, nil
	}

	return !disabled, nil
}

func (i *InMemoryConnectorStore) GetConnector(ctx context.Context, name models.ConnectorProvider) (*models.Connector, error) {
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

	if err = connector.ParseConfig(to); err != nil {
		return err
	}

	return nil
}

func (i *InMemoryConnectorStore) CreateNewTransfer(ctx context.Context, name models.ConnectorProvider, source, destination, currency string, amount *big.Int) (models.Transfer, error) {
	return models.Transfer{}, nil
}

func (i *InMemoryConnectorStore) ListTransfers(ctx context.Context, name models.ConnectorProvider) ([]models.Transfer, error) {
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
