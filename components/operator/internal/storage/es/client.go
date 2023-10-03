package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/go-logr/logr"
	"github.com/opensearch-project/opensearch-go"
)

func NewElasticSearchClient(config *v1beta3.ElasticSearchConfig) (*opensearch.Client, error) {
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

const (
	stacksIndex = "stacks"
)

func DropESIndex(client *opensearch.Client, logger logr.Logger, stackName string, ctx context.Context) error {
	var (
		buf bytes.Buffer
	)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"stack": stackName,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Error(err, "ELK: Error during json encoding")

		return err
	}
	body := bytes.NewReader(buf.Bytes())
	response, err := client.DeleteByQuery([]string{stacksIndex}, body, client.DeleteByQuery.WithContext(ctx))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return fmt.Errorf("ELK status: %d", response.StatusCode)
	}

	return nil
}
