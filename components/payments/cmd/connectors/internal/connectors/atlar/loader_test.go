package atlar

import (
	"context"
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

	assert.Equal(t, Name, loader.Name())
	assert.Equal(t, 50, loader.AllowTasks())
	assert.Equal(t, Config{
		Name:          "ATLAR",
		BaseUrl:       "https://api.atlar.com",
		PollingPeriod: connectors.Duration{Duration: 2 * time.Minute},
		ApiConfig:     ApiConfig{PageSize: 25},
	}, loader.ApplyDefaults(config))

	assert.EqualValues(t, newConnector(logger, config), loader.Load(logger, config))
}
