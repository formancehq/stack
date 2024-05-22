package webhookcontroller

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	testutils "github.com/formancehq/webhooks/internal/tests"
)

var Database storage.PostgresStore


func TestMain(m *testing.M){
	testutils.StartPostgresServer()
	var err error 
	Database, err = testutils.GetStoreProvider()
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	m.Run()
	testutils.StopPostgresServer()
}


