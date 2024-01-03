package storage

import (
	"context"
	"testing"
	"time"

	//lint:ignore ST1001 shared definitions
	. "github.com/formancehq/webhooks/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigStore(t *testing.T, newConfig func(t *testing.T) ConfigStore) {
	t.Parallel()
	run := func(name string, testFunc func(*testing.T, func(*testing.T) ConfigStore)) {
		t.Run(name, func(t *testing.T) {
			testFunc(t, newConfig)
		})
	}
	run("FindManyConfigs", testFindManyConfigs)
	run("DeleteOneConfig", testDeleteOneConfig)
	run("UpdateOneConfigActivation", testUpdateOneConfigActivation)
	run("UpdateOneConfigSecret", testUpdateOneConfigSecret)
}

func testFindManyConfigs(t *testing.T, newStore func(*testing.T) ConfigStore) {
	store := newStore(t)
	ctx, cancel := newTestContext()
	defer cancel()

	basicConfig, err := store.InsertOneConfig(ctx, ConfigUser{
		Endpoint:   "http://example.com/webhook",
		EventTypes: []string{"someType"},
	})
	assert.NoErrorf(t, err, "couldn't insert fixture")

	idleConfig, err := store.InsertOneConfig(ctx, ConfigUser{
		Endpoint:   "http://example.com/webhook",
		EventTypes: []string{"otherType"},
	})
	assert.NoErrorf(t, err, "couldn't insert fixture")
	idleConfig, err = store.UpdateOneConfigActivation(ctx, idleConfig.ID, false)
	assert.NoErrorf(t, err, "couldn't deactivate fixture")

	otherConfig, err := store.InsertOneConfig(ctx, ConfigUser{
		Endpoint:   "http://example.com/webhook2",
		EventTypes: []string{"someType", "otherType"},
	})
	assert.NoErrorf(t, err, "couldn't insert fixture")

	t.Parallel()
	t.Run("nil filter", func(t *testing.T) {
		result, err := store.FindManyConfigs(ctx, nil)

		assert.NoError(t, err)
		assertConfigSliceEqual(t,
			[]Config{otherConfig, idleConfig, basicConfig},
			result,
		)
	})
	t.Run("empty filter", func(t *testing.T) {
		result, err := store.FindManyConfigs(ctx, map[string]any{})

		assert.NoError(t, err)
		assertConfigSliceEqual(t,
			[]Config{otherConfig, idleConfig, basicConfig},
			result,
		)
	})
	t.Run("filter by id", func(t *testing.T) {
		result, err := store.FindManyConfigs(ctx, map[string]any{
			"id": basicConfig.ID,
		})

		assert.NoError(t, err)
		assertConfigSliceEqual(t,
			[]Config{basicConfig},
			result,
		)
	})
	t.Run("filter by endpoint", func(t *testing.T) {
		result, err := store.FindManyConfigs(ctx, map[string]any{
			"endpoint": "http://example.com/webhook",
		})

		assert.NoError(t, err)
		assertConfigSliceEqual(t,
			[]Config{idleConfig, basicConfig},
			result,
		)
	})
	t.Run("filter active configs", func(t *testing.T) {
		result, err := store.FindManyConfigs(ctx, map[string]any{
			"active": true,
		})

		assert.NoError(t, err)
		assertConfigSliceEqual(t,
			[]Config{otherConfig, basicConfig},
			result,
		)
	})
	t.Run("filter by event type", func(t *testing.T) {
		result, err := store.FindManyConfigs(ctx, map[string]any{
			"event_types": "someType",
		})

		assert.NoError(t, err)
		assertConfigSliceEqual(t,
			[]Config{otherConfig, basicConfig},
			result,
		)
	})
	t.Run("cumulated fields", func(t *testing.T) {
		result, err := store.FindManyConfigs(ctx, map[string]any{
			"event_types": "someType",
			"active":      true,
			"endpoint":    "http://example.com/webhook2",
		})

		assert.NoError(t, err)
		assertConfigSliceEqual(t,
			[]Config{otherConfig},
			result,
		)
	})
}

