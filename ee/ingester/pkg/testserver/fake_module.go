package testserver

import (
	"encoding/json"
	"github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
)

type Publisher interface {
	Publish(t require.TestingT, stack string, module string, data []byte)
}
type PublisherFn func(t require.TestingT, stack string, module string, data []byte)
func (fn PublisherFn) Publish(t require.TestingT, stack string, module string, data []byte) {
	fn(t, stack, module, data)
}

type FakeModule struct {
	mu        sync.Mutex
	stack     string
	name      string
	messages  []ingester.Log
	server    *httptest.Server
	publisher Publisher
}

func (m *FakeModule) PushLogs(t require.TestingT, logs ...ingester.Log) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, logs...)
	for _, log := range logs {
		data, err := json.Marshal(log)
		require.NoError(t, err)

		m.publisher.Publish(t, m.stack, m.name, data)
	}
}

func (m *FakeModule) URL() string {
	return m.server.URL
}

func NewFakeModule(t T, stack, name string, publisher Publisher) *FakeModule {
	ret := &FakeModule{
		stack:     stack,
		name:      name,
		publisher: publisher,
	}
	ret.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ret.mu.Lock()
		defer ret.mu.Unlock()

		api.RenderCursor(w, bunpaginate.Cursor[ingester.Log]{
			Data: ret.messages,
		})
	}))
	t.Cleanup(ret.server.Close)
	return ret
}
