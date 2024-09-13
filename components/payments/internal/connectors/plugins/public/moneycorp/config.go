package moneycorp

import (
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type Config struct {
	ClientID string `json:"clientID"`
	APIKey   string `json:"apiKey"`
	Endpoint string `json:"endpoint"`
}

func (c Config) validate() error {
	if c.ClientID == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing clientID in config")
	}

	if c.APIKey == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing api key in config")
	}

	if c.Endpoint == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing endpoint in config")
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