func assertConfigSliceEqual(t *testing.T, want, got []Config) {
	t.Helper()
	require.Equal(t, len(want), len(got), "slices should have the same length")
	for i := range want {
		assertConfigEqual(t, want[i], got[i])
	}
}

func assertConfigEqual(t *testing.T, want, got Config) {
	t.Helper()
	require.Equal(t, want.ConfigUser, got.ConfigUser, "Config.ConfigUser")
	require.Equal(t, want.Active, got.Active, "Config.Active")
}

func testDeleteOneConfig(t *testing.T, newStore func(t *testing.T) ConfigStore) {
	store := newStore(t)
	ctx, cancel := newTestContext()
	defer cancel()

	t.Parallel()
	t.Run("config does not exist", func(t *testing.T) {
		err := store.DeleteOneConfig(ctx, uuid.NewString())
		assert.ErrorIs(t, err, ErrConfigNotFound)
	})
	t.Run("nominal", func(t *testing.T) {
		basicConfig, err := store.InsertOneConfig(ctx, ConfigUser{
			Endpoint:   "http://example.com/webhook",
			EventTypes: []string{"someType"},
		})
		assert.NoErrorf(t, err, "couldn't insert fixture")

		err = store.DeleteOneConfig(ctx, basicConfig.ID)
		assert.NoError(t, err)

		result, err := store.FindManyConfigs(ctx, map[string]any{"id": basicConfig.ID})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func testUpdateOneConfigActivation(t *testing.T, newStore func(*testing.T) ConfigStore) {
	store := newStore(t)
	ctx, cancel := newTestContext()
	defer cancel()

	t.Run("config does not exist", func(t *testing.T) {
		_, err := store.UpdateOneConfigActivation(ctx, uuid.NewString(), false)
		assert.ErrorIs(t, err, ErrConfigNotFound)
	})

	t.Run("config not modified", func(t *testing.T) {
		cfg, err := store.InsertOneConfig(ctx, ConfigUser{
			Endpoint:   "http://example.com/webhook",
			Secret:     "secret",
			EventTypes: []string{"someType"},
		})
		assert.NoError(t, err)
		_, err = store.UpdateOneConfigActivation(ctx, cfg.ID, true)
		assert.ErrorIs(t, err, ErrConfigNotModified)
	})

	t.Run("nominal", func(t *testing.T) {
		cfg, err := store.InsertOneConfig(ctx, ConfigUser{
			Endpoint:   "http://example.com/webhook",
			Secret:     "secret",
			EventTypes: []string{"someType"},
		})
		assert.NoError(t, err)
		result, err := store.UpdateOneConfigActivation(ctx, cfg.ID, false)
		assert.NoError(t, err)
		assert.Equal(t, cfg.ID, result.ID)
		assert.False(t, result.Active)
	})
}

func testUpdateOneConfigSecret(t *testing.T, newStore func(*testing.T) ConfigStore) {
	store := newStore(t)
	ctx, cancel := newTestContext()
	defer cancel()

	t.Run("config does not exist", func(t *testing.T) {
		_, err := store.UpdateOneConfigSecret(ctx, uuid.NewString(), "secret")
		assert.ErrorIs(t, err, ErrConfigNotFound)
	})

	t.Run("config not modified", func(t *testing.T) {
		cfg, err := store.InsertOneConfig(ctx, ConfigUser{
			Endpoint:   "http://example.com/webhook",
			Secret:     "secret",
			EventTypes: []string{"someType"},
		})
		assert.NoError(t, err)
		_, err = store.UpdateOneConfigSecret(ctx, cfg.ID, "secret")
		assert.ErrorIs(t, err, ErrConfigNotModified)
	})

	t.Run("nominal", func(t *testing.T) {
		cfg, err := store.InsertOneConfig(ctx, ConfigUser{
			Endpoint:   "http://example.com/webhook",
			Secret:     "secret",
			EventTypes: []string{"someType"},
		})
		assert.NoError(t, err)
		result, err := store.UpdateOneConfigSecret(ctx, cfg.ID, "new secret")
		assert.NoError(t, err)
		assert.Equal(t, cfg.ID, result.ID)
		assert.Equal(t, "new secret", result.Secret)
	})
}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}
