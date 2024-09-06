package testserver

import (
	"context"
	"encoding/json"
	"github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/drivers/elasticsearch"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/olivere/elastic/v7"
	"sync"
)

type ElasticConnector struct {
	mu       sync.Mutex
	endpoint string
	client   *elastic.Client
}

func (h *ElasticConnector) Clear(ctx context.Context) error {
	_, err := h.client.Delete().Index(elasticsearch.DefaultIndex).Do(ctx)
	return err
}

func (h *ElasticConnector) ReadMessages(ctx context.Context) ([]ingester.LogWithModule, error) {

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client == nil {
		var err error
		h.client, err = elastic.NewClient(elastic.SetURL(h.endpoint))
		if err != nil {
			return nil, err
		}
	}

	response, err := h.client.
		Search(elasticsearch.DefaultIndex).
		Size(1000).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	return collectionutils.Map(response.Hits.Hits, func(from *elastic.SearchHit) ingester.LogWithModule {
		ret := ingester.LogWithModule{}
		if err := json.Unmarshal(from.Source, &ret); err != nil {
			panic(err)
		}

		return ret
	}), nil
}

func (h *ElasticConnector) Config() map[string]any {
	return map[string]any{
		"endpoint": h.endpoint,
	}
}

func (h *ElasticConnector) Name() string {
	return "elasticsearch"
}

var _ Connector = &ElasticConnector{}

func NewElasticConnector(endpoint string) Connector {
	return &ElasticConnector{
		endpoint: endpoint,
	}
}
