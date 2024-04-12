package bunconnect_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testContext struct {
	context.Context
	doneCalled chan struct{}
}

func (t *testContext) Done() <-chan struct{} {
	select {
	case <-t.doneCalled:
	default:
		close(t.doneCalled)
	}
	return t.Context.Done()
}

func TestConnect(t *testing.T) {
	require.NoError(t, pgtesting.CreatePostgresServer())
	t.Cleanup(func() {
		require.NoError(t, pgtesting.DestroyPostgresServer())
	})

	ctx := logging.TestingContext()
	dbName := uuid.NewString()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	ctx = &testContext{ctx, make(chan struct{})}

	done := make(chan struct{})
	go func() {
		db, err := bunconnect.OpenSQLDB(ctx, bunconnect.ConnectionOptions{
			DatabaseSourceName: fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
				pgtesting.Server().GetUsername(),
				pgtesting.Server().GetPassword(),
				pgtesting.Server().GetHost(),
				pgtesting.Server().GetPort(),
				dbName,
			),
		})
		require.NoError(t, err)
		require.NoError(t, db.Close())
		close(done)
	}()

	select {
	case <-ctx.(*testContext).doneCalled:
	case <-time.After(2 * time.Second):
		require.Fail(t, "Done() should have been called on context")
	}

	_ = pgtesting.NewNamedPostgresDatabase(t, dbName)

	select {
	case <-done:
	case <-time.After(bunconnect.PingInterval * 2):
		require.Fail(t, "unable to connect to the database")
	}
}
