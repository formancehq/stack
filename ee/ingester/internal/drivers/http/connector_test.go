package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/drivers"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHTTPConnector(t *testing.T) {
	t.Parallel()

	messages := make(chan []ingester.LogWithModule, 1)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newMessages := make([]ingester.LogWithModule, 0)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&newMessages))

		messages <- newMessages
	}))
	t.Cleanup(testServer.Close)

	// Create our connector
	connector, err := NewConnector(drivers.NewServiceConfig(uuid.NewString(), testing.Verbose()), Config{
		URL: testServer.URL,
	}, logging.Testing())
	require.NoError(t, err)

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
				Payload: json.RawMessage(`{}`),
			},
		)
	}

	// Send all logs to the connector
	itemsErrors, err := connector.Accept(context.TODO(), logs...)
	require.NoError(t, err)
	require.Len(t, itemsErrors, numberOfLogs)
	for index := range logs {
		require.Nil(t, itemsErrors[index])
	}

	// Ensure data has been inserted
	select {
	case receivedMessages := <-messages:
		require.Len(t, receivedMessages, numberOfLogs)
	default:
		require.Fail(t, fmt.Sprintf("should have received %d messages", numberOfLogs))
	}
}
