package adyen

import (
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type Config struct {
	APIKey             string `json:"apiKey"`
	WebhookUsername    string `json:"webhookUsername"`
	WebhookPassword    string `json:"webhookPassword"`
	CompanyID          string `json:"companyID"`
	LiveEndpointPrefix string `json:"liveEndpointPrefix"`
}

func (c Config) validate() error {
	if c.APIKey == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing apiKey in config")
	}

	if c.CompanyID == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing companyID in config")
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
