package otel

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	once   sync.Once
	tracer trace.Tracer
)

func Tracer() trace.Tracer {
	once.Do(func() {
		tracer = otel.Tracer("com.formance.payments")
	})

	return tracer
}

func RecordError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
