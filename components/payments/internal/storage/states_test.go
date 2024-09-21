package storage

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	defaultStates = []models.State{
		{
			ID: models.StateID{
				Reference:   "test1",
				ConnectorID: defaultConnector.ID,
			},
			ConnectorID: defaultConnector.ID,
			State:       []byte(`{}`),
		},
		{
			ID: models.StateID{
				Reference:   "test2",
				ConnectorID: defaultConnector.ID,
			},
			ConnectorID: defaultConnector.ID,
			State:       []byte(`{"foo":"bar"}`),
		},
		{
			ID: models.StateID{
				Reference:   "test3",
				ConnectorID: defaultConnector.ID,
			},
			ConnectorID: defaultConnector.ID,
			State:       []byte(`{"foo3":"bar3"}`),
		},
	}
)

func upsertState(t *testing.T, ctx context.Context, storage Storage, state models.State) {
	require.NoError(t, storage.StatesUpsert(ctx, state))
}

func TestStatesUpsert(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, state := range defaultStates {
		upsertState(t, ctx, store, state)
	}

	t.Run("upsert with unknown connector id", func(t *testing.T) {
		c := models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}
		s := models.State{
			ID: models.StateID{
				Reference:   "test4",
				ConnectorID: c,
			},
			ConnectorID: c,
			State:       []byte(`{}`),
		}

		require.Error(t, store.StatesUpsert(ctx, s))
	})

	t.Run("upsert with same id", func(t *testing.T) {
		s := models.State{
			ID:          defaultStates[0].ID,
			ConnectorID: defaultConnector.ID,
			State:       json.RawMessage(`{"foo":"bar"}`),
		}

		upsertState(t, ctx, store, s)

		// Should update the state
		state, err := store.StatesGet(ctx, s.ID)
		require.NoError(t, err)
		require.Equal(t, s, state)
	})
}

func TestStatesGet(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, state := range defaultStates {
		upsertState(t, ctx, store, state)
	}

	t.Run("get state", func(t *testing.T) {
		for _, state := range defaultStates {
			s, err := store.StatesGet(ctx, state.ID)
			require.NoError(t, err)
			require.Equal(t, state, s)
		}
	})

	t.Run("get state with unknown id", func(t *testing.T) {
		_, err := store.StatesGet(ctx, models.StateID{
			Reference:   "unknown",
			ConnectorID: defaultConnector.ID,
		})
		require.Error(t, err)
	})
}

func TestDeleteStatesFromConnectorID(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, state := range defaultStates {
		upsertState(t, ctx, store, state)
	}

	t.Run("delete states with unknown connector id", func(t *testing.T) {
		require.NoError(t, store.StatesDeleteFromConnectorID(ctx, models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}))

		for _, state := range defaultStates {
			s, err := store.StatesGet(ctx, state.ID)
			require.NoError(t, err)
			require.Equal(t, state, s)
		}
	})

	t.Run("delete states", func(t *testing.T) {
		require.NoError(t, store.StatesDeleteFromConnectorID(ctx, defaultConnector.ID))

		for _, state := range defaultStates {
			_, err := store.StatesGet(ctx, state.ID)
			require.Error(t, err)
		}
	})
}
