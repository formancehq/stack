package modules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/httpclient"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/ThreeDotsLabs/watermill/message"
)

//go:generate mockgen -source module.go -destination module_generated.go -package modules . Module
type Module interface {
	Subscribe(ctx context.Context) (<-chan *message.Message, error)
	Pull(ctx context.Context, cursor string) (*bunpaginate.Cursor[ingester.Log], error)
}

// stackModule implement Module using default tooling of the stack
// It uses the shared library to connect to the broker and a plain old http client to query modules for logs
type stackModule struct {
	subscriber        message.Subscriber
	module            string
	stack             string
	httpClient        *httpclient.StackAuthenticatedClient
	logger            logging.Logger
	pullConfiguration PullConfiguration
}

func (s *stackModule) Subscribe(ctx context.Context) (<-chan *message.Message, error) {
	topic := fmt.Sprintf("%s.%s", s.stack, s.module)
	s.logger.Infof("Subscribing to %s", topic)

	return s.subscriber.Subscribe(ctx, topic)
}

func (s *stackModule) Pull(ctx context.Context, cursor string) (*bunpaginate.Cursor[ingester.Log], error) {

	buf := bytes.NewBuffer([]byte{})
	if err := s.pullConfiguration.ModuleURLTpl.Execute(buf, map[string]any{
		"module": s.module,
	}); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, buf.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.URL.RawQuery = url.Values{
		"cursor":                     []string{cursor},
		bunpaginate.QueryKeyPageSize: []string{fmt.Sprintf("%d", s.pullConfiguration.PullPageSize)},
	}.Encode()

	rsp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	ret := api.BaseResponse[ingester.Log]{}
	if err := json.NewDecoder(rsp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return ret.Cursor, nil
}

var _ Module = (*stackModule)(nil)

//go:generate mockgen -source module.go -destination module_generated.go -package modules . Factory
type Factory interface {
	Create(name string) Module
}

type stackModuleFactory struct {
	subscriber        message.Subscriber
	stack             string
	httpClient        *httpclient.StackAuthenticatedClient
	logger            logging.Logger
	pullConfiguration PullConfiguration
}

func (s *stackModuleFactory) Create(name string) Module {
	return &stackModule{
		subscriber:        s.subscriber,
		module:            name,
		stack:             s.stack,
		httpClient:        s.httpClient,
		pullConfiguration: s.pullConfiguration,
		logger: s.logger.WithFields(map[string]any{
			"component": "module",
			"module":    name,
		}),
	}
}

func NewModuleFactory(
	subscriber message.Subscriber,
	stack string,
	httpClient *httpclient.StackAuthenticatedClient,
	pullConfiguration PullConfiguration,
	logger logging.Logger,
) Factory {
	return &stackModuleFactory{
		subscriber:        subscriber,
		stack:             stack,
		httpClient:        httpClient,
		logger:            logger,
		pullConfiguration: pullConfiguration,
	}
}

var _ Factory = (*stackModuleFactory)(nil)
