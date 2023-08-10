package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type MetricsRegistry interface {
	HTTPCallLatencies() metric.Int64Histogram
	HTTPCallStatusCodes() metric.Int64Counter
	ServerMessageReceivedByType() metric.Int64Counter
}

type metricsRegistry struct {
	httpCallLatencies           metric.Int64Histogram
	httpCallStatusCodes         metric.Int64Counter
	serverMessageReceivedByType metric.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("client")

	httpCallLatencies, err := meter.Int64Histogram(
		"http_call_latencies",
		metric.WithUnit("ms"),
		metric.WithDescription("Latency of HTTP calls"),
	)
	if err != nil {
		return nil, err
	}

	httpCallStatusCodes, err := meter.Int64Counter(
		"http_call_status_codes",
		metric.WithUnit("1"),
		metric.WithDescription("HTTP status codes of HTTP calls"),
	)
	if err != nil {
		return nil, err
	}

	serverMessageReceivedByType, err := meter.Int64Counter(
		"server_message_received_by_type",
		metric.WithUnit("1"),
		metric.WithDescription("Server message received by type"),
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

func (m *metricsRegistry) HTTPCallLatencies() metric.Int64Histogram {
	return m.httpCallLatencies
}

func (m *metricsRegistry) HTTPCallStatusCodes() metric.Int64Counter {
	return m.httpCallStatusCodes
}

func (m *metricsRegistry) ServerMessageReceivedByType() metric.Int64Counter {
	return m.serverMessageReceivedByType
}

type NoOpMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoOpMetricsRegistry {
	return &NoOpMetricsRegistry{}
}

func (m *NoOpMetricsRegistry) HTTPCallLatencies() metric.Int64Histogram {
	histogram, _ := otel.GetMeterProvider().Meter("client").Int64Histogram("http_call_latencies")
	return histogram
}

func (m *NoOpMetricsRegistry) HTTPCallStatusCodes() metric.Int64Counter {
	counter, _ := otel.GetMeterProvider().Meter("client").Int64Counter("http_call_status_codes")
	return counter
}

func (m *NoOpMetricsRegistry) ServerMessageReceivedByType() metric.Int64Counter {
	counter, _ := otel.GetMeterProvider().Meter("client").Int64Counter("server_message_received_by_type")
	return counter
}
