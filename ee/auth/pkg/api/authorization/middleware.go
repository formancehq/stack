package authorization

import (
	"errors"
	"net/http"

	"github.com/zitadel/oidc/v2/pkg/op"
)

var (
	ErrMissingAuthHeader   = errors.New("missing authorization header")
	ErrMalformedAuthHeader = errors.New("malformed authorization header")
	ErrVerifyAuthToken     = errors.New("could not verify access token")
)

func Middleware(o op.OpenIDProvider) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if err := verifyAccessToken(r, o); err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				h.ServeHTTP(w, r)
			})
	}
}
