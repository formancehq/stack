package internal

import (
	"net/url"
	"testing"

	"github.com/formancehq/formance-sdk-go"
)

var sdkClient *formance.APIClient

func configureSDK() {
	gatewayUrl, err := url.Parse(gatewayServer.URL)
	if err != nil {
		panic(err)
	}

	configuration := formance.NewConfiguration()
	configuration.Debug = testing.Verbose()
	configuration.Host = gatewayUrl.Host
	sdkClient = formance.NewAPIClient(configuration)
}

func Client() *formance.APIClient {
	return sdkClient
}
