package engine

import "go.opentelemetry.io/otel"

var Tracer = otel.Tracer("connectors")
