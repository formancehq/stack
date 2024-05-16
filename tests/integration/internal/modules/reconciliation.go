package modules

import (
	"github.com/formancehq/reconciliation/cmd"
	"github.com/formancehq/stack/tests/integration/internal"
)

var Reconciliation = internal.NewModule("reconciliation").
	WithCreateDatabase().
	WithServices(
		internal.NewCommandService("reconciliation", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"serve",
					"--auth-enabled=false",
					"--stack-client-id=global",
					"--stack-client-secret=global",
					"--stack-url=" + test.GatewayURL(),
					"--listen=0.0.0.0:0",
					"--debug",
					"--postgres-uri=" + test.GetDatabaseSourceName("reconciliation"),
					"--auto-migrate",
				}
			}),
	)
