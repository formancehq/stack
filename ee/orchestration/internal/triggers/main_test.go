package triggers

import (
	"log"
	"os"
	"testing"

	"github.com/formancehq/orchestration/internal/workflow"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/operatorservice/v1"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.temporal.io/sdk/testsuite"

	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	flag "github.com/spf13/pflag"
)

var (
	devServer *testsuite.DevServer
)

func TestMain(m *testing.M) {
	flag.Parse()

	if err := pgtesting.CreatePostgresServer(); err != nil {
		log.Fatal(err)
	}

	var err error
	devServer, err = testsuite.StartDevServer(logging.TestingContext(), testsuite.DevServerOptions{
		LogLevel: "warn",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = devServer.Client().OperatorService().AddSearchAttributes(logging.TestingContext(), &operatorservice.AddSearchAttributesRequest{
		SearchAttributes: map[string]enums.IndexedValueType{
			workflow.SearchAttributeWorkflowID: enums.INDEXED_VALUE_TYPE_TEXT,
			SearchAttributeTriggerID:           enums.INDEXED_VALUE_TYPE_TEXT,
		},
		Namespace: "default",
	})
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	if err := devServer.Stop(); err != nil {
		log.Println("unable to stop temporal server", err)
	}
	if err := pgtesting.DestroyPostgresServer(); err != nil {
		log.Println("unable to stop postgres server", err)
	}
	os.Exit(code)
}
