package wise

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type Config struct {
	APIKey           string `json:"apiKey"`
	WebhookPublicKey string `json:"webhookPublicKey"`

	webhookPublicKey *rsa.PublicKey `json:"-"`
}

func (c *Config) validate() error {
	if c.APIKey == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing api key in config")
	}

	if c.WebhookPublicKey == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing webhook public key in config")
	}

	publicKey, err := x509.ParsePKCS1PublicKey([]byte(c.WebhookPublicKey))
	if err != nil {
		return errors.Wrap(models.ErrInvalidConfig, "invalid webhook public key in config")
	}
	c.webhookPublicKey = publicKey

	return nil
}

func unmarshalAndValidateConfig(payload json.RawMessage) (Config, error) {
	var config Config
	if err := json.Unmarshal(payload, &config); err != nil {
		return Config{}, errors.Wrap(models.ErrInvalidConfig, err.Error())
	}

	return config, config.validate()
}
