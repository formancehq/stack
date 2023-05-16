package opentelemetry

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type MetricsRegistry interface {
	ReceivedHTTPCallByPath() instrument.Int64Counter
}

type metricsRegistry struct {
	receivedHTTPCallByPath instrument.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("server_http")

	receivedHTTPCallByPath, err := meter.Int64Counter(
		"received_http_call_by_path",
		instrument.WithUnit("1"),
		instrument.WithDescription("Received HTTP call by path"),
	)
	if err != nil {
		return nil, err
	}

	return &metricsRegistry{
		receivedHTTPCallByPath: receivedHTTPCallByPath,
	}, nil
}

func (m *metricsRegistry) ReceivedHTTPCallByPath() instrument.Int64Counter {
	return m.receivedHTTPCallByPath
}

type NoOpMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoOpMetricsRegistry {
	return &NoOpMetricsRegistry{}
}

func (m *NoOpMetricsRegistry) ReceivedHTTPCallByPath() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("received_http_call_by_path")
	return counter
}
