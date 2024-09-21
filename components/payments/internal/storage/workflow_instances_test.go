package storage

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/go-libs/query"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	defaultWorkflowInstances = []models.Instance{
		{
			ID:          "test1",
			ScheduleID:  defaultSchedules[0].ID,
			ConnectorID: defaultConnector.ID,
			CreatedAt:   now.Add(-60 * time.Minute).UTC().Time,
			UpdatedAt:   now.Add(-60 * time.Minute).UTC().Time,
			Terminated:  false,
		},
		{
			ID:          "test2",
			ScheduleID:  defaultSchedules[0].ID,
			ConnectorID: defaultConnector.ID,
			CreatedAt:   now.Add(-30 * time.Minute).UTC().Time,
			UpdatedAt:   now.Add(-30 * time.Minute).UTC().Time,
			Terminated:  false,
		},
		{
			ID:           "test3",
			ScheduleID:   defaultSchedules[2].ID,
			ConnectorID:  defaultConnector.ID,
			CreatedAt:    now.Add(-55 * time.Minute).UTC().Time,
			UpdatedAt:    now.Add(-55 * time.Minute).UTC().Time,
			Terminated:   true,
			TerminatedAt: pointer.For(now.UTC().Time),
			Error:        pointer.For("test error"),
		},
	}
)

func upsertInstance(t *testing.T, ctx context.Context, storage Storage, instance models.Instance) {
	require.NoError(t, storage.InstancesUpsert(ctx, instance))
}

func TestInstancesUpsert(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, schedule := range defaultSchedules {
		upsertSchedule(t, ctx, store, schedule)
	}
	for _, instance := range defaultWorkflowInstances {
		upsertInstance(t, ctx, store, instance)
	}

	t.Run("same id upsert", func(t *testing.T) {
		instance := defaultWorkflowInstances[0]
		instance.Terminated = true
		instance.TerminatedAt = pointer.For(now.UTC().Time)
		instance.Error = pointer.For("test error")

		upsertInstance(t, ctx, store, instance)

		actual, err := store.InstancesGet(ctx, instance.ID, instance.ScheduleID, instance.ConnectorID)
		require.NoError(t, err)
		require.Equal(t, defaultWorkflowInstances[0], *actual)
	})

	t.Run("unknown connector id", func(t *testing.T) {
		instance := defaultWorkflowInstances[0]
		instance.ConnectorID = models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}

		err := store.InstancesUpsert(ctx, instance)
		require.Error(t, err)
	})

	t.Run("unknown schedule id", func(t *testing.T) {
		instance := defaultWorkflowInstances[0]
		instance.ScheduleID = uuid.New().String()

		err := store.InstancesUpsert(ctx, instance)
		require.Error(t, err)
	})
}

func TestInstancesUpdate(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, schedule := range defaultSchedules {
		upsertSchedule(t, ctx, store, schedule)
	}
	for _, instance := range defaultWorkflowInstances {
		upsertInstance(t, ctx, store, instance)
	}

	t.Run("update instance error", func(t *testing.T) {
		instance := defaultWorkflowInstances[0]
		instance.Error = pointer.For("test error")
		instance.Terminated = true
		instance.TerminatedAt = pointer.For(now.UTC().Time)

		err := store.InstancesUpdate(ctx, instance)
		require.NoError(t, err)

		actual, err := store.InstancesGet(ctx, instance.ID, instance.ScheduleID, instance.ConnectorID)
		require.NoError(t, err)
		require.Equal(t, instance, *actual)
	})

	t.Run("update instance already on error", func(t *testing.T) {
		instance := defaultWorkflowInstances[2]
		instance.Error = pointer.For("test error2")
		instance.Terminated = true
		instance.TerminatedAt = pointer.For(now.UTC().Time)

		err := store.InstancesUpdate(ctx, instance)
		require.NoError(t, err)

		actual, err := store.InstancesGet(ctx, instance.ID, instance.ScheduleID, instance.ConnectorID)
		require.NoError(t, err)
		require.Equal(t, instance, *actual)
	})
}

