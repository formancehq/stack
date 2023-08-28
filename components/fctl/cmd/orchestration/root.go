package orchestration

import (
	"github.com/formancehq/fctl/cmd/orchestration/instances"
	"github.com/formancehq/fctl/cmd/orchestration/workflows"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("orchestration",
		fctl.WithAliases("orch", "or"),
		fctl.WithShortDescription("Orchestration"),
		fctl.WithHidden(),
		fctl.WithChildCommands(
			instances.NewCommand(),
			workflows.NewCommand(),
		),
		fctl.WithCommandScopesFlags(fctl.Organization, fctl.Stack),
	)
}
