package es

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/go-logr/logr"
)

const (
	stacksIndex = "stacks"
)

func DropESIndex(config *v1beta3.ElasticSearchConfig, logger logr.Logger, stackName string) error {
	client, err := NewElasticSearchClient(config)
	if err != nil {
		logger.Error(err, "ES client error")
		return err
	}

	var (
		buf bytes.Buffer
		res Response
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
	response, err := client.DeleteByQuery([]string{stacksIndex}, body)

	if err != nil {
		logger.Error(err, "ELK: During index drop es")

		return err
	}

	defer response.Body.Close()
	if response.StatusCode > 300 {
		logger.Error(fmt.Errorf("ELK: Status code: %d", response.StatusCode), "ELK: Error during index drop")
		return fmt.Errorf("ELK: Status code: %d", response.StatusCode)
	}

	logger.Info(fmt.Sprintf("ELK Status code: %d", response.StatusCode))

	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		logger.Error(err, "ELK: Error decoding es body")
		return err
	}

	if response.StatusCode == 200 {
		logger.Info("ES Index Dropped")
		logger.Info(fmt.Sprint("Total: ", res.Total))
		logger.Info(fmt.Sprint("Deleted: ", res.Deleted))
		logger.Info(fmt.Sprint("Failures: ", res.Failures))
	}

	return nil
}
