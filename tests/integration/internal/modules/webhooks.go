package modules

import (
	"fmt"
	"github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/webhooks/cmd"
)

var Webhooks = internal.NewModule("webhooks").
	WithCreateDatabase().
	WithServices(
		internal.NewCommandService("webhooks", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"serve",
					"--storage-postgres-conn-string=" + test.GetDatabaseSourceName("webhooks"),
					"--listen=0.0.0.0:0",
					"--worker",
					"--publisher-nats-enabled",
					"--publisher-nats-client-id=webhooks",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--kafka-topics=%s-ledger", test.ID()),
					"--retries-cron=1s",
					"--retries-schedule=1s,1s",
				}
			}),
	)
