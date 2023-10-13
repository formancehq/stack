package modules

import (
	"fmt"
	"github.com/formancehq/payments/cmd"
	"github.com/formancehq/stack/tests/integration/internal"
)

var Payments = internal.NewModule("payments").
	WithCreateDatabase().
	WithServices(
		internal.NewCommandService("payments", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"serve",
					"--postgres-uri=" + test.GetDatabaseSourceName("payments"),
					"--config-encryption-key=encryption-key",
					"--publisher-nats-enabled",
					"--publisher-nats-client-id=payments",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--publisher-topic-mapping=*:%s-payments", test.ID()),
					"--listen=0.0.0.0:0",
					"--auto-migrate",
				}
			}),
	)
