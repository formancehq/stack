package service

import (
	"net/http"
	"net/http/httputil"

	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	data []byte
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.data = append(w.data, data...)
	return w.ResponseWriter.Write(w.data)
}

func OTLPMiddleware(serverName string) func(h http.Handler) http.Handler {
	m := otelchi.Middleware(serverName)
	return func(h http.Handler) http.Handler {
		return m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if IsDebug(r.Context()) {
				data, err := httputil.DumpRequest(r, true)
				if err != nil {
					panic(err)
				}
				trace.SpanFromContext(r.Context()).
					SetAttributes(attribute.String("http.request", string(data)))

				rw := &responseWriter{w, make([]byte, 0, 1024)}
				defer func() {
					trace.SpanFromContext(r.Context()).
						SetAttributes(attribute.String("http.response", string(rw.data)))
				}()

				h.ServeHTTP(rw, r)
				return
			}

			h.ServeHTTP(w, r)
		}))
	}
}
