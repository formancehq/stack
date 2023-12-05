package storage_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const encryptionOptions = "compress-algo=1, cipher-algo=aes256"
const encryptionKey = "test"

// Helpers to add test data
func installConnector(t *testing.T, store *storage.Storage) models.ConnectorID {
	db := store.DB()

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}
	connector := &models.Connector{
		ID:        connectorID,
		Name:      "test_connector",
		CreatedAt: time.Date(2023, 11, 13, 0, 0, 0, 0, time.UTC),
		Provider:  models.ConnectorProviderDummyPay,
	}

	_, err := db.NewInsert().Model(connector).Exec(context.Background())
	require.NoError(t, err)

	_, err = db.NewUpdate().
		Model(&models.Connector{}).
		Set("config = pgp_sym_encrypt(?::TEXT, ?, ?)", json.RawMessage(`{}`), encryptionKey, encryptionOptions).
		Where("id = ?", connectorID). // Connector name is unique
		Exec(context.Background())
	require.NoError(t, err)

	return connectorID
}
