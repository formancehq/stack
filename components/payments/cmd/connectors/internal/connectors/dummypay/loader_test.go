package dummypay

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

	assert.Equal(t, name, loader.Name())
	assert.Equal(t, 10, loader.AllowTasks())
	assert.Equal(t, Config{
		Name:              "DUMMY-PAY",
		FilePollingPeriod: connectors.Duration{Duration: 10 * time.Second},
	}, loader.ApplyDefaults(config))

	assert.EqualValues(t, newConnector(logger, config, newFS()), loader.Load(logger, config))
}
