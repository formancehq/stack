package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient *http.Client

	username string
	password string

	endpoint              string
	authorizationEndpoint string

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

func New(
	username, password,
	endpoint, authorizationEndpoint,
	uCertificate, uCertificateKey string,
) (*Client, error) {
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
	}

	return c, nil
}
