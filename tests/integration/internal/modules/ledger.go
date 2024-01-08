package modules

import (
	"fmt"

	"github.com/formancehq/ledger/cmd"
	"github.com/formancehq/stack/tests/integration/internal"
)

var Ledger = internal.NewModule("ledger").
	WithCreateDatabase().
	WithServices(
		internal.NewCommandService("ledger", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"serve",
					"--publisher-nats-enabled",
					"--auth-enabled=false",
					"--publisher-nats-client-id=ledger",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--publisher-topic-mapping=*:%s-ledger", test.ID()),
					"--storage-postgres-conn-string=" + test.GetDatabaseSourceName("ledger"),
					"--json-formatting-logger=false",
					"--bind=0.0.0.0:0", // Random port
					"--debug",
				}
			}),
	)
