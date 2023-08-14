package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alitto/pond"
	"github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/client/metrics"
	"github.com/formancehq/stack/components/stargate/internal/opentelemetry"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

type WorkerPoolConfig struct {
	MaxWorkers int
	MaxTasks   int
}

func NewWorkerPoolConfig(maxWorkers, maxTasks int) WorkerPoolConfig {
	return WorkerPoolConfig{
		MaxWorkers: maxWorkers,
		MaxTasks:   maxTasks,
	}
}

type Config struct {
	OrganizationID          string
	StackID                 string
	ChanSize                int
	GatewayUrl              string
	HTTPClientTimeout       time.Duration
	HTTPMaxIdleConns        int
	HTTPMaxIdleConnsPerHost int
}

func NewClientConfig(
	organizationID string,
	stackID string,
	chanSize int,
	gatewayUrl string,
	httpClientTimeout time.Duration,
	httpMaxIdleConns int,
	httpMaxIdleConnsPerHost int,
) Config {
	return Config{
		OrganizationID:          organizationID,
		StackID:                 stackID,
		ChanSize:                chanSize,
		GatewayUrl:              gatewayUrl,
		HTTPClientTimeout:       httpClientTimeout,
		HTTPMaxIdleConns:        httpMaxIdleConns,
		HTTPMaxIdleConnsPerHost: httpMaxIdleConnsPerHost,
	}
}

type Client struct {
	logger         logging.Logger
	config         Config
	stargateClient api.StargateServiceClient
	httpClient     *http.Client

	workerPool      *pond.WorkerPool
	metricsRegistry metrics.MetricsRegistry
}

func NewClient(
	l logging.Logger,
	stargateClient api.StargateServiceClient,
	clientConfig Config,
	workerPoolConfig WorkerPoolConfig,
	metricsRegistry metrics.MetricsRegistry,
) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = clientConfig.HTTPMaxIdleConns
	transport.MaxIdleConnsPerHost = clientConfig.HTTPMaxIdleConnsPerHost

	clientConfig.GatewayUrl = strings.TrimSuffix(clientConfig.GatewayUrl, "/")

	return &Client{
		logger:          l,
		stargateClient:  stargateClient,
		config:          clientConfig,
		workerPool:      pond.New(workerPoolConfig.MaxWorkers, workerPoolConfig.MaxTasks),
		metricsRegistry: metricsRegistry,
		httpClient: &http.Client{
			Timeout:   clientConfig.HTTPClientTimeout,
			Transport: transport,
		},
	}
}

type ResponseChanEvent struct {
	msg *api.StargateClientMessage
	err error
}

func (c *Client) Run(ctx context.Context) error {
	c.logger.Info("starting client...")

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"organization-id", c.config.OrganizationID,
		"stack-id", c.config.StackID,
	)

	c.logger.WithFields(map[string]any{
		"organization_id": c.config.OrganizationID,
		"stack_id":        c.config.StackID,
	}).Info("connecting to stargate server...")

	stream, err := c.stargateClient.Stargate(ctx)
	if err != nil {
		return err
	}

	c.logger.WithFields(map[string]any{
		"organization_id": c.config.OrganizationID,
		"stack_id":        c.config.StackID,
	}).Info("connected to stargate server")
	defer c.logger.WithFields(map[string]any{
		"organization_id": c.config.OrganizationID,
		"stack_id":        c.config.StackID,
	}).Info("disconnected from stargate server")

	responseChan := make(chan *ResponseChanEvent, c.config.ChanSize)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		for {
			in, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}

				return err
			}

			c.logger.WithFields(map[string]any{
				"event": in,
			}).Debug("received message from server")

			c.workerPool.Submit(func() {
				out := c.Forward(ctx, in)
				select {
				case <-ctx.Done():
					return
				case responseChan <- out:
				}
			})
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case response := <-responseChan:
				if response.err != nil {
					// Note: how should we handle errors here?
					return response.err
				}

				if response.msg == nil {
					continue
				}

				c.logger.WithFields(map[string]any{
					"response": response,
				}).Debug("sending response message to server")

				err := stream.Send(response.msg)
				if err != nil {
					return err
				}
			}
		}
	})

	return eg.Wait()
}

func (c *Client) Forward(ctx context.Context, in *api.StargateServerMessage) *ResponseChanEvent {
	attrs := []attribute.KeyValue{}

	switch ev := in.Event.(type) {
	case *api.StargateServerMessage_ApiCall:

		ctx = opentelemetry.Propagator.Extract(ctx, propagation.MapCarrier(ev.ApiCall.OtlpContext))

		attrs = append(attrs, attribute.String("message_type", "api_call"))
		c.metricsRegistry.ServerMessageReceivedByType().Add(ctx, 1, metric.WithAttributes(attrs...))

		attrs = append(attrs, attribute.String("path", ev.ApiCall.Path))
		path := strings.TrimPrefix(ev.ApiCall.Path, "/")

		req, err := http.NewRequestWithContext(ctx, ev.ApiCall.Method, c.config.GatewayUrl+"/"+path, bytes.NewReader(ev.ApiCall.Body))
		if err != nil {
			return &ResponseChanEvent{
				err: err,
			}
		}

		opentelemetry.Propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

		q := req.URL.Query()
		for k, v := range ev.ApiCall.Query {
			for _, vv := range v.Values {
				q.Add(k, vv)
			}
		}
		req.URL.RawQuery = q.Encode()

		for k, v := range ev.ApiCall.Headers {
			for _, vv := range v.Values {
				req.Header.Add(k, vv)
			}
		}

		now := time.Now()
		resp, err := c.httpClient.Do(req)
		if err != nil {
			c.logger.Errorf("error making http request: %v", err)
			return &ResponseChanEvent{
				err: nil,
				msg: &api.StargateClientMessage{
					CorrelationId: in.CorrelationId,
					Event: &api.StargateClientMessage_ApiCallResponse{ApiCallResponse: &api.StargateClientMessage_APICallResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       []byte{},
						Headers:    map[string]*api.Values{},
					}},
				},
			}
		}
		latency := time.Since(now)
		c.metricsRegistry.HTTPCallLatencies().Record(ctx, latency.Milliseconds(), metric.WithAttributes(attrs...))

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return &ResponseChanEvent{
				err: err,
			}
		}

		headers := make(map[string]*api.Values)
		for k, v := range resp.Header {
			headers[k] = &api.Values{
				Values: v,
			}
		}

		attrs = append(attrs, attribute.Int("status_code", resp.StatusCode))
		c.metricsRegistry.HTTPCallStatusCodes().Add(ctx, 1, metric.WithAttributes(attrs...))

		return &ResponseChanEvent{
			err: nil,
			msg: &api.StargateClientMessage{
				CorrelationId: in.CorrelationId,
				Event: &api.StargateClientMessage_ApiCallResponse{ApiCallResponse: &api.StargateClientMessage_APICallResponse{
					StatusCode: int32(resp.StatusCode),
					Body:       body,
					Headers:    headers,
				}},
			},
		}
	case *api.StargateServerMessage_Ping_:
		return &ResponseChanEvent{
			err: nil,
			msg: &api.StargateClientMessage{
				CorrelationId: in.CorrelationId,
				Event: &api.StargateClientMessage_Pong_{
					Pong: &api.StargateClientMessage_Pong{},
				},
			},
		}
	}

	return &ResponseChanEvent{
		err: nil,
		msg: nil,
	}
}

func (c *Client) Close() error {
	c.workerPool.StopAndWait()
	return nil
}
