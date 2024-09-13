package bankingcircle

import (
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type Config struct {
	Username              string `json:"username" yaml:"username" `
	Password              string `json:"password" yaml:"password" `
	Endpoint              string `json:"endpoint" yaml:"endpoint"`
	AuthorizationEndpoint string `json:"authorizationEndpoint" yaml:"authorizationEndpoint" `
	UserCertificate       string `json:"userCertificate" yaml:"userCertificate" `
	UserCertificateKey    string `json:"userCertificateKey" yaml:"userCertificateKey"`
}

func (c Config) validate() error {
	if c.Username == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing username in config")
	}

	if c.Password == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing password in config")
	}

	if c.Endpoint == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing endpoint in config")
	}

	if c.AuthorizationEndpoint == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing authorization endpoint in config")
	}

	if c.UserCertificate == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing user certificate in config")
	}

	if c.UserCertificateKey == "" {
		return errors.Wrap(models.ErrInvalidConfig, "missing user certificate key in config")
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
