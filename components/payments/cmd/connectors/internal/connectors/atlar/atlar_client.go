package atlar

import (
	atlar_client "github.com/get-momo/atlar-v1-go-client/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func createAtlarClient(config *Config) *atlar_client.Rest {
	transport := httptransport.New(
		config.BaseUrl.Host,
		config.BaseUrl.Path,
		[]string{config.BaseUrl.Scheme},
	)
	basicAuth := httptransport.BasicAuth(config.AccessKey, config.Secret)
	transport.DefaultAuthentication = basicAuth
	client := atlar_client.New(transport, strfmt.Default)
	return client
}
