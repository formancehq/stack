package testserver

import (
	"context"
	"encoding/json"
	"github.com/formancehq/stack/ee/ingester/internal"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
)

type HTTPConnector struct {
	srv       *httptest.Server
	collector *Collector
}

func (h *HTTPConnector) Clear(_ context.Context) error {
	h.collector.messages = nil
	return nil
}

func (h *HTTPConnector) Config() map[string]any {
	return map[string]any{
		"url": h.srv.URL,
	}
}

func (h *HTTPConnector) Name() string {
	return "http"
}

func (h *HTTPConnector) ReadMessages(_ context.Context) ([]ingester.LogWithModule, error) {
	h.collector.mu.Lock()
	defer h.collector.mu.Unlock()

	return h.collector.messages[:], nil
}

var _ Connector = &HTTPConnector{}

func NewHTTPConnector(t T, collector *Collector) Connector {
	ret := &HTTPConnector{
		collector: collector,
	}

	ret.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newMessages := make([]ingester.LogWithModule, 0)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&newMessages))

		ret.collector.mu.Lock()
		defer ret.collector.mu.Unlock()

		for _, message := range newMessages {
			exists := false
			for _, existingMessage := range ret.collector.messages {
				if existingMessage.Module == message.Module &&
					existingMessage.Shard == message.Shard &&
					existingMessage.ID == message.ID {
					exists = true
					break
				}
			}
			if !exists {
				ret.collector.messages = append(ret.collector.messages, message)
			}
		}

	}))
	t.Cleanup(ret.srv.Close)

	return ret
}

type Collector struct {
	mu       sync.Mutex
	messages []ingester.LogWithModule
}

func NewCollector() *Collector {
	return &Collector{}
}
