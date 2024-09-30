package stripe

import (
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type Config struct {
	APIKey string `json:"apiKey"`
}

func (c Config) validate() error {
	if c.APIKey == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing api key in config")
	}
	return nil
}

func unmarshalAndValidateConfig(payload json.RawMessage) (Config, error) {
	var config Config
	if err := json.Unmarshal(payload, &config); err != nil {
		return Config{}, errors.Wrap(models.ErrInvalidConfig, err.Error())
	}

	return config, config.validate()
}
