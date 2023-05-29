package opentelemetry

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

var (
	ClientsConnected atomic.Int64
)

type MetricsRegistry interface {
	UnAuthenticatedCalls() instrument.Int64Counter
	ClientsConnected() instrument.Int64ObservableGauge
	StreamErrors() instrument.Int64Counter
	GRPCLatencies() instrument.Int64Histogram
	CorrelationIDNotFound() instrument.Int64Counter
}

type metricsRegistry struct {
	unAuthenticatedCalls  instrument.Int64Counter
	clientsConnected      instrument.Int64ObservableGauge
	streamErrors          instrument.Int64Counter
	grpcLatencies         instrument.Int64Histogram
	correlationIDNotFound instrument.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("server_grpc")

	unAuthenticatedCalls, err := meter.Int64Counter(
		"stargate_server_unauthenticated_calls",
		instrument.WithUnit("1"),
		instrument.WithDescription("Unauthenticated calls"),
	)
	if err != nil {
		return nil, err
	}

	clientsConnected, err := meter.Int64ObservableGauge(
		"stargate_server_clients_connected",
		instrument.WithUnit("1"),
		instrument.WithDescription("Number of connected clients"),
		instrument.WithInt64Callback(func(ctx context.Context, obs instrument.Int64Observer) error {
			obs.Observe(ClientsConnected.Load())
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	streamErrors, err := meter.Int64Counter(
		"stargate_server_stream_errors",
		instrument.WithUnit("1"),
		instrument.WithDescription("Stream errors"),
	)
	if err != nil {
		return nil, err
	}

	grpcLatencies, err := meter.Int64Histogram(
		"stargate_server_grpc_latencies",
		instrument.WithUnit("ms"),
		instrument.WithDescription("Latency of gRPC calls"),
	)
	if err != nil {
		return nil, err
	}

	correlationIDNotFound, err := meter.Int64Counter(
		"stargate_server_correlation_id_not_found",
		instrument.WithUnit("1"),
		instrument.WithDescription("Correlation ID not found"),
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

func (m *metricsRegistry) UnAuthenticatedCalls() instrument.Int64Counter {
	return m.unAuthenticatedCalls
}

func (m *metricsRegistry) GRPCLatencies() instrument.Int64Histogram {
	return m.grpcLatencies
}

func (m *metricsRegistry) CorrelationIDNotFound() instrument.Int64Counter {
	return m.correlationIDNotFound
}

func (m *metricsRegistry) ClientsConnected() instrument.Int64ObservableGauge {
	return m.clientsConnected
}

func (m *metricsRegistry) StreamErrors() instrument.Int64Counter {
	return m.streamErrors
}

type NoOpMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoOpMetricsRegistry {
	return &NoOpMetricsRegistry{}
}

func (m *NoOpMetricsRegistry) UnAuthenticatedCalls() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("stargate_server_unauthenticated_calls")
	return counter
}

func (m *NoOpMetricsRegistry) GRPCLatencies() instrument.Int64Histogram {
	histogram, _ := metric.NewNoopMeter().Int64Histogram("stargate_server_grpc_latencies")
	return histogram
}

func (m *NoOpMetricsRegistry) CorrelationIDNotFound() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("stargate_server_correlation_id_not_found")
	return counter
}

func (m *NoOpMetricsRegistry) ClientsConnected() instrument.Int64ObservableGauge {
	counter, _ := metric.NewNoopMeter().Int64ObservableGauge("stargate_server_clients_connected")
	return counter
}

func (m *NoOpMetricsRegistry) StreamErrors() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("stargate_server_stream_errors")
	return counter
}
