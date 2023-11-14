package driver_test

import (
	"testing"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/systemstore"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"

	"github.com/formancehq/ledger/internal/storage/storagetesting"
	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	t.Parallel()

	d := storagetesting.StorageDriver(t)
	ctx := logging.TestingContext()

	require.NoError(t, d.GetSystemStore().InsertConfiguration(ctx, "foo", "bar"))
	bar, err := d.GetSystemStore().GetConfiguration(ctx, "foo")
	require.NoError(t, err)
	require.Equal(t, "bar", bar)
}

func TestConfigurationError(t *testing.T) {
	t.Parallel()

	d := storagetesting.StorageDriver(t)
	ctx := logging.TestingContext()

	_, err := d.GetSystemStore().GetConfiguration(ctx, "not_existing")
	require.Error(t, err)
	require.True(t, sqlutils.IsNotFoundError(err))
}

func TestErrorOnOutdatedBucket(t *testing.T) {
	t.Parallel()

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

func TestGetLedgerFromDefaultBucket(t *testing.T) {
	t.Parallel()

	d := storagetesting.StorageDriver(t)
	ctx := logging.TestingContext()

	name := uuid.NewString()
	_, err := d.GetLedgerStore(ctx, name)
	require.NoError(t, err)
}

func TestGetLedgerFromAlternateBucket(t *testing.T) {
	t.Parallel()

	d := storagetesting.StorageDriver(t)
	ctx := logging.TestingContext()

	ledgerName := "ledger0"
	bucketName := "bucket0"
	_, err := d.GetSystemStore().RegisterLedger(ctx, &systemstore.Ledger{
		Name:    ledgerName,
		AddedAt: ledger.Now(),
		Bucket:  bucketName,
	})
	require.NoError(t, err)

	_, err = d.GetLedgerStore(ctx, ledgerName)
	require.NoError(t, err)
}
