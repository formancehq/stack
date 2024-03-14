package dummypay

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

// Config is the configuration for the dummy payment connector.
type Config struct {
	Name string `json:"name" yaml:"name" bson:"name"`

	// Directory is the directory where the files are stored.
	Directory string `json:"directory" yaml:"directory" bson:"directory"`

	// FilePollingPeriod is the period between file polling.
	FilePollingPeriod connectors.Duration `json:"filePollingPeriod" yaml:"filePollingPeriod" bson:"filePollingPeriod"`

	// PrefixFileToIngest is the prefix of the file to ingest.
	PrefixFileToIngest string `json:"prefixFileToIngest" yaml:"prefixFileToIngest" bson:"prefixFileToIngest"`

	// NumberOfAccountsPreGenerated is the number of accounts to pre-generate.
	NumberOfAccountsPreGenerated int `json:"numberOfAccountsPreGenerated" yaml:"numberOfAccountsPreGenerated" bson:"numberOfAccountsPreGenerated"`
	// NumberOfPaymentsPreGenerated is the number of payments to pre-generate.
	NumberOfPaymentsPreGenerated int `json:"numberOfPaymentsPreGenerated" yaml:"numberOfPaymentsPreGenerated" bson:"numberOfPaymentsPreGenerated"`
}

// String returns a string representation of the configuration.
func (c Config) String() string {
	return fmt.Sprintf("directory=%s, filePollingPeriod=%s",
		c.Directory, c.FilePollingPeriod.String())
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) ConnectorName() string {
	return c.Name
}

// Validate validates the configuration.
func (c Config) Validate() error {
	// require directory path to be present
	if c.Directory == "" {
		return ErrMissingDirectory
	}

	// check if file polling period is set properly
	if c.FilePollingPeriod.Duration <= 0 {
		return fmt.Errorf("filePollingPeriod must be greater than 0: %w",
			ErrFilePollingPeriodInvalid)
	}

	if c.Name == "" {
		return fmt.Errorf("name must be set: %w", ErrMissingName)
	}

	return nil
}

func (c Config) BuildTemplate() (string, configtemplate.Config) {
	cfg := configtemplate.NewConfig()

	cfg.AddParameter("name", configtemplate.TypeString, name.String(), false)
	cfg.AddParameter("directory", configtemplate.TypeString, "", true)
	cfg.AddParameter("filePollingPeriod", configtemplate.TypeDurationNs, "", true)
	cfg.AddParameter("prefixFileToIngest", configtemplate.TypeString, "", false)
	cfg.AddParameter("numberOfAccountsPreGenerated", configtemplate.TypeDurationUnsignedInteger, "0", false)
	cfg.AddParameter("numberOfPaymentsPreGenerated", configtemplate.TypeDurationUnsignedInteger, "0", false)

	return name.String(), cfg
}
