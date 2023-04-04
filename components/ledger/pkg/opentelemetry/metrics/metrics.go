package metrics

import (
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type GlobalMetricsRegistry struct {
	// API Latencies
	APILatencies instrument.Int64Histogram
	StatusCodes  instrument.Int64Counter
}

func RegisterGlobalMetricsRegistry(meterProvider *sdkmetric.MeterProvider) (*GlobalMetricsRegistry, error) {
	meter := meterProvider.Meter("global")

	apiLatencies, err := meter.Int64Histogram(
		"api_latencies",
		instrument.WithUnit("ms"),
		instrument.WithDescription("Latency of API calls"),
	)
	if err != nil {
		return nil, err
	}

	statusCodes, err := meter.Int64Counter(
		"status_codes",
		instrument.WithUnit("1"),
		instrument.WithDescription("Status codes of API calls"),
	)
	if err != nil {
		return nil, err
	}

	return &GlobalMetricsRegistry{
		APILatencies: apiLatencies,
		StatusCodes:  statusCodes,
	}, nil
}

type PerLedgerMetricsRegistry struct {
}

func RegisterPerLedgerMetricsRegistry(ledger string) (*PerLedgerMetricsRegistry, error) {
	// we can now use the global meter provider to create a meter
	// since it was created by the fx
	meter := global.MeterProvider().Meter(ledger)
	_ = meter

	return &PerLedgerMetricsRegistry{
		// APILatencies: apiLatencies,
	}, nil
}
