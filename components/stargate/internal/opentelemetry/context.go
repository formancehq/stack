package opentelemetry

import (
	"go.opentelemetry.io/otel/propagation"
)

var (
	Propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
)
