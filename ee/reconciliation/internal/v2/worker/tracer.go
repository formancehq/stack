package worker

import (
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("worker")
