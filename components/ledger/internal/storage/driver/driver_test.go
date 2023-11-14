package driver_test

import (
	"testing"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"

	"github.com/formancehq/ledger/internal/storage/storagetesting"
	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	d := storagetesting.StorageDriver(t)
	ctx := logging.TestingContext()

	require.NoError(t, d.GetSystemStore().InsertConfiguration(ctx, "foo", "bar"))
	bar, err := d.GetSystemStore().GetConfiguration(ctx, "foo")
	require.NoError(t, err)
	require.Equal(t, "bar", bar)
}

func TestConfigurationError(t *testing.T) {
	d := storagetesting.StorageDriver(t)
	ctx := logging.TestingContext()

	_, err := d.GetSystemStore().GetConfiguration(ctx, "not_existing")
	require.Error(t, err)
	require.True(t, sqlutils.IsNotFoundError(err))
}

func TestErrorOnOutdatedBucket(t *testing.T) {
	ctx := logging.TestingContext()
	d := storagetesting.StorageDriver(t)

	name := uuid.NewString()

	_, err := d.GetSystemStore().RegisterBucket(ctx, name)
	require.NoError(t, err)

	b, err := d.GetBucket(ctx, name)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = b.Close()
	})

	upToDate, err := b.IsUpToDate(ctx)
	require.NoError(t, err)
	require.False(t, upToDate)
}