func TestInstancesDeleteFromConnectorID(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, schedule := range defaultSchedules {
		upsertSchedule(t, ctx, store, schedule)
	}
	for _, instance := range defaultWorkflowInstances {
		upsertInstance(t, ctx, store, instance)
	}

	t.Run("delete instances from unknown connector", func(t *testing.T) {
		unknownConnectorID := models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}

		require.NoError(t, store.InstancesDeleteFromConnectorID(ctx, unknownConnectorID))

		for _, instance := range defaultWorkflowInstances {
			actual, err := store.InstancesGet(ctx, instance.ID, instance.ScheduleID, instance.ConnectorID)
			require.NoError(t, err)
			require.Equal(t, instance, *actual)
		}
	})

	t.Run("delete instances from default connector", func(t *testing.T) {
		require.NoError(t, store.InstancesDeleteFromConnectorID(ctx, defaultConnector.ID))

		for _, instance := range defaultWorkflowInstances {
			_, err := store.InstancesGet(ctx, instance.ID, instance.ScheduleID, instance.ConnectorID)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrNotFound)
		}
	})
}

func TestInstancesList(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	for _, schedule := range defaultSchedules {
		upsertSchedule(t, ctx, store, schedule)
	}
	for _, instance := range defaultWorkflowInstances {
		upsertInstance(t, ctx, store, instance)
	}

	t.Run("list instances by schedule_id", func(t *testing.T) {
		q := NewListInstancesQuery(
			bunpaginate.NewPaginatedQueryOptions(InstanceQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("schedule_id", defaultSchedules[0].ID)),
		)

		cursor, err := store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Equal(t, 2, len(cursor.Data))
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		require.Equal(t, defaultWorkflowInstances[1], cursor.Data[0])
		require.Equal(t, defaultWorkflowInstances[0], cursor.Data[1])
	})

	t.Run("list instances by unknown schedule_id", func(t *testing.T) {
		q := NewListInstancesQuery(
			bunpaginate.NewPaginatedQueryOptions(InstanceQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("schedule_id", uuid.New().String())),
		)

		cursor, err := store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Empty(t, cursor.Data)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
	})

	t.Run("list instances by connector_id", func(t *testing.T) {
		q := NewListInstancesQuery(
			bunpaginate.NewPaginatedQueryOptions(InstanceQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("connector_id", defaultConnector.ID)),
		)

		cursor, err := store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Equal(t, 3, len(cursor.Data))
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		require.Equal(t, defaultWorkflowInstances[1], cursor.Data[0])
		require.Equal(t, defaultWorkflowInstances[2], cursor.Data[1])
		require.Equal(t, defaultWorkflowInstances[0], cursor.Data[2])
	})

	t.Run("list instances by unknown connector_id", func(t *testing.T) {
		q := NewListInstancesQuery(
			bunpaginate.NewPaginatedQueryOptions(InstanceQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("connector_id", models.ConnectorID{
					Reference: uuid.New(),
					Provider:  "unknown",
				})),
		)

		cursor, err := store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Empty(t, cursor.Data)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
	})

	t.Run("list instances test cursor", func(t *testing.T) {
		q := NewListInstancesQuery(
			bunpaginate.NewPaginatedQueryOptions(InstanceQuery{}).
				WithPageSize(1),
		)

		cursor, err := store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, defaultWorkflowInstances[1], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Next, &q)
		require.NoError(t, err)
		cursor, err = store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.True(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, defaultWorkflowInstances[2], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Next, &q)
		require.NoError(t, err)
		cursor, err = store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.False(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		require.Equal(t, defaultWorkflowInstances[0], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &q)
		require.NoError(t, err)
		cursor, err = store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.True(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, defaultWorkflowInstances[2], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &q)
		require.NoError(t, err)
		cursor, err = store.InstancesList(ctx, q)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, defaultWorkflowInstances[1], cursor.Data[0])
	})
}
