package workflow

import (
	"go.opentelemetry.io/otel"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
)

var Tracer = otel.Tracer("runner")
var TracerOptions = opentelemetry.TracerOptions{
	Tracer: Tracer,
}
var TemporalTracer interceptor.Tracer

func init() {
	var err error
	TemporalTracer, err = opentelemetry.NewTracer(TracerOptions)
	if err != nil {
		panic(err)
	}
}
