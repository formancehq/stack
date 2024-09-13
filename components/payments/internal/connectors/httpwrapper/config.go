package httpwrapper

import (
	"net/http"
	"time"
)

type Config struct {
	HttpErrorCheckerFn func(code int) error

	Timeout   time.Duration
	Transport http.RoundTripper
}
