package es

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/opensearch-project/opensearch-go"
)

type Response struct {
	Total    int                      `json:"total"`
	Deleted  int                      `json:"deleted"`
	Failures []map[string]interface{} `json:"failures"`
	Took     int                      `json:"took"`
	TimedOut bool                     `json:"timed_out"`
}

func NewElasticSearchClient(config *v1beta3.ElasticSearchConfig) (*opensearch.Client, error) {
	httpTransport := http.DefaultTransport
	httpTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	opensearchConfig := opensearch.Config{
		Addresses: []string{config.Scheme + "://" + config.Host + ":" + strconv.FormatUint(uint64(config.Port), 10)},
		Transport: httpTransport,

		UseResponseCheckOnly: true,
	}

	if config.BasicAuth != nil {
		opensearchConfig.Username = config.BasicAuth.Username
		opensearchConfig.Password = config.BasicAuth.Password
	}

	return opensearch.NewClient(opensearchConfig)
}
