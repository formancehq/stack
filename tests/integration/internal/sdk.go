package internal

import (
	"net/http"
	"net/url"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/httpclient"
)

var sdkClient *formance.APIClient

func configureSDK() {
	gatewayUrl, err := url.Parse(gatewayServer.URL)
	if err != nil {
		panic(err)
	}

	configuration := formance.NewConfiguration()
	configuration.Host = gatewayUrl.Host
	configuration.HTTPClient = &http.Client{
		Transport: httpclient.NewDebugHTTPTransport(http.DefaultTransport),
	}
	sdkClient = formance.NewAPIClient(configuration)
}

func Client() *formance.APIClient {
	return sdkClient
}
