package metrics

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	ClientsConnected atomic.Int64

	// TODO(polo): this is not ideal, we should be able to have the pod name
	// instead
	attrs []attribute.KeyValue
)

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	n := source.Int63()

	attrs = append(attrs, attribute.Int64("pod_id", n))
}

type MetricsRegistry interface {
	UnAuthenticatedCalls() metric.Int64Counter
	ClientsConnected() metric.Int64ObservableGauge
	StreamErrors() metric.Int64Counter
	GRPCLatencies() metric.Int64Histogram
	CorrelationIDNotFound() metric.Int64Counter
}

type metricsRegistry struct {
	unAuthenticatedCalls  metric.Int64Counter
	clientsConnected      metric.Int64ObservableGauge
	streamErrors          metric.Int64Counter
	grpcLatencies         metric.Int64Histogram
	correlationIDNotFound metric.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("server_grpc")

	unAuthenticatedCalls, err := meter.Int64Counter(
		"stargate_server_unauthenticated_calls",
		metric.WithUnit("1"),
		metric.WithDescription("Unauthenticated calls"),
	)
	if err != nil {
		return nil, err
	}

	clientsConnected, err := meter.Int64ObservableGauge(
		"stargate_server_clients_connected",
		metric.WithUnit("1"),
		metric.WithDescription("Number of connected clients"),
		metric.WithInt64Callback(func(ctx context.Context, obs metric.Int64Observer) error {
			obs.Observe(ClientsConnected.Load(), metric.WithAttributes(attrs...))
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	streamErrors, err := meter.Int64Counter(
		"stargate_server_stream_errors",
		metric.WithUnit("1"),
		metric.WithDescription("Stream errors"),
	)
	if err != nil {
		return nil, err
	}

	grpcLatencies, err := meter.Int64Histogram(
		"stargate_server_grpc_latencies",
		metric.WithUnit("ms"),
		metric.WithDescription("Latency of gRPC calls"),
	)
	if err != nil {
		return nil, err
	}

	correlationIDNotFound, err := meter.Int64Counter(
		"stargate_server_correlation_id_not_found",
		metric.WithUnit("1"),
		metric.WithDescription("Correlation ID not found"),
	)
	if err != nil {
		return nil, err
	}

	return &metricsRegistry{
		unAuthenticatedCalls:  unAuthenticatedCalls,
		clientsConnected:      clientsConnected,
		streamErrors:          streamErrors,
		grpcLatencies:         grpcLatencies,
		correlationIDNotFound: correlationIDNotFound,
	}, nil
}

func (m *metricsRegistry) UnAuthenticatedCalls() metric.Int64Counter {
	return m.unAuthenticatedCalls
}

func (m *metricsRegistry) GRPCLatencies() metric.Int64Histogram {
	return m.grpcLatencies
}

func (m *metricsRegistry) CorrelationIDNotFound() metric.Int64Counter {
	return m.correlationIDNotFound
}

func (m *metricsRegistry) ClientsConnected() metric.Int64ObservableGauge {
	return m.clientsConnected
}

func (m *metricsRegistry) StreamErrors() metric.Int64Counter {
	return m.streamErrors
}

type NoOpMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoOpMetricsRegistry {
	return &NoOpMetricsRegistry{}
}

func (m *NoOpMetricsRegistry) UnAuthenticatedCalls() metric.Int64Counter {
	counter, _ := noop.NewMeterProvider().Meter("server_grpc").Int64Counter("stargate_server_unauthenticated_calls")
	return counter
}

func (m *NoOpMetricsRegistry) GRPCLatencies() metric.Int64Histogram {
	histogram, _ := noop.NewMeterProvider().Meter("server_grpc").Int64Histogram("stargate_server_grpc_latencies")
	return histogram
}

func (m *NoOpMetricsRegistry) CorrelationIDNotFound() metric.Int64Counter {
	counter, _ := noop.NewMeterProvider().Meter("server_grpc").Int64Counter("stargate_server_correlation_id_not_found")
	return counter
}

func (m *NoOpMetricsRegistry) ClientsConnected() metric.Int64ObservableGauge {
	counter, _ := noop.NewMeterProvider().Meter("server_grpc").Int64ObservableGauge("stargate_server_clients_connected")
	return counter
}

func (m *NoOpMetricsRegistry) StreamErrors() metric.Int64Counter {
	counter, _ := noop.NewMeterProvider().Meter("server_grpc").Int64Counter("stargate_server_stream_errors")
	return counter
}
