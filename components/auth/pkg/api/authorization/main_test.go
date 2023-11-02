package authorization

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/zitadel/logging"
)

func TestMain(t *testing.M) {
	if err := pgtesting.CreatePostgresServer(); err != nil {
		logging.Errorf("Unable to start postgres server: %s", err)
		os.Exit(1)
	}
	code := t.Run()

	if err := pgtesting.DestroyPostgresServer(); err != nil {
		logging.Errorf("Unable to stop postgres server: %s", err)
	}
	os.Exit(code)
}
