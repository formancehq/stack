package metrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const (
	ObjectAttributeKey = "object"
)

type MetricsRegistry interface {
	ConnectorCurrencyNotSupported() metric.Int64Counter
	ConnectorObjects() metric.Int64Counter
	ConnectorObjectsLatency() metric.Int64Histogram
	ConnectorObjectsErrors() metric.Int64Counter
}

type metricsRegistry struct {
	connectorCurrencyNotSupported metric.Int64Counter
	connectorObjects              metric.Int64Counter
	connectorObjectsLatency       metric.Int64Histogram
	connectorObjectsErrors        metric.Int64Counter
}

func RegisterMetricsRegistry(meterProvider metric.MeterProvider) (MetricsRegistry, error) {
	meter := meterProvider.Meter("payments")

	connectorCurrencyNotSupported, err := meter.Int64Counter(
		"payments_connectors_currency_not_supported",
		metric.WithUnit("1"),
		metric.WithDescription("Currency not supported by connector"),
	)
	if err != nil {
		return nil, err
	}

	connectorObjects, err := meter.Int64Counter(
		"payments_connectors_objects",
		metric.WithUnit("1"),
		metric.WithDescription("Object fetch from connectors (accounts, payments, balances, ...)"),
	)
	if err != nil {
		return nil, err
	}

	connectorObjectLatencies, err := meter.Int64Histogram(
		"payments_connectors_object_latencies",
		metric.WithUnit("ms"),
		metric.WithDescription("Object latencies from connectors (accounts, payments, balances, ...)"),
	)
	if err != nil {
		return nil, err
	}

	connectorObjectErrors, err := meter.Int64Counter(
		"payments_connectors_object_errors",
		metric.WithUnit("1"),
		metric.WithDescription("Obejct errors from connectors (accounts, payments, balances, ...)"),
	)
	if err != nil {
		return nil, err
	}

	return &metricsRegistry{
		connectorCurrencyNotSupported: connectorCurrencyNotSupported,
		connectorObjects:              connectorObjects,
		connectorObjectsLatency:       connectorObjectLatencies,
		connectorObjectsErrors:        connectorObjectErrors,
	}, nil
}

func (m *metricsRegistry) ConnectorCurrencyNotSupported() metric.Int64Counter {
	return m.connectorCurrencyNotSupported
}

func (m *metricsRegistry) ConnectorObjects() metric.Int64Counter {
	return m.connectorObjects
}

func (m *metricsRegistry) ConnectorObjectsLatency() metric.Int64Histogram {
	return m.connectorObjectsLatency
}

func (m *metricsRegistry) ConnectorObjectsErrors() metric.Int64Counter {
	return m.connectorObjectsErrors
}

type NoopMetricsRegistry struct{}

func NewNoOpMetricsRegistry() *NoopMetricsRegistry {
	return &NoopMetricsRegistry{}
}

func (m *NoopMetricsRegistry) ConnectorCurrencyNotSupported() metric.Int64Counter {
	counter, _ := noop.NewMeterProvider().Meter("payments").Int64Counter("payments_connectors_currency_not_supported")
	return counter
}

func (m *NoopMetricsRegistry) ConnectorObjects() metric.Int64Counter {
	counter, _ := noop.NewMeterProvider().Meter("payments").Int64Counter("payments_connectors_objects")
	return counter
}

func (m *NoopMetricsRegistry) ConnectorObjectsLatency() metric.Int64Histogram {
	histogram, _ := noop.NewMeterProvider().Meter("payments").Int64Histogram("payments_connectors_object_latencies")
	return histogram
}

func (m *NoopMetricsRegistry) ConnectorObjectsErrors() metric.Int64Counter {
	counter, _ := noop.NewMeterProvider().Meter("payments").Int64Counter("payments_connectors_object_errors")
	return counter
}
