package modules

import (
	"fmt"

	"github.com/formancehq/orchestration/cmd"
	"github.com/formancehq/stack/tests/integration/internal"
)

var Orchestration = internal.NewModule("orchestration").
	WithCreateDatabase().
	WithServices(
		internal.NewCommandService("orchestration", cmd.NewRootCommand).
			WithArgs(func(test *internal.Test) []string {
				return []string{
					"serve",
					"--auth-enabled=false",
					"--listen=0.0.0.0:0",
					"--postgres-uri=" + test.GetDatabaseSourceName("orchestration"),
					"--stack-client-id=global",
					"--stack-client-secret=global",
					"--stack-url=" + test.GatewayURL(),
					"--temporal-address=" + internal.GetTemporalAddress(),
					"--temporal-task-queue=" + test.ID(),
					"--temporal-init-search-attributes",
					"--worker",
					"--publisher-nats-enabled",
					"--publisher-nats-client-id=orchestration",
					"--publisher-nats-url=" + internal.GetNatsAddress(),
					fmt.Sprintf("--publisher-topic-mapping=*:%s-orchestration", test.ID()),
					fmt.Sprintf("--topics=%s-ledger", test.ID()),
					fmt.Sprintf("--topics=%s-payments", test.ID()),
					"--debug",
				}
			}),
	)
