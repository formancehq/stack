package webhookworker

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	component "github.com/formancehq/webhooks/internal/components/commons"
	webhookworker "github.com/formancehq/webhooks/internal/components/webhook_worker"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	testutils "github.com/formancehq/webhooks/internal/tests"
)

var Database storage.PostgresStore
var WebhookWorker webhookworker.Worker

func TestMain(m *testing.M){
	testutils.StartPostgresServer()
	var err error 
	Database, err = testutils.GetStoreProvider()
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	WebhookWorker = *webhookworker.NewWorker(component.DefaultRunnerParams(), 
	Database, testutils.NewHTTPClient())

	m.Run()
	testutils.StopPostgresServer()
}


