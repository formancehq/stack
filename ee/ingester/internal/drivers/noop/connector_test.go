package noop

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNoOpConnector(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// Create our connector
	connector, err := NewConnector(drivers.NewServiceConfig(uuid.NewString(), testing.Verbose()), struct{}{}, logging.Testing())
	require.NoError(t, err)
	require.NoError(t, connector.Start(ctx))
	t.Cleanup(func() {
		require.NoError(t, connector.Stop(ctx))
	})

	// We will insert numberOfLogs logs split across numberOfModules modules
	const (
		numberOfLogs    = 50
		numberOfModules = 2
	)
	logs := make([]ingester.LogWithModule, numberOfLogs)
	for i := 0; i < numberOfLogs; i++ {
		logs[i] = ingester.NewLogWithModule(
			fmt.Sprintf("module%d", i%numberOfModules),
			ingester.Log{
				Shard:   "test",
				ID:      fmt.Sprint(i),
				Date:    time.Now(),
				Type:    "test",
				Payload: json.RawMessage(``),
			},
		)
	}

	// Send all logs to the connector
	itemsErrors, err := connector.Accept(ctx, logs...)
	require.NoError(t, err)
	require.Len(t, itemsErrors, numberOfLogs)
	for index := range logs {
		require.Nil(t, itemsErrors[index])
	}
}
