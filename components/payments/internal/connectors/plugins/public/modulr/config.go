package modulr

import (
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type Config struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
	Endpoint  string `json:"endpoint"`
}

func (c Config) validate() error {
	if c.APIKey == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing api key in config")
	}

	if c.APISecret == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing api secret in config")
	}

	if c.Endpoint == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing endpoint in config")
	}

	return nil
}

func unmarshalAndValidateConfig(payload []byte) (Config, error) {
	var config Config
	if err := json.Unmarshal(payload, &config); err != nil {
		return Config{}, errors.Wrap(models.ErrInvalidConfig, err.Error())
	}

	return config, config.validate()
}
