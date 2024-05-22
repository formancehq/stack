package oidc

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/op"
)

const AuthorizeCallbackPath = "/authorize/callback"

func AddRoutes(r chi.Router, provider op.OpenIDProvider, storage Storage, relyingParty rp.RelyingParty) {
	r.Group(func(r chi.Router) {
		if relyingParty != nil {
			r.Use(func(h http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == AuthorizeCallbackPath {
						if code := r.URL.Query().Get("code"); code != "" {
							authorizeCallbackHandler(provider, storage, relyingParty).ServeHTTP(w, r)
							return
						} else if err := r.URL.Query().Get("error"); err != "" {
							authorizeErrorHandler().ServeHTTP(w, r)
							return
						}
					}
					h.ServeHTTP(w, r)
				})
			})
		}

		r.
			With(func(handler http.Handler) http.Handler {
				// The otelchi middleware does not see matching route as it is matched in a subrouter
				// So the span name terminated with just "/"
				// This middleware make the hack
				// We can do this because url does not contain any dynamic variables.
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					currentSpan := trace.SpanFromContext(r.Context())
					currentSpan.SetName(r.URL.Path)
					currentSpan.SetAttributes(attribute.String("http.route", r.URL.Path))
					handler.ServeHTTP(w, r)
				})
			}).
			Mount("/", provider.HttpHandler())
	})
}
