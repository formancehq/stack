package httpclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/stack/ee/ingester/internal/drivers"

	"github.com/formancehq/stack/libs/go-libs/httpclient"

	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/zitadel/oidc/v2/pkg/client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type OAuth2Config struct {
	Issuer       string
	ClientID     string
	ClientSecret string
}

type StackAuthenticatedClient struct {
	*http.Client
}

func NewStackAuthenticatedClientFromHTTPClient(client *http.Client) *StackAuthenticatedClient {
	return &StackAuthenticatedClient{
		Client: client,
	}
}

// NewStackAuthenticatedClient create a new http client instrumentalized with open telemetry.
// It also adds OAuth2 instrumentation if the stackClientID parameter is not empty.
// If stackClientID parameter is not empty, stackIssuer and stackClientSecret must be set.
func NewStackAuthenticatedClient(authConfig OAuth2Config, httpClient *http.Client, modules ...string) (*StackAuthenticatedClient, error) {

	if authConfig.ClientID == "" {
		return NewStackAuthenticatedClientFromHTTPClient(httpClient), nil
	}
	scopes := []string{"openid"}
	for _, module := range modules {
		scopes = append(scopes, fmt.Sprintf("%s:read", module))
	}

	discovery, err := client.Discover(authConfig.Issuer, httpClient)
	if err != nil {
		return nil, err
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		TokenURL:     discovery.TokenEndpoint,
		Scopes:       scopes,
	}
	return NewStackAuthenticatedClientFromHTTPClient(oauthConfig.Client(context.WithValue(context.Background(), oauth2.HTTPClient, httpClient))), nil
}

// NewClient create a new http client instrumentalized with open telemetry.
func NewClient(serviceConfig drivers.ServiceConfig) *http.Client {
	httpClient := &http.Client{
		Transport: otlp.NewRoundTripper(http.DefaultTransport, serviceConfig.Debug),
	}
	if serviceConfig.Debug {
		httpClient.Transport = httpclient.NewDebugHTTPTransport(httpClient.Transport)
	}

	return httpClient
}
