package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type Client struct {
	httpClient httpwrapper.Client

	username string
	password string

	endpoint              string
	authorizationEndpoint string

	accessToken          string
	accessTokenExpiresAt time.Time
}

func New(
	username, password,
	endpoint, authorizationEndpoint,
	uCertificate, uCertificateKey string,
) (*Client, error) {
	cert, err := tls.X509KeyPair([]byte(uCertificate), []byte(uCertificateKey))
	if err != nil {
		return nil, err
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	config := &httpwrapper.Config{
		Transport: tr,
	}
	httpClient, err := httpwrapper.NewClient(config)
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
