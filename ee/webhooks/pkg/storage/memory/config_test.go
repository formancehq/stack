package memory

import (
	"testing"

	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
)

func TestConfigStore(t *testing.T) {
	newStore := func(*testing.T) webhooks.ConfigStore {
		return NewConfigStore()
	}
	storage.TestConfigStore(t, newStore)
}
