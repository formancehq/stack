package storage_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/stretchr/testify/require"
)

func TestConnectors(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testIsInstalledConnectors(t, store)
	testUpdateConfig(t, store)
	testEnableConnectors(t, store)
	testUninstallConnectors(t, store)
	testAfterInstallationConnectors(t, store)
}

func testInstallConnectors(t *testing.T, store *storage.Storage) {
	err := store.Install(
		context.Background(),
		models.ConnectorProviderDummyPay,
		json.RawMessage([]byte(`{"foo": "bar"}`)),
	)
	require.NoError(t, err)

	err = store.Install(
		context.Background(),
		models.ConnectorProviderDummyPay,
		json.RawMessage([]byte(`{"foo": "bar"}`)),
	)
	require.Equal(t, storage.ErrDuplicateKeyValue, err)

	testGetConnector(t, store, []byte(`{"foo": "bar"}`))

}

func testGetConnector(t *testing.T, store *storage.Storage, expectedConfig []byte) {
	var config json.RawMessage
	err := store.GetConfig(context.Background(), models.ConnectorProviderDummyPay, &config)
	require.NoError(t, err)
	require.Equal(t, json.RawMessage(expectedConfig), config)
}

func testUpdateConfig(t *testing.T, store *storage.Storage) {
	err := store.UpdateConfig(context.Background(), models.ConnectorProviderDummyPay, json.RawMessage([]byte(`{"foo2": "bar2"}`)))
	require.NoError(t, err)

	testGetConnector(t, store, []byte(`{"foo2": "bar2"}`))
}

func testEnableConnectors(t *testing.T, store *storage.Storage) {
	// Update dos not fail when connector is not installed
	err := store.Enable(context.Background(), models.ConnectorProviderMangopay)
	require.NoError(t, err)

	err = store.Enable(context.Background(), models.ConnectorProviderDummyPay)
	require.NoError(t, err)

	isEnabled, err := store.IsEnabled(context.Background(), models.ConnectorProviderMangopay)
	require.EqualError(t, err, storage.ErrNotFound.Error())
	require.False(t, isEnabled)

	isEnabled, err = store.IsEnabled(context.Background(), models.ConnectorProviderDummyPay)
	require.NoError(t, err)
	require.True(t, isEnabled)

	// Disable connectors
	err = store.Disable(context.Background(), models.ConnectorProviderMangopay)
	require.NoError(t, err)

	err = store.Disable(context.Background(), models.ConnectorProviderDummyPay)
	require.NoError(t, err)

	isEnabled, err = store.IsEnabled(context.Background(), models.ConnectorProviderDummyPay)
	require.NoError(t, err)
	require.False(t, isEnabled)
}

func testIsInstalledConnectors(t *testing.T, store *storage.Storage) {
	isInstalled, err := store.IsInstalled(context.Background(), models.ConnectorProviderMangopay)
	require.NoError(t, err)
	require.False(t, isInstalled)

	isInstalled, err = store.IsInstalled(context.Background(), models.ConnectorProviderDummyPay)
	require.NoError(t, err)
	require.True(t, isInstalled)
}

func testUninstallConnectors(t *testing.T, store *storage.Storage) {
	// No error if deleting an unknown connector
	err := store.Uninstall(context.Background(), models.ConnectorProviderMangopay)
	require.NoError(t, err)

	err = store.Uninstall(context.Background(), models.ConnectorProviderDummyPay)
	require.NoError(t, err)
}

func testAfterInstallationConnectors(t *testing.T, store *storage.Storage) {
	isInstalled, err := store.IsInstalled(context.Background(), models.ConnectorProviderDummyPay)
	require.NoError(t, err)
	require.False(t, isInstalled)
}
