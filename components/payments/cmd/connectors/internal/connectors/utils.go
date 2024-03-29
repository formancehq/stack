package connectors

import (
	"context"
	"os"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DeferrableFunc func(ctx context.Context, timeSince time.Time)

func ClientMetrics(ctx context.Context, connectorName, operation string) DeferrableFunc {
	attributes := []attribute.KeyValue{
		attribute.String("connector", connectorName),
		attribute.String("operation", operation),
	}

	stack := os.Getenv("STACK")
	if stack != "" {
		attributes = append(attributes, attribute.String("stack", stack))
	}

	metrics.GetMetricsRegistry().ConnectorPSPCalls().Add(ctx, 1, metric.WithAttributes(attributes...))

	return func(ctx context.Context, timeSince time.Time) {
		metrics.GetMetricsRegistry().ConnectorPSPCallLatencies().Record(ctx, time.Since(timeSince).Milliseconds(), metric.WithAttributes(attributes...))
	}
}

func StartSpan(
	ctx context.Context,
	spanName string,
	attributes ...attribute.KeyValue,
) (context.Context, trace.Span) {
	parentSpan := trace.SpanFromContext(ctx)
	return otel.Tracer().Start(
		ctx,
		spanName,
		trace.WithNewRoot(),
		trace.WithLinks(trace.Link{
			SpanContext: parentSpan.SpanContext(),
		}),
		trace.WithAttributes(
			attributes...,
		),
	)
}
