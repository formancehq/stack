package metrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

const (
	ObjectAttributeKey = "object"
)

type MetricsRegistry interface {
	ConnectorObjects() instrument.Int64Counter
	ConnectorObjectsLatency() instrument.Int64Histogram
	ConnectorObjectsErrors() instrument.Int64Counter
}

type metricsRegistry struct {
	connectorObjects        instrument.Int64Counter
	connectorObjectsLatency instrument.Int64Histogram
	connectorObjectsErrors  instrument.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("payments")

	connectorObjects, err := meter.Int64Counter(
		"payments_connectors_objects",
		instrument.WithUnit("1"),
		instrument.WithDescription("Object fetch from connectors (accounts, payments, balances, ...)"),
	)
	if err != nil {
		return nil, err
	}

	connectorObjectLatencies, err := meter.Int64Histogram(
		"payments_connectors_object_latencies",
		instrument.WithUnit("ms"),
		instrument.WithDescription("Object latencies from connectors (accounts, payments, balances, ...)"),
	)
	if err != nil {
		return nil, err
	}

	connectorObjectErrors, err := meter.Int64Counter(
		"payments_connectors_object_errors",
		instrument.WithUnit("1"),
		instrument.WithDescription("Obejct errors from connectors (accounts, payments, balances, ...)"),
	)
	if err != nil {
		return nil, err
	}

	return &metricsRegistry{
		connectorObjects:        connectorObjects,
		connectorObjectsLatency: connectorObjectLatencies,
		connectorObjectsErrors:  connectorObjectErrors,
	}, nil
}

func (m *metricsRegistry) ConnectorObjects() instrument.Int64Counter {
	return m.connectorObjects
}

func (m *metricsRegistry) ConnectorObjectsLatency() instrument.Int64Histogram {
	return m.connectorObjectsLatency
}

func (m *metricsRegistry) ConnectorObjectsErrors() instrument.Int64Counter {
	return m.connectorObjectsErrors
}

type NoopMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoopMetricsRegistry {
	return &NoopMetricsRegistry{}
}

func (m *NoopMetricsRegistry) ConnectorObjects() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("payments_connectors_objects")
	return counter
}

func (m *NoopMetricsRegistry) ConnectorObjectsLatency() instrument.Int64Histogram {
	histogram, _ := metric.NewNoopMeter().Int64Histogram("payments_connectors_object_latencies")
	return histogram
}

func (m *NoopMetricsRegistry) ConnectorObjectsErrors() instrument.Int64Counter {
	counter, _ := metric.NewNoopMeter().Int64Counter("payments_connectors_object_errors")
	return counter
}
