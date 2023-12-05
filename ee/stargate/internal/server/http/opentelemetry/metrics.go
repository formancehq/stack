package opentelemetry

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

type MetricsRegistry interface {
	ReceivedHTTPCallByPath() metric.Int64Counter
}

type metricsRegistry struct {
	receivedHTTPCallByPath metric.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("server_http")

	receivedHTTPCallByPath, err := meter.Int64Counter(
		"stargate_server_received_http_call_by_path",
		metric.WithUnit("1"),
		metric.WithDescription("Received HTTP call by path"),
	)
	if err != nil {
		return nil, err
	}

	return &metricsRegistry{
		receivedHTTPCallByPath: receivedHTTPCallByPath,
	}, nil
}

func (m *metricsRegistry) ReceivedHTTPCallByPath() metric.Int64Counter {
	return m.receivedHTTPCallByPath
}

type NoOpMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoOpMetricsRegistry {
	return &NoOpMetricsRegistry{}
}

func (m *NoOpMetricsRegistry) ReceivedHTTPCallByPath() metric.Int64Counter {
	counter, _ := noop.NewMeterProvider().Meter("server_http").Int64Counter("stargate_server_received_http_call_by_path")
	return counter
}
