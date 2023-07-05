package bankingcircle

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors"
	"github.com/formancehq/payments/internal/app/connectors/configtemplate"
)

// PFX is not handle very well in Go if we have more than one certificate
// in the pfx data.
// To be safe for every user to pass the right data to the connector, let's
// use two config parameters instead of one: the user certificate and the key
// associated.
// To extract them for a pfx file, you can use the following commands:
// openssl pkcs12 -in PC20230412293693.pfx -clcerts -nokeys | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > clientcert.cer
// openssl pkcs12 -in PC20230412293693.pfx -nocerts -nodes | sed -ne '/-BEGIN PRIVATE KEY-/,/-END PRIVATE KEY-/p' > clientcert.key
type Config struct {
	Username              string              `json:"username" yaml:"username" bson:"username"`
	Password              string              `json:"password" yaml:"password" bson:"password"`
	Endpoint              string              `json:"endpoint" yaml:"endpoint" bson:"endpoint"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint" yaml:"authorizationEndpoint" bson:"authorizationEndpoint"`
	UserCertificate       string              `json:"userCertificate" yaml:"userCertificate" bson:"userCertificate"`
	UserCertificateKey    string              `json:"userCertificateKey" yaml:"userCertificateKey" bson:"userCertificateKey"`
	PollingPeriod         connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

// String obfuscates sensitive fields and returns a string representation of the config.
// This is used for logging.
func (c Config) String() string {
	return fmt.Sprintf("username=%s, password=****, endpoint=%s, authorizationEndpoint=%s", c.Username, c.Endpoint, c.AuthorizationEndpoint)
}

func (c Config) Validate() error {
	if c.Username == "" {
		return ErrMissingUsername
	}

	if c.Password == "" {
		return ErrMissingPassword
	}

	if c.Endpoint == "" {
		return ErrMissingEndpoint
	}

	if c.AuthorizationEndpoint == "" {
		return ErrMissingAuthorizationEndpoint
	}

	if c.UserCertificate == "" {
		return ErrMissingUserCertificate
	}

	if c.UserCertificateKey == "" {
		return ErrMissingUserCertificatePassphrase
	}

	return nil
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("username", configtemplate.TypeString, true)
	cfg.AddParameter("password", configtemplate.TypeString, true)
	cfg.AddParameter("endpoint", configtemplate.TypeString, true)
	cfg.AddParameter("authorizationEndpoint", configtemplate.TypeString, true)
	cfg.AddParameter("userCertificate", configtemplate.TypeLongString, true)
	cfg.AddParameter("userCertificateKey", configtemplate.TypeLongString, true)
	cfg.AddParameter("pollingPeriod", configtemplate.TypeDurationNs, false)

	return Name.String(), cfg
}
