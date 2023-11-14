package integration_test

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/cmd/connectors/internal/integration"
	"github.com/formancehq/payments/internal/models"
)

type connector struct {
	name     string
	id       models.ConnectorID
	config   json.RawMessage
	provider models.ConnectorProvider
}

type InMemoryConnectorStore struct {
	connectorsByID   map[string]*connector
	connectorsByName map[string]*connector
}

func (i *InMemoryConnectorStore) Uninstall(ctx context.Context, connectorID models.ConnectorID) error {
	connector, ok := i.connectorsByID[connectorID.String()]
	if !ok {
		return nil
	}

	delete(i.connectorsByID, connectorID.String())
	delete(i.connectorsByName, connector.name)

	return nil
}

func (i *InMemoryConnectorStore) ListConnectors(_ context.Context) ([]*models.Connector, error) {
	connectors := make([]*models.Connector, 0, len(i.connectorsByID))
	for _, c := range i.connectorsByID {
		connectors = append(connectors, &models.Connector{
			ID:       c.id,
			Name:     c.name,
			Config:   c.config,
			Provider: c.provider,
		})
	}
	return connectors, nil
}

func (i *InMemoryConnectorStore) IsInstalledByConnectorID(ctx context.Context, connectorID models.ConnectorID) (bool, error) {
	_, ok := i.connectorsByID[connectorID.String()]
	return ok, nil
}

func (i *InMemoryConnectorStore) IsInstalledByConnectorName(ctx context.Context, name string) (bool, error) {
	_, ok := i.connectorsByName[name]
	return ok, nil
}

func (i *InMemoryConnectorStore) Install(ctx context.Context, newConnector *models.Connector, config json.RawMessage) error {
	c := &connector{
		name:     newConnector.Name,
		id:       newConnector.ID,
		config:   config,
		provider: newConnector.Provider,
	}
	i.connectorsByID[newConnector.ID.String()] = c
	i.connectorsByName[newConnector.Name] = c

	return nil
}

func (i *InMemoryConnectorStore) UpdateConfig(ctx context.Context, connectorID models.ConnectorID, config json.RawMessage) error {
	i.connectorsByID[connectorID.String()].config = config
	return nil
}

func (i *InMemoryConnectorStore) GetConnector(ctx context.Context, connectorID models.ConnectorID) (*models.Connector, error) {
	c, ok := i.connectorsByID[connectorID.String()]
	if !ok {
		return nil, integration.ErrNotFound
	}

	return &models.Connector{
		ID:       c.id,
		Name:     c.name,
		Config:   c.config,
		Provider: c.provider,
	}, nil
}

func (i *InMemoryConnectorStore) ReadConfig(ctx context.Context, connectorID models.ConnectorID, to interface{}) error {
	connector, err := i.GetConnector(ctx, connectorID)
	if err != nil {
		return err
	}

	if err = connector.ParseConfig(to); err != nil {
		return err
	}

	return nil
}

var _ integration.Repository = &InMemoryConnectorStore{}

func NewInMemoryStore() *InMemoryConnectorStore {
	return &InMemoryConnectorStore{
		connectorsByID:   make(map[string]*connector),
		connectorsByName: make(map[string]*connector),
	}
}
