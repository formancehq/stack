package service

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"

	testutils "github.com/formancehq/webhooks/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.StartPostgresServer()
	var err error
	database, err = testutils.GetStoreProvider()
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	m.Run()
	testutils.StopPostgresServer()
}
