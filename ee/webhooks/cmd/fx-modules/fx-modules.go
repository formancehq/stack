package fxmodules

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"

	"go.opentelemetry.io/otel"
)

var Tracer = otel.Tracer("webhook")

func FxProvideHttpClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			str := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			if len(r.URL.Query()) == 0 {
				return str
			}
			return fmt.Sprintf("%s?%s", str, r.URL.Query().Encode())
		})),
	}

}
