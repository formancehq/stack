package storage_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	connectorID = models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}
)

func TestConnectors(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testIsInstalledConnectors(t, store)
	testUpdateConfig(t, store)
	testUninstallConnectors(t, store)
	testAfterInstallationConnectors(t, store)
}

func testInstallConnectors(t *testing.T, store *storage.Storage) {
	connector1 := models.Connector{
		ID:       connectorID,
		Name:     "test1",
		Provider: models.ConnectorProviderDummyPay,
	}
	err := store.Install(
		context.Background(),
		&connector1,
		json.RawMessage([]byte(`{"foo": "bar"}`)),
	)
	require.NoError(t, err)

	err = store.Install(
		context.Background(),
		&connector1,
		json.RawMessage([]byte(`{"foo": "bar"}`)),
	)
	require.Equal(t, storage.ErrDuplicateKeyValue, err)

	testGetConnector(t, store, connectorID, []byte(`{"foo": "bar"}`))
}

func testGetConnector(t *testing.T, store *storage.Storage, connectorID models.ConnectorID, expectedConfig []byte) {
	var config json.RawMessage
	err := store.GetConfig(context.Background(), connectorID, &config)
	require.NoError(t, err)
	require.Equal(t, json.RawMessage(expectedConfig), config)
}

func testUpdateConfig(t *testing.T, store *storage.Storage) {
	err := store.UpdateConfig(context.Background(), connectorID, json.RawMessage([]byte(`{"foo2": "bar2"}`)))
	require.NoError(t, err)

	testGetConnector(t, store, connectorID, []byte(`{"foo2": "bar2"}`))
}

func testIsInstalledConnectors(t *testing.T, store *storage.Storage) {
	isInstalled, err := store.IsInstalledByConnectorID(
		context.Background(),
		models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		})
	require.NoError(t, err)
	require.False(t, isInstalled)

	isInstalled, err = store.IsInstalledByConnectorID(context.Background(), connectorID)
	require.NoError(t, err)
	require.True(t, isInstalled)
}

func testUninstallConnectors(t *testing.T, store *storage.Storage) {
	// No error if deleting an unknown connector
	err := store.Uninstall(context.Background(), models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	})
	require.NoError(t, err)

	err = store.Uninstall(context.Background(), connectorID)
	require.NoError(t, err)
}

func testAfterInstallationConnectors(t *testing.T, store *storage.Storage) {
	isInstalled, err := store.IsInstalledByConnectorID(context.Background(), connectorID)
	require.NoError(t, err)
	require.False(t, isInstalled)
}
