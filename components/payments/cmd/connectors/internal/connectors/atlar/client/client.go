package client

import (
	"net/url"

	"github.com/go-openapi/strfmt"

	atlar_client "github.com/get-momo/atlar-v1-go-client/client"

	httptransport "github.com/go-openapi/runtime/client"
)

type Client struct {
	client *atlar_client.Rest
}

func NewClient(baseURL url.URL, accessKey, secret string) *Client {
	return &Client{
		client: createAtlarClient(baseURL, accessKey, secret),
	}
}

func createAtlarClient(baseURL url.URL, accessKey, secret string) *atlar_client.Rest {
	transport := httptransport.New(
		baseURL.Host,
		baseURL.Path,
		[]string{baseURL.Scheme},
	)
	basicAuth := httptransport.BasicAuth(accessKey, secret)
	transport.DefaultAuthentication = basicAuth
	client := atlar_client.New(transport, strfmt.Default)
	return client
}
