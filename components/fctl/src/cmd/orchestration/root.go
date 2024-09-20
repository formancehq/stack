package orchestration

import (
	"github.com/formancehq/fctl/cmd/orchestration/instances"
	"github.com/formancehq/fctl/cmd/orchestration/triggers"
	"github.com/formancehq/fctl/cmd/orchestration/workflows"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("orchestration",
		fctl.WithAliases("orch", "or"),
		fctl.WithShortDescription("Orchestration"),
		fctl.WithHidden(),
		fctl.WithChildCommands(
			instances.NewCommand(),
			workflows.NewCommand(),
			triggers.NewCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			return fctl.NewStackStore(cmd)
		}),
	)
}
