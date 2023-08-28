package workflows

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("workflows",
		fctl.WithAliases("w", "work"),
		fctl.WithShortDescription("Workflows management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewCreateCommand(),
			NewRunCommand(),
			NewShowCommand(),
			NewDeleteCommand(),
		),
		fctl.WithCommandScopesFlags(fctl.Organization, fctl.Stack),
	)
}
