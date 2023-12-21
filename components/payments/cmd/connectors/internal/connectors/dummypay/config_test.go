package dummypay

import (
	"os"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/stretchr/testify/assert"
)

// TestConfigString tests the string representation of the config.
func TestConfigString(t *testing.T) {
	t.Parallel()

	config := Config{
		Directory:         "test",
		FilePollingPeriod: connectors.Duration{Duration: time.Second},
	}

	assert.Equal(t, "directory=test, filePollingPeriod=1s", config.String())
}

// TestConfigValidate tests the validation of the config.
func TestConfigValidate(t *testing.T) {
	t.Parallel()

	var config Config

	config.Name = "test1"

	// fail on missing directory
	assert.EqualError(t, config.Validate(), ErrMissingDirectory.Error())

	// fail on missing RW access to directory
	config.Directory = "/non-existing"
	assert.Error(t, config.Validate())

	// set directory with RW access
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}

	config.Directory = userHomeDir

	// fail on invalid file polling period
	config.FilePollingPeriod.Duration = -1
	assert.ErrorIs(t, config.Validate(), ErrFilePollingPeriodInvalid)

	// success
	config.FilePollingPeriod.Duration = 1
	assert.NoError(t, config.Validate())
}
