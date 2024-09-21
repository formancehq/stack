package storage

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	defaultWebhooksConfigs = []models.WebhookConfig{
		{
			Name:        "test1",
			ConnectorID: defaultConnector.ID,
			URLPath:     "/test1",
		},
		{
			Name:        "test2",
			ConnectorID: defaultConnector.ID,
			URLPath:     "/test2",
		},
		{
			Name:        "test3",
			ConnectorID: defaultConnector.ID,
			URLPath:     "/test3",
		},
	}
)

func upsertWebhookConfigs(t *testing.T, ctx context.Context, storage Storage, webhookConfigs []models.WebhookConfig) {
	require.NoError(t, storage.WebhooksConfigsUpsert(ctx, webhookConfigs))
}

func TestWebhooksConfigsUpsert(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertWebhookConfigs(t, ctx, store, defaultWebhooksConfigs)

	t.Run("same name and connector id insert", func(t *testing.T) {
		w := models.WebhookConfig{
			Name:        "test1",
			ConnectorID: defaultConnector.ID,
			URLPath:     "/test3",
		}

		require.NoError(t, store.WebhooksConfigsUpsert(ctx, []models.WebhookConfig{w}))

		actual, err := store.WebhooksConfigsGet(ctx, w.Name, w.ConnectorID)
		require.NoError(t, err)
		require.Equal(t, defaultWebhooksConfigs[0], *actual)
	})

	t.Run("unknown connector id", func(t *testing.T) {
		w := models.WebhookConfig{
			Name: "test1",
			ConnectorID: models.ConnectorID{
				Reference: uuid.New(),
				Provider:  "unknown",
			},
			URLPath: "/test3",
		}

		require.Error(t, store.WebhooksConfigsUpsert(ctx, []models.WebhookConfig{w}))
	})
}

func TestWebhooksConfigsGet(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertWebhookConfigs(t, ctx, store, defaultWebhooksConfigs)

	t.Run("get webhook config", func(t *testing.T) {
		for _, w := range defaultWebhooksConfigs {
			actual, err := store.WebhooksConfigsGet(ctx, w.Name, w.ConnectorID)
			require.NoError(t, err)
			require.Equal(t, w, *actual)
		}
	})

	t.Run("unknown webhook config", func(t *testing.T) {
		_, err := store.WebhooksConfigsGet(ctx, "unknown", defaultConnector.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)
	})
}

func TestWebhooksConfigsDeleteFromConnectorID(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertWebhookConfigs(t, ctx, store, defaultWebhooksConfigs)

	t.Run("delete webhooks configs from unknown connector id", func(t *testing.T) {
		require.NoError(t, store.WebhooksConfigsDeleteFromConnectorID(ctx, models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}))

		for _, w := range defaultWebhooksConfigs {
			actual, err := store.WebhooksConfigsGet(ctx, w.Name, w.ConnectorID)
			require.NoError(t, err)
			require.Equal(t, w, *actual)
		}
	})

	t.Run("delete webhooks configs from connector id", func(t *testing.T) {
		require.NoError(t, store.WebhooksConfigsDeleteFromConnectorID(ctx, defaultConnector.ID))

		for _, w := range defaultWebhooksConfigs {
			_, err := store.WebhooksConfigsGet(ctx, w.Name, w.ConnectorID)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrNotFound)
		}
	})
}
