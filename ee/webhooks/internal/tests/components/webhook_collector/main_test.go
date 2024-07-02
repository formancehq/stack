package webhookcollector

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	component "github.com/formancehq/webhooks/internal/components/commons"
	webhookcollector "github.com/formancehq/webhooks/internal/components/webhook_collector"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	testutils "github.com/formancehq/webhooks/internal/tests"
)

var Database storage.PostgresStore
var WebhookCollector webhookcollector.Collector

func TestMain(m *testing.M) {
	testutils.StartPostgresServer()
	var err error
	Database, err = testutils.GetStoreProvider()
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	WebhookCollector = *webhookcollector.NewCollector(component.DefaultRunnerParams(),
		Database, testutils.NewHTTPClient())

	m.Run()
	testutils.StopPostgresServer()
}
