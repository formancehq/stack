package modules

import (
	"github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/wallets/cmd"
)

var Wallets = internal.NewModule("wallets").
	WithServices(
		internal.NewCommandService("wallets", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"serve",
					"--auth-enabled=false",
					"--stack-client-id=global",
					"--stack-client-secret=global",
					"--stack-url=" + test.GatewayURL(),
					"--listen=0.0.0.0:0",
				}
			}),
	)
