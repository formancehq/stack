package webhookrunner

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	component "github.com/formancehq/webhooks/internal/components/commons"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	testutils "github.com/formancehq/webhooks/internal/tests"
)

var Database storage.PostgresStore
var WebhookRunner component.WebhookRunner

func TestMain(m *testing.M){
	testutils.StartPostgresServer()
	var err error 
	Database, err = testutils.GetStoreProvider()
	WebhookRunner = *component.NewWebhookRunner(component.DefaultRunnerParams(), 
												nil, 
												Database, 
												testutils.NewHTTPClient())
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	m.Run()
	testutils.StopPostgresServer()
}


