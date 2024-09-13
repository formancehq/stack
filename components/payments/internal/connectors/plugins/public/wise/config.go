package wise

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

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

	p, _ := pem.Decode([]byte(c.WebhookPublicKey))
	if p == nil {
		return errors.Wrap(models.ErrInvalidConfig, "invalid webhook public key in config")
	}

	publicKey, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return errors.Wrap(models.ErrInvalidConfig, fmt.Sprintf("invalid webhook public key in config: %v", err))
	}

	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		c.webhookPublicKey = pub
	default:
		return errors.Wrap(models.ErrInvalidConfig, "invalid webhook public key in config")
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
