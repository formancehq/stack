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
					"--auth-enabled=false",
					"--postgres-uri=" + test.GetDatabaseSourceName("webhooks"),
					"--listen=0.0.0.0:0",
					"--worker",
					"--publisher-nats-enabled",
					"--publisher-nats-client-id=webhooks",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--kafka-topics=%s-ledger", test.ID()),
					fmt.Sprintf("--kafka-topics=%s-payments", test.ID()),
					"--retry-period=1s",
					"--min-backoff-delay=1s",
					"--abort-after=3s",
					"--auto-migrate=true",
				}
			}),
	)
