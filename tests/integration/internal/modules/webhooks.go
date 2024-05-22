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
					"start",
					"--auth-enabled=false",
					"--postgres-uri=" + test.GetDatabaseSourceName("webhooks"),
					"--listen=0.0.0.0:0",
					"--publisher-nats-enabled",
					"--publisher-nats-client-id=webhooks",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--kafka-topics=%s-ledger", test.ID()),
					fmt.Sprintf("--kafka-topics=%s-payments", test.ID()),
					"--max-call=20",
					"--max-retry=60",
					"--time-out=2000",
					"--delay-pull=1",
					"--auto-migrate=true",
				}
			}),
	)
