package es

import (
	"crypto/tls"
	"net/http"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/opensearch-project/opensearch-go"
)

func NewElasticSearchClient(config v1beta3.ElasticSearchConfig) (*opensearch.Client, error) {
	httpTransport := http.DefaultTransport
	httpTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: config.TLS.SkipCertVerify,
	}

	opensearchConfig := opensearch.Config{
		Addresses:            []string{config.Endpoint()},
		Transport:            httpTransport,
		UseResponseCheckOnly: true,
	}

	if config.BasicAuth != nil {
		opensearchConfig.Username = config.BasicAuth.Username
		opensearchConfig.Password = config.BasicAuth.Password
	}

	return opensearch.NewClient(opensearchConfig)
}
