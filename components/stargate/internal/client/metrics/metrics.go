package metrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type MetricsRegistry interface {
	HTTPCallLatencies() instrument.Int64Histogram
	HTTPCallStatusCodes() instrument.Int64Counter
	ServerMessageReceivedByType() instrument.Int64Counter
}

type metricsRegistry struct {
	httpCallLatencies           instrument.Int64Histogram
	httpCallStatusCodes         instrument.Int64Counter
	serverMessageReceivedByType instrument.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("client")

	httpCallLatencies, err := meter.Int64Histogram(
		"http_call_latencies",
		instrument.WithUnit("ms"),
		instrument.WithDescription("Latency of HTTP calls"),
	)
	if err != nil {
		return nil, err
	}

	httpCallStatusCodes, err := meter.Int64Counter(
		"http_call_status_codes",
		instrument.WithUnit("1"),
		instrument.WithDescription("HTTP status codes of HTTP calls"),
	)
	if err != nil {
		return nil, err
	}

	serverMessageReceivedByType, err := meter.Int64Counter(
		"server_message_received_by_type",
		instrument.WithUnit("1"),
		instrument.WithDescription("Server message received by type"),
	)
	if err != nil {
		return nil, err
	}

	return &metricsRegistry{
		httpCallLatencies:           httpCallLatencies,
		httpCallStatusCodes:         httpCallStatusCodes,
		serverMessageReceivedByType: serverMessageReceivedByType,
	}, nil
}

func (m *metricsRegistry) HTTPCallLatencies() instrument.Int64Histogram {
	return m.httpCallLatencies
}

func (m *metricsRegistry) HTTPCallStatusCodes() instrument.Int64Counter {
	return m.httpCallStatusCodes
}

func (m *metricsRegistry) ServerMessageReceivedByType() instrument.Int64Counter {
	return m.serverMessageReceivedByType
}

type NoOpMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoOpMetricsRegistry {
	return &NoOpMetricsRegistry{}
}

func (m *NoOpMetricsRegistry) HTTPCallLatencies() instrument.Int64Histogram {
	histogram, _ := metric.NewNoopMeter().Int64Histogram("http_call_latencies")
	return histogram
}

func (m *NoOpMetricsRegistry) HTTPCallStatusCodes() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("http_call_status_codes")
	return counter
}

func (m *NoOpMetricsRegistry) ServerMessageReceivedByType() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("server_message_received_by_type")
	return counter
}
