package httpwrapper

import (
	"net/http"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

type Config struct {
	HttpErrorCheckerFn func(code int) error

	Timeout     time.Duration
	Transport   http.RoundTripper
	OAuthConfig *clientcredentials.Config
}
