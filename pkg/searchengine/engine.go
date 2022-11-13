package searchengine

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/formancehq/search/pkg/es"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/pkg/errors"
)

type Engine interface {
	doRequest(ctx context.Context, m map[string]interface{}) (*es.Response, error)
}

type EngineFn func(ctx context.Context, m map[string]interface{}) (*es.Response, error)

func (fn EngineFn) doRequest(ctx context.Context, m map[string]interface{}) (*es.Response, error) {
	return fn(ctx, m)
}

var NotImplementedEngine = EngineFn(func(ctx context.Context, m map[string]interface{}) (*es.Response, error) {
	return nil, errors.New("not implemented")
})

type DefaultEngineOption interface {
	apply(*DefaultEngine)
}
type DefaultEngineOptionFn func(engine *DefaultEngine)

func (fn DefaultEngineOptionFn) apply(engine *DefaultEngine) {
	fn(engine)
}

func WithESIndices(esIndices ...string) DefaultEngineOptionFn {
	return func(engine *DefaultEngine) {
		engine.indices = esIndices
	}
}

func WithRequestOption(opt func(req *opensearchapi.SearchRequest)) DefaultEngineOptionFn {
	return func(engine *DefaultEngine) {
		engine.requestOptions = append(engine.requestOptions, opt)
	}
}

var DefaultEsIndices = []string{"ledger"}

var DefaultEngineOptions = []DefaultEngineOption{
	WithESIndices(DefaultEsIndices...),
}

type Response map[string][]interface{}

type DefaultEngine struct {
	openSearchClient *opensearch.Client
	indices          []string
	requestOptions   []func(req *opensearchapi.SearchRequest)
}

func (e *DefaultEngine) doRequest(ctx context.Context, m map[string]interface{}) (*es.Response, error) {

	data, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling query")
	}

	httpResponse, err := e.openSearchClient.Search(
		e.openSearchClient.Search.WithBody(bytes.NewReader(data)),
		e.openSearchClient.Search.WithStoredFields("_all"),
		e.openSearchClient.Search.WithSource("*"),
		e.openSearchClient.Search.WithIndex(e.indices...),
		e.openSearchClient.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.IsError() {
		if httpResponse.StatusCode == 404 {
			return &es.Response{}, nil
		}
		return nil, errors.New(httpResponse.Status())
	}

	res := &es.Response{}
	err = json.NewDecoder(httpResponse.Body).Decode(res)
	if err != nil {
		return nil, errors.Wrap(err, "decoding result")
	}
	return res, nil
}

func NewDefaultEngine(openSearchClient *opensearch.Client, opts ...DefaultEngineOption) *DefaultEngine {

	engine := &DefaultEngine{
		openSearchClient: openSearchClient,
	}
	opts = append(DefaultEngineOptions, opts...)
	for _, opt := range opts {
		opt.apply(engine)
	}
	return engine
}
