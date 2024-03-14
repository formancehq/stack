package atlar

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/assert"
)

// TestLoader tests the loader.
func TestLoader(t *testing.T) {
	t.Parallel()

	config := Config{}
	logger := logging.FromContext(context.TODO())

	loader := NewLoader()

	assert.Equal(t, name, loader.Name())
	assert.Equal(t, 50, loader.AllowTasks())

	baseUrl, err := url.Parse("https://api.atlar.com")
	assert.Nil(t, err)

	assert.Equal(t, Config{
		Name:                                  "ATLAR",
		BaseUrl:                               *baseUrl,
		PollingPeriod:                         connectors.Duration{Duration: 2 * time.Minute},
		TransferInitiationStatusPollingPeriod: connectors.Duration{Duration: 2 * time.Minute},
		ApiConfig:                             ApiConfig{PageSize: 25},
	}, loader.ApplyDefaults(config))

	assert.EqualValues(t, newConnector(logger, config), loader.Load(logger, config))
}
