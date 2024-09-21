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
	defaultWebhooks = []models.Webhook{
		{
			ID:          "test1",
			ConnectorID: defaultConnector.ID,
			QueryValues: map[string][]string{
				"foo": {"bar"},
			},
			Headers: map[string][]string{
				"foo2": {"bar2"},
			},
			Body: []byte(`{}`),
		},
		{
			ID:          "test2",
			ConnectorID: defaultConnector.ID,
			QueryValues: map[string][]string{
				"foo3": {"bar3"},
			},
			Headers: map[string][]string{
				"foo4": {"bar4"},
			},
			Body: []byte(`{}`),
		},
	}
)

func upsertWebhook(t *testing.T, ctx context.Context, storage Storage, webhook models.Webhook) {
	require.NoError(t, storage.WebhooksInsert(ctx, webhook))
}

func TestWebhooksInsert(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, webhook := range defaultWebhooks {
		upsertWebhook(t, ctx, store, webhook)
	}

	t.Run("same id upsert", func(t *testing.T) {
		webhook := defaultWebhooks[0]
		webhook.QueryValues = map[string][]string{
			"changed": {"changed"},
		}

		require.NoError(t, store.WebhooksInsert(ctx, webhook))

		// should not have been changed
		actual, err := store.WebhooksGet(ctx, webhook.ID)
		require.NoError(t, err)
		require.Equal(t, defaultWebhooks[0], actual)
	})

	t.Run("unknown connector id", func(t *testing.T) {
		webhook := defaultWebhooks[0]
		webhook.ID = "unknown"
		webhook.ConnectorID = models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}

		require.Error(t, store.WebhooksInsert(ctx, webhook))
	})
}

func TestWebhooksGet(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, webhook := range defaultWebhooks {
		upsertWebhook(t, ctx, store, webhook)
	}

	t.Run("get webhook", func(t *testing.T) {
		for _, webhook := range defaultWebhooks {
			actual, err := store.WebhooksGet(ctx, webhook.ID)
			require.NoError(t, err)
			require.Equal(t, webhook, actual)
		}
	})

	t.Run("get unknown webhook", func(t *testing.T) {
		_, err := store.WebhooksGet(ctx, "unknown")
		require.Error(t, err)
	})
}

func TestWebhooksDeleteFromConnectorID(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, webhook := range defaultWebhooks {
		upsertWebhook(t, ctx, store, webhook)
	}

	t.Run("delete unknown connector id", func(t *testing.T) {
		require.NoError(t, store.WebhooksDeleteFromConnectorID(ctx, models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}))

		for _, webhook := range defaultWebhooks {
			actual, err := store.WebhooksGet(ctx, webhook.ID)
			require.NoError(t, err)
			require.Equal(t, webhook, actual)
		}
	})

	t.Run("delete webhooks", func(t *testing.T) {
		require.NoError(t, store.WebhooksDeleteFromConnectorID(ctx, defaultConnector.ID))

		for _, webhook := range defaultWebhooks {
			_, err := store.WebhooksGet(ctx, webhook.ID)
			require.Error(t, err)
		}
	})
}
