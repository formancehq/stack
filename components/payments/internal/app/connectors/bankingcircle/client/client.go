package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient *http.Client

	username string
	password string

	endpoint              string
	authorizationEndpoint string

	logger logging.Logger

	accessToken          string
	accessTokenExpiresAt time.Time
}

func newHTTPClient(userCertificate, userCertificateKey string) (*http.Client, error) {
	cert, err := tls.X509KeyPair([]byte(userCertificate), []byte(userCertificateKey))
	if err != nil {
		return nil, err
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: otelhttp.NewTransport(tr),
	}, nil
}

func NewClient(
	username, password,
	endpoint, authorizationEndpoint,
	uCertificate, uCertificateKey string,
	logger logging.Logger) (*Client, error) {
	httpClient, err := newHTTPClient(uCertificate, uCertificateKey)
	if err != nil {
		return nil, err
	}

	c := &Client{
		httpClient: httpClient,

		username:              username,
		password:              password,
		endpoint:              endpoint,
		authorizationEndpoint: authorizationEndpoint,

		logger: logger,
	}

	if err := c.login(context.TODO()); err != nil {
		return nil, err
	}

	return c, nil
}
