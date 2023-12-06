package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/hmac"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type apiTransport struct {
	apiKey     string
	headers    map[string]string
	underlying http.RoundTripper
}

func (t *apiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", t.apiKey)

	return t.underlying.RoundTrip(req)
}

type responseWrapper[t any] struct {
	Content    t   `json:"content"`
	Size       int `json:"size"`
	TotalSize  int `json:"totalSize"`
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
}

type Client struct {
	httpClient *http.Client
	endpoint   string
}

func (m *Client) buildEndpoint(path string, args ...interface{}) string {
	endpoint := strings.TrimSuffix(m.endpoint, "/")
	return fmt.Sprintf("%s/%s", endpoint, fmt.Sprintf(path, args...))
}

const sandboxAPIEndpoint = "https://api-sandbox.modulrfinance.com/api-sandbox-token"

func NewClient(apiKey, apiSecret, endpoint string) (*Client, error) {
	if endpoint == "" {
		endpoint = sandboxAPIEndpoint
	}

	headers, err := hmac.GenerateHeaders(apiKey, apiSecret, "", false)
	if err != nil {
		return nil, fmt.Errorf("failed to generate headers: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Transport: &apiTransport{
				headers:    headers,
				apiKey:     apiKey,
				underlying: otelhttp.NewTransport(http.DefaultTransport),
			},
		},
		endpoint: endpoint,
	}, nil
}

type ErrorResponse struct {
	Field         string `json:"field"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	ErrorCode     string `json:"errorCode"`
	SourceService string `json:"sourceService"`
}
